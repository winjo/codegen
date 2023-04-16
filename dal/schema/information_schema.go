package schema

import (
	"database/sql"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/samber/lo"
	"github.com/winjo/codegen/dal/util"
)

type (
	InformationSchemaModel struct {
		dbName string
		db     *sql.DB
	}

	DbColumn struct {
		ColumnName      string
		DataType        string
		ColumnType      string
		Extra           string
		ColumnComment   string
		ColumnDefault   *string
		IsNullable      string
		OrdinalPosition int
	}

	DbIndex struct {
		IndexName  string
		NonUnique  int
		SeqInIndex int
		ColumnName string
	}

	Column struct {
		Name     string  `json:"name"`
		Type     string  `json:"type"`
		Length   *int    `json:"length"`
		Unsigned bool    `json:"unsigned"`
		Nullable bool    `json:"nullable"`
		Default  *string `json:"default"`
		Extra    string  `json:"extra"`
		Comment  string  `json:"comment"`
	}

	Index struct {
		Name    string   `json:"name"`
		Unique  bool     `json:"unique"`
		Seqs    []int    `json:"-"`
		Columns []string `json:"columns"`
	}

	Table struct {
		DB      string    `json:"db"`
		Name    string    `json:"name"`
		Columns []*Column `json:"columns"`
		Indexes []*Index  `json:"indexes"`
	}
)

func NewInformationSchemaModel(dsn string) *InformationSchemaModel {
	config, err := mysql.ParseDSN(dsn)
	util.AssertNotNil(err)

	infoSchemaConfig := config.Clone()
	infoSchemaConfig.DBName = "information_schema"
	connector, err := mysql.NewConnector(infoSchemaConfig)
	util.AssertNotNil(err)

	db := sql.OpenDB(connector)

	return &InformationSchemaModel{db: db, dbName: config.DBName}
}

func (m *InformationSchemaModel) GetAllTables() []*Table {
	query := `select TABLE_NAME from TABLES where TABLE_SCHEMA = ?`
	rows, err := m.db.Query(query, m.dbName)
	util.AssertNotNil(err)

	var tables []*Table
	for rows.Next() {
		var tableName string
		rows.Scan(&tableName)
		table := m.newTableData(tableName)
		tables = append(tables, table)
	}

	return tables
}

func (m *InformationSchemaModel) newTableData(tableName string) *Table {
	columns := lo.Map(m.FindColumns(tableName), func(c *DbColumn, index int) *Column {
		temp := lo.Map(strings.Split(c.ColumnType, " "), func(t string, i int) string {
			return strings.TrimSpace(t)
		})
		ret := regexp.MustCompile(`^(?:.*)\((.*)\)$`).FindStringSubmatch(temp[0])
		var length *int
		if ret != nil {
			v, err := strconv.Atoi(ret[1])
			util.AssertNotNil(err)
			length = &v
		}
		column := &Column{
			Name:     c.ColumnName,
			Type:     strings.ToLower(c.DataType),
			Length:   length,
			Nullable: c.IsNullable == "YES",
			Unsigned: strings.Contains(strings.ToLower(c.ColumnType), "unsigned"),
			Default:  c.ColumnDefault,
			Extra:    c.Extra,
			Comment:  c.ColumnComment,
		}
		return column
	})

	dbIndexes := m.FindIndexes(tableName)
	indexes := make([]*Index, 0, len(dbIndexes))
	indexGroup := make(map[string]*Index)
	for _, index := range dbIndexes {
		indexItem, ok := indexGroup[index.IndexName]
		if !ok {
			indexGroup[index.IndexName] = &Index{
				Name:    index.IndexName,
				Unique:  index.NonUnique == 0,
				Columns: []string{index.ColumnName},
				Seqs:    []int{index.SeqInIndex},
			}
			indexes = append(indexes, indexGroup[index.IndexName])
		} else {
			indexItem.Columns = append(indexItem.Columns, index.ColumnName)
			indexItem.Seqs = append(indexItem.Seqs, index.SeqInIndex)
		}
	}

	for _, index := range indexes {
		sort.Slice(index.Columns, func(i, j int) bool {
			return index.Seqs[i] < index.Seqs[j]
		})
	}

	return &Table{
		DB:      m.dbName,
		Name:    tableName,
		Columns: columns,
		Indexes: indexes,
	}
}

func (m *InformationSchemaModel) FindColumns(tableName string) []*DbColumn {
	query := "SELECT COLUMN_NAME,DATA_TYPE,COLUMN_TYPE,EXTRA,COLUMN_COMMENT,COLUMN_DEFAULT,IS_NULLABLE from COLUMNS WHERE TABLE_SCHEMA = ? and TABLE_NAME = ? ORDER BY ORDINAL_POSITION"
	rows, err := m.db.Query(query, m.dbName, tableName)
	util.AssertNotNil(err)

	var reply []*DbColumn
	for rows.Next() {
		var c DbColumn
		util.AssertNotNil(rows.Scan(
			&c.ColumnName,
			&c.DataType,
			&c.ColumnType,
			&c.Extra,
			&c.ColumnComment,
			&c.ColumnDefault,
			&c.IsNullable,
		))
		reply = append(reply, &c)
	}

	return reply
}

func (m *InformationSchemaModel) FindIndexes(tableName string) []*DbIndex {
	query := "SELECT INDEX_NAME,NON_UNIQUE,SEQ_IN_INDEX,COLUMN_NAME from STATISTICS WHERE TABLE_SCHEMA = ? and TABLE_NAME = ?"
	rows, err := m.db.Query(query, m.dbName, tableName)
	util.AssertNotNil(err)

	var reply []*DbIndex
	for rows.Next() {
		var item DbIndex
		util.AssertNotNil(rows.Scan(
			&item.IndexName,
			&item.NonUnique,
			&item.SeqInIndex,
			&item.ColumnName,
		))
		reply = append(reply, &item)
	}

	return reply
}
