package schema

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/samber/lo"
	"github.com/winjo/codegen/dal/util"
)

func GenStructure(tables []*Table) {
	for _, table := range tables {
		genStructureJSON(table)
		genStructureSql(table)
	}
}

func genStructureJSON(table *Table) {
	path := fmt.Sprintf("dal/structure/%s.json", table.Name)
	tableBuf, err := json.MarshalIndent(table, "", "  ")
	util.AssertNotNil(err)
	os.WriteFile(path, tableBuf, os.ModePerm)
}

func genStructureSql(table *Table) {
	path := fmt.Sprintf("dal/structure/%s.sql", table.Name)
	var items []string
	for _, c := range table.Columns {
		var item []string
		item = append(item, wrapName(c.Name))
		if c.Length != nil {
			item = append(item, fmt.Sprintf("%s(%d)", c.Type, *c.Length))
		} else {
			item = append(item, c.Type)
		}
		if c.Unsigned {
			item = append(item, "unsigned")
		}
		if c.Nullable {
			item = append(item, "NULL")
		} else {
			item = append(item, "NOT NULL")
		}
		if c.Extra == "auto_increment" {
			item = append(item, "AUTO_INCREMENT")
		} else if strings.Contains(strings.ToLower(c.Extra), "on update current_timestamp") {
			item = append(item, "on update CURRENT_TIMESTAMP")
		}

		defaultValue := getDefault(c)
		if defaultValue != "" {
			item = append(item, fmt.Sprintf("DEFAULT %s", defaultValue))
		}

		if c.Comment != "" {
			item = append(item, fmt.Sprintf("COMMENT '%s'", strings.ReplaceAll(c.Comment, "\n", "\\n")))
		}
		items = append(items, strings.Join(item, " "))
	}

	var primayIndex string
	var uniqueIndex []string
	var notUniqueIndex []string
	for _, i := range table.Indexes {
		if i.Name == "PRIMARY" {
			primayIndex = fmt.Sprintf("PRIMARY KEY (%s)", wrapName(i.Columns[0]))
		} else {
			columns := lo.Map(i.Columns, func(c string, index int) string {
				return wrapName(c)
			})
			str := fmt.Sprintf("KEY %s (%s)", wrapName(i.Name), strings.Join(columns, ","))
			if i.Unique {
				uniqueIndex = append(uniqueIndex, "UNIQUE "+str)
			} else {
				notUniqueIndex = append(notUniqueIndex, str)
			}
		}
	}
	if primayIndex != "" {
		items = append(items, primayIndex)
	}
	if len(uniqueIndex) != 0 {
		items = append(items, uniqueIndex...)
	}
	if len(notUniqueIndex) != 0 {
		items = append(items, notUniqueIndex...)
	}

	// tableName := fmt.Sprintf("`%s`.`%s`", table.DB, table.Name)
	tableOptions := strings.Join(lo.Map(items, func(item string, index int) string {
		return " " + item
	}), ",\n")

	sql := strings.Join([]string{
		fmt.Sprintf("DROP TABLE IF EXISTS %s;", table.Name),
		fmt.Sprintf("CREATE TABLE %s (", table.Name),
		tableOptions,
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;\n",
	}, "\n")

	os.WriteFile(path, []byte(sql), os.ModePerm)
}

func wrapName(name string) string {
	return fmt.Sprintf("`%s`", name)
}

func getDefault(c *Column) string {
	if c.Default == nil && c.Nullable {
		return "NULL"
	}

	if c.Type == "timestamp" && *c.Default == "CURRENT_TIMESTAMP" {
		return *c.Default
	}

	if c.Default != nil {
		return fmt.Sprintf("'%s'", *c.Default)
	}

	return ""
}
