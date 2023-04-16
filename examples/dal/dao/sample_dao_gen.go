package dao

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/samber/mo"
	"gopkg.in/guregu/null.v4"
	"strings"
	"time"
)

type (
	baseSampleDAO struct {
		db          *sql.DB
		table       string
		columnNames []string
		allColumns  string
	}

	SampleData struct {
		Id          int64       `db:"id" json:"id"`
		GmtCreate   time.Time   `db:"gmt_create" json:"gmtCreate"`
		GmtModified time.Time   `db:"gmt_modified" json:"gmtModified"`
		RInt        int64       `db:"r_int" json:"rInt"`
		NInt        null.Int    `db:"n_int" json:"nInt"`
		RFloat      float64     `db:"r_float" json:"rFloat"`
		NFloat      null.Float  `db:"n_float" json:"nFloat"`
		RString     string      `db:"r_string" json:"rString"`
		NString     null.String `db:"n_string" json:"nString"`
		RTime       time.Time   `db:"r_time" json:"rTime"`
		NTime       null.Time   `db:"n_time" json:"nTime"`
		Union1      string      `db:"union1" json:"union1"`
		Union2      string      `db:"union2" json:"union2"`
		Union3      null.String `db:"union3" json:"union3"`
	}

	SamplePageData struct {
		Data  []*SampleData
		Count int
	}

	SamplePartialData struct {
		RInt    mo.Option[int64]       `db:"r_int" json:"rInt"`
		NInt    mo.Option[null.Int]    `db:"n_int" json:"nInt"`
		RFloat  mo.Option[float64]     `db:"r_float" json:"rFloat"`
		NFloat  mo.Option[null.Float]  `db:"n_float" json:"nFloat"`
		RString mo.Option[string]      `db:"r_string" json:"rString"`
		NString mo.Option[null.String] `db:"n_string" json:"nString"`
		RTime   mo.Option[time.Time]   `db:"r_time" json:"rTime"`
		NTime   mo.Option[null.Time]   `db:"n_time" json:"nTime"`
		Union1  mo.Option[string]      `db:"union1" json:"union1"`
		Union2  mo.Option[string]      `db:"union2" json:"union2"`
		Union3  mo.Option[null.String] `db:"union3" json:"union3"`
	}
)

func newSampleData() (*SampleData, []any) {
	var d SampleData
	ptrs := []any{
		&d.Id,
		&d.GmtCreate,
		&d.GmtModified,
		&d.RInt,
		&d.NInt,
		&d.RFloat,
		&d.NFloat,
		&d.RString,
		&d.NString,
		&d.RTime,
		&d.NTime,
		&d.Union1,
		&d.Union2,
		&d.Union3,
	}
	return &d, ptrs
}

func newBaseSampleDAO(db *sql.DB) *baseSampleDAO {
	columnNames := []string{"id", "gmt_create", "gmt_modified", "r_int", "n_int", "r_float", "n_float", "r_string", "n_string", "r_time", "n_time", "union1", "union2", "union3"}
	return &baseSampleDAO{
		db:          db,
		table:       "`sample`",
		columnNames: columnNames,
		allColumns:  columnsToRow(columnNames),
	}
}

func (dao *baseSampleDAO) Find(ctx context.Context, options ...QueryOption) ([]*SampleData, error) {
	option := resoveOption(options)
	query := fmt.Sprintf("select %s from %s %s %s", dao.allColumns, dao.table, option.order, option.sort)
	rows, err := dao.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return scanRows(rows, newSampleData)
}

func (dao *baseSampleDAO) Page(ctx context.Context, pageIndex int, pageSize int, options ...QueryOption) (*SamplePageData, error) {
	option := resoveOption(options)
	var data []*SampleData
	var count int

	err := parallelQuery(
		func() error {
			query := fmt.Sprintf("select %s from %s %s %s limit %d, %d", dao.allColumns, dao.table, option.order, option.sort, (pageIndex-1)*pageSize, pageSize)
			rows, err := dao.db.QueryContext(ctx, query)
			if err != nil {
				return err
			}
			data, err = scanRows(rows, newSampleData)
			return err
		},
		func() error {
			var err error
			count, err = dao.Count(ctx)
			return err
		},
	)

	if err != nil {
		return nil, err
	}

	return &SamplePageData{
		Data:  data,
		Count: count,
	}, nil
}

