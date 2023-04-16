package gen

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

func (g *generator) genTest() string {
	createField, existCreateField := lo.Find(g.Fields, func(f *Field) bool {
		return f.Name == CreateField
	})
	if !existCreateField {
		panic(fmt.Sprintf("miss %s in %s", CreateField, g.Object.Name))
	}

	modifiedField, existModifiedField := lo.Find(g.Fields, func(f *Field) bool {
		return f.Name == ModifiedField
	})
	if !existModifiedField {
		panic(fmt.Sprintf("miss %s in %s", ModifiedField, g.Object.Name))
	}

	filterString := func(list []string) []string {
		return lo.Filter(list, func(item string, i int) bool {
			return item != ""
		})
	}

	return tmpl.Exec("test", map[string]any{
		"tableObject":      g.Object.Name,
		"camelCaseObject":  g.Object.CamelCaseName,
		"pascalCaseObject": g.Object.PascalCaseName,
		"primaryField":     g.PrimaryField,
		"createField":      createField,
		"modifiedField":    modifiedField,
		"fieldsData": lo.Map([]int{0, 1}, func(i int, index int) string {
			return strings.Join(filterString(lo.Map(g.Fields, func(f *Field, index int) string {
				if lo.Contains([]string{g.PrimaryField.Name, CreateField, ModifiedField}, f.Name) {
					return ""
				}
				var value string
				switch f.DataType {
				case "int64":
					value = fmt.Sprintf("%d", i)
				case "null.Int":
					value = "null.Int{}"
				case "float64":
					value = fmt.Sprintf("%d", i)
				case "null.Float":
					value = "null.Float{}"
				case "bool":
					value = fmt.Sprintf("%t", i > 0)
				case "null.Bool":
					value = "null.Bool{}"
				case "string":
					value = fmt.Sprintf(`"%s%d"`, string([]rune(f.CamelCaseName)[0]), i)
				case "null.String":
					value = "null.String{}"
				case "time.Time":
					value = fmt.Sprintf("time.Unix(time.Now().Unix() + int64(%d), 0).UTC()", i)
				case "null.Time":
					value = "null.Time{}"
				default:
					value = "nil"
				}
				return fmt.Sprintf(`%s: %s,`, f.PascalCaseName, value)
			})), "\n")
		}),
		"fieldsNoData": strings.Join(filterString(lo.Map(g.Fields, func(f *Field, index int) string {
			if lo.Contains([]string{g.PrimaryField.Name, CreateField, ModifiedField}, f.Name) {
				return ""
			}
			var value string
			switch f.DataType {
			case "int64":
				value = "404"
			case "null.Int":
				value = "null.IntFrom(404)"
			case "float64":
				value = "404"
			case "null.Float":
				value = "null.FloatFrom(404)"
			case "string":
				value = fmt.Sprintf(`"%s%d"`, string([]rune(f.CamelCaseName)[0]), 404)
			case "null.String":
				value = fmt.Sprintf(`null.StringFrom("%s%d")`, string([]rune(f.CamelCaseName)[0]), 404)
			case "time.Time":
				value = "time.Unix(int64(404), 0)"
			case "null.Time":
				value = "null.TimeFrom(time.Unix(int64(404), 0))"
			default:
				value = "nil"
			}
			return fmt.Sprintf(`%s: %s,`, f.PascalCaseName, value)
		})), "\n"),
		"indexList": lo.Map(g.Indexes, func(index *Index, i int) map[string]any {
			return map[string]any{
				"fields": index.Fields,
				"unique": index.Unique,
				"fieldsJoinAnd": strings.Join(lo.Map(index.Fields, func(f *Field, i int) string {
					return f.PascalCaseName
				}), "And"),
				"fieldsRealArgs": strings.Join(lo.Map(index.Fields, func(f *Field, i int) string {
					return fmt.Sprintf("data.%s", f.PascalCaseName)
				}), ","),
				"fieldsRealArgs404": strings.Join(lo.Map(index.Fields, func(f *Field, i int) string {
					return fmt.Sprintf("noData.%s", f.PascalCaseName)
				}), ","),
			}
		}),
	})
}
