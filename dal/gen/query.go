package gen

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

func (g *generator) genQuery() string {
	indexList := lo.Map(g.Indexes, func(index *Index, i int) map[string]any {
		return map[string]any{
			"fieldsJoinAnd": strings.Join(lo.Map(index.Fields, func(f *Field, i int) string {
				return f.PascalCaseName
			}), "And"),
			"fieldsRow": strings.Join(lo.Map(index.Fields, func(f *Field, i int) string {
				return fmt.Sprintf("`%s`", f.Name)
			}), ","),
			"fieldsFormalArgs": strings.Join(lo.Map(index.Fields, func(f *Field, i int) string {
				return fmt.Sprintf("%s %s", f.CamelCaseName, f.DataType)
			}), ","),
			"fieldsRealArgs": strings.Join(lo.Map(index.Fields, func(f *Field, i int) string {
				return f.CamelCaseName
			}), ","),
			"fieldsCondition": strings.Join(lo.Map(index.Fields, func(f *Field, i int) string {
				return fmt.Sprintf("`%s` %s", f.Name, lo.If(f.Nullable, "%s").Else("= ?"))
			}), " and "),
			"fieldsConditionArgs": lo.Reduce(index.Fields, func(str string, f *Field, i int) string {
				if f.Nullable {
					str += fmt.Sprintf(`,ternary(%s.Valid, "= ?", "is NULL")`, f.CamelCaseName)
				}
				return str
			}, ""),
			"fields": index.Fields,
			"unique": index.Unique,
		}
	})
	return tmpl.Exec("query", map[string]any{
		"tableObject":      g.Object.Name,
		"camelCaseObject":  g.Object.CamelCaseName,
		"pascalCaseObject": g.Object.PascalCaseName,
		"primaryField":     g.PrimaryField,
		"indexList":        indexList,
	})
}