func (dao *baseSampleDAO) Count(ctx context.Context) (int, error) {
	query := fmt.Sprintf("select count(0) from %s", dao.table)
	var count int
	row := dao.db.QueryRowContext(ctx, query)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (dao *baseSampleDAO) FindByNTime(ctx context.Context, nTime null.Time, options ...QueryOption) ([]*SampleData, error) {
	option := resoveOption(options)
	query := fmt.Sprintf("select %s from %s where `n_time` %s %s %s", dao.allColumns, dao.table, ternary(nTime.Valid, "= ?", "is NULL"), option.order, option.sort)

	args := make([]any, 0, 1)
	if nTime.Valid {
		args = append(args, nTime)
	}

	rows, err := dao.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return scanRows(rows, newSampleData)
}

func (dao *baseSampleDAO) PageByNTime(ctx context.Context, nTime null.Time, pageIndex int, pageSize int, options ...QueryOption) (*SamplePageData, error) {
	option := resoveOption(options)
	var data []*SampleData
	var count int

	args := make([]any, 0, 1)
	if nTime.Valid {
		args = append(args, nTime)
	}

	err := parallelQuery(
		func() error {
			query := fmt.Sprintf("select %s from %s where `n_time` %s %s %s limit %d, %d", dao.allColumns, dao.table, ternary(nTime.Valid, "= ?", "is NULL"), option.order, option.sort, (pageIndex-1)*pageSize, pageSize)

			rows, err := dao.db.QueryContext(ctx, query, args...)
			if err != nil {
				return err
			}
			data, err = scanRows(rows, newSampleData)
			return err
		},
		func() error {
			query := fmt.Sprintf("select count(0) from (select `id` from %s where `n_time` %s) as a", dao.table, ternary(nTime.Valid, "= ?", "is NULL"))
			row := dao.db.QueryRowContext(ctx, query, args...)
			return row.Scan(&count)
		},
	)

	if err != nil {
		return nil, err
	}

	return &SamplePageData{
		Data:  data,
		Count: count,
	}, nil
}

func (dao *baseSampleDAO) FindInByNTime(ctx context.Context, nTimeList []null.Time, options ...QueryOption) ([]*SampleData, error) {
	option := resoveOption(options)

	containNull := false
	args := make([]any, len(nTimeList))
	for i, item := range nTimeList {
		if item.Valid {
			args[i] = nTimeList[i]
		} else {
			containNull = true
		}
	}

	if len(nTimeList) == 0 {
		return []*SampleData{}, nil
	}

	query := fmt.Sprintf("select %s from %s where `n_time` in (?%s) %s %s %s", dao.allColumns, dao.table, strings.Repeat(",?", len(nTimeList)-1), ternary(containNull, "or `n_time` is NULL", ""), option.order, option.sort)
	rows, err := dao.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return scanRows(rows, newSampleData)
}

func (dao *baseSampleDAO) PageInByNTime(ctx context.Context, nTimeList []null.Time, pageIndex int, pageSize int, options ...QueryOption) (*SamplePageData, error) {
	option := resoveOption(options)
	var data []*SampleData
	var count int

	containNull := false
	args := make([]any, len(nTimeList))
	for i, item := range nTimeList {
		if item.Valid {
			args[i] = nTimeList[i]
		} else {
			containNull = true
		}
	}

	if len(nTimeList) == 0 {
		return &SamplePageData{
			Data:  []*SampleData{},
			Count: 0,
		}, nil
	}

	err := parallelQuery(
		func() error {
			query := fmt.Sprintf("select %s from %s where n_time in (?%s) %s %s %s limit %d, %d", dao.allColumns, dao.table, strings.Repeat(",?", len(nTimeList)-1), ternary(containNull, "or `n_time` is NULL", ""), option.order, option.sort, (pageIndex-1)*pageSize, pageSize)
			rows, err := dao.db.QueryContext(ctx, query, args...)
			if err != nil {
				return err
			}
			data, err = scanRows(rows, newSampleData)
			return err
		},
		func() error {
			query := fmt.Sprintf("select count(0) from (select `n_time` from %s where n_time in (?%s) %s) as a", dao.table, strings.Repeat(",?", len(nTimeList)-1), ternary(containNull, "or `n_time` is NULL", ""))
			row := dao.db.QueryRowContext(ctx, query, args...)
			return row.Scan(&count)
		},
	)

	if err != nil {
		return nil, err
	}

	return &SamplePageData{
		Data:  data,
		Count: count,
	}, nil
}

func (dao *baseSampleDAO) FindByRTime(ctx context.Context, rTime time.Time, options ...QueryOption) ([]*SampleData, error) {
	option := resoveOption(options)
	query := fmt.Sprintf("select %s from %s where `r_time` = ? %s %s", dao.allColumns, dao.table, option.order, option.sort)

	args := make([]any, 0, 1)
	args = append(args, rTime)

	rows, err := dao.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return scanRows(rows, newSampleData)
}

func (dao *baseSampleDAO) PageByRTime(ctx context.Context, rTime time.Time, pageIndex int, pageSize int, options ...QueryOption) (*SamplePageData, error) {
	option := resoveOption(options)
	var data []*SampleData
	var count int

	args := make([]any, 0, 1)
	args = append(args, rTime)

	err := parallelQuery(
		func() error {
			query := fmt.Sprintf("select %s from %s where `r_time` = ? %s %s limit %d, %d", dao.allColumns, dao.table, option.order, option.sort, (pageIndex-1)*pageSize, pageSize)

			rows, err := dao.db.QueryContext(ctx, query, args...)
			if err != nil {
				return err
			}
			data, err = scanRows(rows, newSampleData)
			return err
		},
		func() error {
			query := fmt.Sprintf("select count(0) from (select `id` from %s where `r_time` = ?) as a", dao.table)
			row := dao.db.QueryRowContext(ctx, query, args...)
			return row.Scan(&count)
		},
	)

	if err != nil {
		return nil, err
	}

	return &SamplePageData{
		Data:  data,
		Count: count,
	}, nil
}

func (dao *baseSampleDAO) FindInByRTime(ctx context.Context, rTimeList []time.Time, options ...QueryOption) ([]*SampleData, error) {
	option := resoveOption(options)

	containNull := false
	args := make([]any, len(rTimeList))
	for i, item := range rTimeList {
		args[i] = item
	}

	if len(rTimeList) == 0 {
		return []*SampleData{}, nil
	}

	query := fmt.Sprintf("select %s from %s where `r_time` in (?%s) %s %s %s", dao.allColumns, dao.table, strings.Repeat(",?", len(rTimeList)-1), ternary(containNull, "or `r_time` is NULL", ""), option.order, option.sort)
	rows, err := dao.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return scanRows(rows, newSampleData)
}

func (dao *baseSampleDAO) PageInByRTime(ctx context.Context, rTimeList []time.Time, pageIndex int, pageSize int, options ...QueryOption) (*SamplePageData, error) {
	option := resoveOption(options)
	var data []*SampleData
	var count int

	containNull := false
	args := make([]any, len(rTimeList))
	for i, item := range rTimeList {
		args[i] = item
	}

	if len(rTimeList) == 0 {
		return &SamplePageData{
			Data:  []*SampleData{},
			Count: 0,
		}, nil
	}

	err := parallelQuery(
		func() error {
			query := fmt.Sprintf("select %s from %s where r_time in (?%s) %s %s %s limit %d, %d", dao.allColumns, dao.table, strings.Repeat(",?", len(rTimeList)-1), ternary(containNull, "or `r_time` is NULL", ""), option.order, option.sort, (pageIndex-1)*pageSize, pageSize)
			rows, err := dao.db.QueryContext(ctx, query, args...)
			if err != nil {
				return err
			}
			data, err = scanRows(rows, newSampleData)
			return err
		},
		func() error {
			query := fmt.Sprintf("select count(0) from (select `r_time` from %s where r_time in (?%s) %s) as a", dao.table, strings.Repeat(",?", len(rTimeList)-1), ternary(containNull, "or `r_time` is NULL", ""))
			row := dao.db.QueryRowContext(ctx, query, args...)
			return row.Scan(&count)
		},
	)

	if err != nil {
		return nil, err
	}

	return &SamplePageData{
		Data:  data,
		Count: count,
	}, nil
}

func (dao *baseSampleDAO) FindByUnion1AndUnion3(ctx context.Context, union1 string, union3 null.String, options ...QueryOption) ([]*SampleData, error) {
	option := resoveOption(options)
	query := fmt.Sprintf("select %s from %s where `union1` = ? and `union3` %s %s %s", dao.allColumns, dao.table, ternary(union3.Valid, "= ?", "is NULL"), option.order, option.sort)

	args := make([]any, 0, 2)
	args = append(args, union1)

	if union3.Valid {
		args = append(args, union3)
	}

	rows, err := dao.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return scanRows(rows, newSampleData)
}

func (dao *baseSampleDAO) PageByUnion1AndUnion3(ctx context.Context, union1 string, union3 null.String, pageIndex int, pageSize int, options ...QueryOption) (*SamplePageData, error) {
	option := resoveOption(options)
	var data []*SampleData
	var count int

	args := make([]any, 0, 2)
	args = append(args, union1)

	if union3.Valid {
		args = append(args, union3)
	}

	err := parallelQuery(
		func() error {
			query := fmt.Sprintf("select %s from %s where `union1` = ? and `union3` %s %s %s limit %d, %d", dao.allColumns, dao.table, ternary(union3.Valid, "= ?", "is NULL"), option.order, option.sort, (pageIndex-1)*pageSize, pageSize)

			rows, err := dao.db.QueryContext(ctx, query, args...)
			if err != nil {
				return err
			}
			data, err = scanRows(rows, newSampleData)
			return err
		},
		func() error {
			query := fmt.Sprintf("select count(0) from (select `id` from %s where `union1` = ? and `union3` %s) as a", dao.table, ternary(union3.Valid, "= ?", "is NULL"))
			row := dao.db.QueryRowContext(ctx, query, args...)
			return row.Scan(&count)
		},
	)

	if err != nil {
		return nil, err
	}

	return &SamplePageData{
		Data:  data,
		Count: count,
	}, nil
}

func (dao *baseSampleDAO) GetById(ctx context.Context, id int64) (*SampleData, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", dao.allColumns, dao.table)

	args := make([]any, 0, 1)
	args = append(args, id)

	row := dao.db.QueryRowContext(ctx, query, args...)
	data, err := scanRow(row, newSampleData)
	switch err {
	case nil:
		return data, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func (dao *baseSampleDAO) ExistById(ctx context.Context, id int64) (bool, error) {
	query := fmt.Sprintf("select `id` from %s where `id` = ? limit 1", dao.table)

	args := make([]any, 0, 1)
	args = append(args, id)

	row := dao.db.QueryRowContext(ctx, query, args...)
	var value int64
	err := row.Scan(&value)
	switch err {
	case nil:
		return true, nil
	case sql.ErrNoRows:
		return false, nil
	default:
		return false, err
	}
}

func (dao *baseSampleDAO) FindInById(ctx context.Context, idList []int64, options ...QueryOption) ([]*SampleData, error) {
	option := resoveOption(options)

	containNull := false
	args := make([]any, len(idList))
	for i, item := range idList {
		args[i] = item
	}

	if len(idList) == 0 {
		return []*SampleData{}, nil
	}

	query := fmt.Sprintf("select %s from %s where `id` in (?%s) %s %s %s", dao.allColumns, dao.table, strings.Repeat(",?", len(idList)-1), ternary(containNull, "or `id` is NULL", ""), option.order, option.sort)
	rows, err := dao.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return scanRows(rows, newSampleData)
}

func (dao *baseSampleDAO) PageInById(ctx context.Context, idList []int64, pageIndex int, pageSize int, options ...QueryOption) (*SamplePageData, error) {
	option := resoveOption(options)
	var data []*SampleData
	var count int

	containNull := false
	args := make([]any, len(idList))
	for i, item := range idList {
		args[i] = item
	}

	if len(idList) == 0 {
		return &SamplePageData{
			Data:  []*SampleData{},
			Count: 0,
		}, nil
	}

	err := parallelQuery(
		func() error {
			query := fmt.Sprintf("select %s from %s where id in (?%s) %s %s %s limit %d, %d", dao.allColumns, dao.table, strings.Repeat(",?", len(idList)-1), ternary(containNull, "or `id` is NULL", ""), option.order, option.sort, (pageIndex-1)*pageSize, pageSize)
			rows, err := dao.db.QueryContext(ctx, query, args...)
			if err != nil {
				return err
			}
			data, err = scanRows(rows, newSampleData)
			return err
		},
		func() error {
			query := fmt.Sprintf("select count(0) from (select `id` from %s where id in (?%s) %s) as a", dao.table, strings.Repeat(",?", len(idList)-1), ternary(containNull, "or `id` is NULL", ""))
			row := dao.db.QueryRowContext(ctx, query, args...)
			return row.Scan(&count)
		},
	)

	if err != nil {
		return nil, err
	}

	return &SamplePageData{
		Data:  data,
		Count: count,
	}, nil
}

func (dao *baseSampleDAO) GetByNInt(ctx context.Context, nInt null.Int) (*SampleData, error) {
	query := fmt.Sprintf("select %s from %s where `n_int` %s limit 1", dao.allColumns, dao.table, ternary(nInt.Valid, "= ?", "is NULL"))

	args := make([]any, 0, 1)
	if nInt.Valid {
		args = append(args, nInt)
	}

	row := dao.db.QueryRowContext(ctx, query, args...)
	data, err := scanRow(row, newSampleData)
	switch err {
	case nil:
		return data, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func (dao *baseSampleDAO) ExistByNInt(ctx context.Context, nInt null.Int) (bool, error) {
	query := fmt.Sprintf("select `id` from %s where `n_int` %s limit 1", dao.table, ternary(nInt.Valid, "= ?", "is NULL"))

	args := make([]any, 0, 1)
	if nInt.Valid {
		args = append(args, nInt)
	}

	row := dao.db.QueryRowContext(ctx, query, args...)
	var value int64
	err := row.Scan(&value)
	switch err {
	case nil:
		return true, nil
	case sql.ErrNoRows:
		return false, nil
	default:
		return false, err
	}
}

func (dao *baseSampleDAO) FindInByNInt(ctx context.Context, nIntList []null.Int, options ...QueryOption) ([]*SampleData, error) {
	option := resoveOption(options)

	containNull := false
	args := make([]any, len(nIntList))
	for i, item := range nIntList {
		if item.Valid {
			args[i] = nIntList[i]
		} else {
			containNull = true
		}
	}

	if len(nIntList) == 0 {
		return []*SampleData{}, nil
	}

	query := fmt.Sprintf("select %s from %s where `n_int` in (?%s) %s %s %s", dao.allColumns, dao.table, strings.Repeat(",?", len(nIntList)-1), ternary(containNull, "or `n_int` is NULL", ""), option.order, option.sort)
	rows, err := dao.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return scanRows(rows, newSampleData)
}

func (dao *baseSampleDAO) PageInByNInt(ctx context.Context, nIntList []null.Int, pageIndex int, pageSize int, options ...QueryOption) (*SamplePageData, error) {
	option := resoveOption(options)
	var data []*SampleData
	var count int

	containNull := false
	args := make([]any, len(nIntList))
	for i, item := range nIntList {
		if item.Valid {
			args[i] = nIntList[i]
		} else {
			containNull = true
		}
	}

	if len(nIntList) == 0 {
		return &SamplePageData{
			Data:  []*SampleData{},
			Count: 0,
		}, nil
	}

	err := parallelQuery(
		func() error {
			query := fmt.Sprintf("select %s from %s where n_int in (?%s) %s %s %s limit %d, %d", dao.allColumns, dao.table, strings.Repeat(",?", len(nIntList)-1), ternary(containNull, "or `n_int` is NULL", ""), option.order, option.sort, (pageIndex-1)*pageSize, pageSize)
			rows, err := dao.db.QueryContext(ctx, query, args...)
			if err != nil {
				return err
			}
			data, err = scanRows(rows, newSampleData)
			return err
		},
		func() error {
			query := fmt.Sprintf("select count(0) from (select `n_int` from %s where n_int in (?%s) %s) as a", dao.table, strings.Repeat(",?", len(nIntList)-1), ternary(containNull, "or `n_int` is NULL", ""))
			row := dao.db.QueryRowContext(ctx, query, args...)
			return row.Scan(&count)
		},
	)

	if err != nil {
		return nil, err
	}

	return &SamplePageData{
		Data:  data,
		Count: count,
	}, nil
}

func (dao *baseSampleDAO) GetByRInt(ctx context.Context, rInt int64) (*SampleData, error) {
	query := fmt.Sprintf("select %s from %s where `r_int` = ? limit 1", dao.allColumns, dao.table)

	args := make([]any, 0, 1)
	args = append(args, rInt)

	row := dao.db.QueryRowContext(ctx, query, args...)
	data, err := scanRow(row, newSampleData)
	switch err {
	case nil:
		return data, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func (dao *baseSampleDAO) ExistByRInt(ctx context.Context, rInt int64) (bool, error) {
	query := fmt.Sprintf("select `id` from %s where `r_int` = ? limit 1", dao.table)

	args := make([]any, 0, 1)
	args = append(args, rInt)

	row := dao.db.QueryRowContext(ctx, query, args...)
	var value int64
	err := row.Scan(&value)
	switch err {
	case nil:
		return true, nil
	case sql.ErrNoRows:
		return false, nil
	default:
		return false, err
	}
}

func (dao *baseSampleDAO) FindInByRInt(ctx context.Context, rIntList []int64, options ...QueryOption) ([]*SampleData, error) {
	option := resoveOption(options)

	containNull := false
	args := make([]any, len(rIntList))
	for i, item := range rIntList {
		args[i] = item
	}

	if len(rIntList) == 0 {
		return []*SampleData{}, nil
	}

	query := fmt.Sprintf("select %s from %s where `r_int` in (?%s) %s %s %s", dao.allColumns, dao.table, strings.Repeat(",?", len(rIntList)-1), ternary(containNull, "or `r_int` is NULL", ""), option.order, option.sort)
	rows, err := dao.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return scanRows(rows, newSampleData)
}

func (dao *baseSampleDAO) PageInByRInt(ctx context.Context, rIntList []int64, pageIndex int, pageSize int, options ...QueryOption) (*SamplePageData, error) {
	option := resoveOption(options)
	var data []*SampleData
	var count int

	containNull := false
	args := make([]any, len(rIntList))
	for i, item := range rIntList {
		args[i] = item
	}

	if len(rIntList) == 0 {
		return &SamplePageData{
			Data:  []*SampleData{},
			Count: 0,
		}, nil
	}

	err := parallelQuery(
		func() error {
			query := fmt.Sprintf("select %s from %s where r_int in (?%s) %s %s %s limit %d, %d", dao.allColumns, dao.table, strings.Repeat(",?", len(rIntList)-1), ternary(containNull, "or `r_int` is NULL", ""), option.order, option.sort, (pageIndex-1)*pageSize, pageSize)
			rows, err := dao.db.QueryContext(ctx, query, args...)
			if err != nil {
				return err
			}
			data, err = scanRows(rows, newSampleData)
			return err
		},
		func() error {
			query := fmt.Sprintf("select count(0) from (select `r_int` from %s where r_int in (?%s) %s) as a", dao.table, strings.Repeat(",?", len(rIntList)-1), ternary(containNull, "or `r_int` is NULL", ""))
			row := dao.db.QueryRowContext(ctx, query, args...)
			return row.Scan(&count)
		},
	)

	if err != nil {
		return nil, err
	}

	return &SamplePageData{
		Data:  data,
		Count: count,
	}, nil
}

func (dao *baseSampleDAO) GetByUnion1AndUnion2(ctx context.Context, union1 string, union2 string) (*SampleData, error) {
	query := fmt.Sprintf("select %s from %s where `union1` = ? and `union2` = ? limit 1", dao.allColumns, dao.table)

	args := make([]any, 0, 2)
	args = append(args, union1)

	args = append(args, union2)

	row := dao.db.QueryRowContext(ctx, query, args...)
	data, err := scanRow(row, newSampleData)
	switch err {
	case nil:
		return data, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func (dao *baseSampleDAO) ExistByUnion1AndUnion2(ctx context.Context, union1 string, union2 string) (bool, error) {
	query := fmt.Sprintf("select `id` from %s where `union1` = ? and `union2` = ? limit 1", dao.table)

	args := make([]any, 0, 2)
	args = append(args, union1)

	args = append(args, union2)

	row := dao.db.QueryRowContext(ctx, query, args...)
	var value int64
	err := row.Scan(&value)
	switch err {
	case nil:
		return true, nil
	case sql.ErrNoRows:
		return false, nil
	default:
		return false, err
	}
}

func (dao *baseSampleDAO) Insert(ctx context.Context, data *SampleData) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (`gmt_create`,`gmt_modified`,`r_int`,`n_int`,`r_float`,`n_float`,`r_string`,`n_string`,`r_time`,`n_time`,`union1`,`union2`,`union3`) values (NOW(),NOW(),?,?,?,?,?,?,?,?,?,?,?)", dao.table)
	result, err := dao.db.ExecContext(ctx, query, data.RInt, data.NInt, data.RFloat, data.NFloat, data.RString, data.NString, data.RTime, data.NTime, data.Union1, data.Union2, data.Union3)
	return result, err
}

func (dao *baseSampleDAO) Update(ctx context.Context, id int64, data *SampleData) error {
	query := fmt.Sprintf("update %s set `gmt_modified` = NOW(),`r_int` = ?,`n_int` = ?,`r_float` = ?,`n_float` = ?,`r_string` = ?,`n_string` = ?,`r_time` = ?,`n_time` = ?,`union1` = ?,`union2` = ?,`union3` = ? where `id` = ?", dao.table)
	_, err := dao.db.ExecContext(ctx, query, data.RInt, data.NInt, data.RFloat, data.NFloat, data.RString, data.NString, data.RTime, data.NTime, data.Union1, data.Union2, data.Union3, id)
	return err
}

func (dao *baseSampleDAO) UpdatePartial(ctx context.Context, id int64, data *SamplePartialData) error {
	sets := make([]string, 0, 14)
	args := make([]any, 0, 14)

	sets = append(sets, "`gmt_modified` = NOW()")

	if v, present := data.RInt.Get(); present {
		sets = append(sets, "`r_int` = ?")
		args = append(args, v)
	}

	if v, present := data.NInt.Get(); present {
		sets = append(sets, "`n_int` = ?")
		args = append(args, v)
	}

	if v, present := data.RFloat.Get(); present {
		sets = append(sets, "`r_float` = ?")
		args = append(args, v)
	}

	if v, present := data.NFloat.Get(); present {
		sets = append(sets, "`n_float` = ?")
		args = append(args, v)
	}

	if v, present := data.RString.Get(); present {
		sets = append(sets, "`r_string` = ?")
		args = append(args, v)
	}

	if v, present := data.NString.Get(); present {
		sets = append(sets, "`n_string` = ?")
		args = append(args, v)
	}

	if v, present := data.RTime.Get(); present {
		sets = append(sets, "`r_time` = ?")
		args = append(args, v)
	}

	if v, present := data.NTime.Get(); present {
		sets = append(sets, "`n_time` = ?")
		args = append(args, v)
	}

	if v, present := data.Union1.Get(); present {
		sets = append(sets, "`union1` = ?")
		args = append(args, v)
	}

	if v, present := data.Union2.Get(); present {
		sets = append(sets, "`union2` = ?")
		args = append(args, v)
	}

	if v, present := data.Union3.Get(); present {
		sets = append(sets, "`union3` = ?")
		args = append(args, v)
	}

	args = append(args, id)

	query := fmt.Sprintf("update %s set %s where `id` = ?", dao.table, strings.Join(sets, ","))
	_, err := dao.db.ExecContext(ctx, query, args...)
	return err
}

func (dao *baseSampleDAO) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", dao.table)
	_, err := dao.db.ExecContext(ctx, query, id)
	return err
}
