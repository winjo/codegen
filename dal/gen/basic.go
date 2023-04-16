package gen

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

func (g *generator) genBasic() string {
	return tmpl.Exec("basic", map[string]any{
		"tableObject":      g.Object.Name,
		"camelCaseObject":  g.Object.CamelCaseName,
		"pascalCaseObject": g.Object.PascalCaseName,
		"primaryField":     g.PrimaryField,
		"fieldNames": lo.Map(g.Fields, func(f *Field, i int) string {
			return f.PascalCaseName
		}),
		"fieldsSlice": strings.Join(lo.Map(g.Fields, func(f *Field, i int) string {
			return fmt.Sprintf("\"%s\"", f.Name)
		}), ","),
		"fieldsStruct": strings.Join(lo.Map(g.Fields, func(f *Field, i int) string {
			return fmt.Sprintf("%s %s `db:\"%s\" json:\"%s\"` %s", f.PascalCaseName, f.DataType, f.Name, f.CamelCaseName, lo.Ternary(f.Comment == "", "", fmt.Sprintf("// %s", f.Comment)))
		}), "\n"),
		"fieldsPartialStruct": strings.Join(lo.Map(g.Fields, func(f *Field, i int) string {
			if lo.Contains([]string{g.PrimaryField.Name, CreateField, ModifiedField}, f.Name) {
				return ""
			}
			return fmt.Sprintf("%s mo.Option[%s] `db:\"%s\" json:\"%s\"` %s", f.PascalCaseName, f.DataType, f.Name, f.CamelCaseName, lo.Ternary(f.Comment == "", "", fmt.Sprintf("// %s", f.Comment)))
		}), "\n"),
	})
}
