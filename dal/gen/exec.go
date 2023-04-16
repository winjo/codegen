package gen

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

func (g *generator) genExec() string {
	autoSetFields := []string{"id", "gmt_create", "gmt_modified"}
	filterFields := lo.Filter(g.Fields, func(f *Field, i int) bool {
		return !lo.Contains(autoSetFields, f.Name)
	})
	execArgs := lo.Map(filterFields, func(f *Field, i int) string {
		return fmt.Sprintf("data.%s", f.PascalCaseName)
	})

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

	insertColumns := make([]string, 0)
	insertValues := make([]string, 0)
	insertArgs := make([]string, 0)
	for _, f := range g.Fields {
		if f == g.PrimaryField {
			if !strings.Contains(g.PrimaryField.Extra, "auto_increment") {
				insertColumns = append(insertColumns, fmt.Sprintf("`%s`", f.Name))
				insertValues = append(insertValues, "?")
				insertArgs = append(insertArgs, fmt.Sprintf("data.%s", f.PascalCaseName))
			}

		} else {
			insertColumns = append(insertColumns, fmt.Sprintf("`%s`", f.Name))
			if f == createField || f == modifiedField {
				insertValues = append(insertValues, getTimeFieldValue(f))
			} else {
				insertValues = append(insertValues, "?")
				insertArgs = append(insertArgs, fmt.Sprintf("data.%s", f.PascalCaseName))
			}
		}
	}

	updateSets := make([]string, 0)
	updateArgs := make([]string, 0)
	for _, f := range g.Fields {
		if f == g.PrimaryField || f == createField {
			continue
		}
		if f == modifiedField {
			updateSets = append(updateSets, fmt.Sprintf("`%s` = %s", f.Name, getTimeFieldValue(f)))
		} else {
			updateSets = append(updateSets, fmt.Sprintf("`%s` = ?", f.Name))
			updateArgs = append(updateArgs, fmt.Sprintf("data.%s", f.PascalCaseName))
		}
	}

	return tmpl.Exec("exec", map[string]any{
		"tableObject":          g.Object.Name,
		"camelCaseObject":      g.Object.CamelCaseName,
		"pascalCaseObject":     g.Object.PascalCaseName,
		"primaryField":         g.PrimaryField,
		"primaryAutoIncrement": strings.Contains(g.PrimaryField.Extra, "auto_increment"),
		"createField":          createField,
		"modifiedField":        modifiedField,
		"insertColumns":        strings.Join(insertColumns, ","),
		"insertValues":         strings.Join(insertValues, ","),
		"insertArgs":           strings.Join(insertArgs, ","),
		"updateSets":           strings.Join(updateSets, ","),
		"updateArgs":           strings.Join(updateArgs, ","),
		"fields":               g.Fields,
		"execArgs":             strings.Join(execArgs, ","),
	})
}

func getTimeFieldValue(f *Field) string {
	if f.DataType == "int" { // 秒级时间戳
		return "UNIX_TIMESTAMP()"
	} else if f.DataType == "bigint" { // 毫秒级暑假戳
		return "ROUND(UNIX_TIMESTAMP(CURTIME(4) * 1000))"
	} else { // 等同 CURRENT_TIMESTAMPE
		length := ""
		if f.Length != nil {
			length = strconv.Itoa(*f.Length)
		}
		return fmt.Sprintf("NOW(%s)", length)
	}
}
