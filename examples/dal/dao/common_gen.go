package dao

import (
	"database/sql"
	"github.com/samber/mo"
	"strings"
	"sync/atomic"
)

type QueryOption struct {
	OrderBy mo.Option[string]
	Desc    mo.Option[bool]
}

type queryOption struct {
	order string
	sort  string
}

func resoveOption(options []QueryOption) queryOption {
	defaultOption := queryOption{
		order: "order by id",
		sort:  "",
	}
	switch len(options) {
	case 0:
	case 1:
		option := options[0]
		if v, present := option.OrderBy.Get(); present {
			defaultOption.order = "order by " + v
		}
		if v, present := option.Desc.Get(); present && v {
			defaultOption.sort = "desc"
		}
	default:
		panic("two many options parameters")
	}
	return defaultOption
}

func columnsToRow(columns []string) string {
	wrapNames := make([]string, len(columns))

	for i, c := range columns {
		wrapNames[i] = "`" + c + "`"
	}

	return strings.Join(wrapNames, ",")
}

func omitColumns(columns []string, excludeColumns ...string) []string {
	ret := make([]string, 0, len(columns))

	m := make(map[string]struct{})
	for _, c := range excludeColumns {
		m[c] = struct{}{}
	}
	for _, c := range columns {
		if _, ok := m[c]; !ok {
			ret = append(ret, c)
		}
	}
	return ret
}

func scanRow[T any](row *sql.Row, newObject func() (T, []any)) (T, error) {
	o, dest := newObject()
	err := row.Scan(dest...)
	return o, err
}

func scanRows[T any](rows *sql.Rows, newObject func() (T, []any)) ([]T, error) {
	defer rows.Close()
	var items []T
	for rows.Next() {
		o, dest := newObject()
		err := rows.Scan(dest...)
		if err != nil {
			return nil, err
		}
		items = append(items, o)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func parallelQuery(fns ...func() error) (err error) {
	done := make(chan struct{})
	num := int32(len(fns))
	cancelled := int32(0)

	for _, f := range fns {
		go func(f func() error) {
			if e := f(); e != nil {
				if atomic.CompareAndSwapInt32(&cancelled, 0, 1) {
					err = e
					close(done)
				}
			} else {
				if atomic.AddInt32(&num, -1) == 0 {
					close(done)
				}
			}
		}(f)
	}

	<-done

	return
}

func ternary[T any](condition bool, left T, right T) T {
	if condition {
		return left
	}
	return right
}
