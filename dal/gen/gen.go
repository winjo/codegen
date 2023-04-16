package gen

import (
	"fmt"
	"go/format"
	"os"
	"strings"

	"github.com/samber/lo"
	"github.com/winjo/codegen/dal/schema"
	"github.com/winjo/codegen/dal/util"
)

var formatCode = true

type (
	Object struct {
		Name           string
		CamelCaseName  string
		PascalCaseName string
	}

	Field struct {
		*schema.Column
		CamelCaseName  string
		PascalCaseName string
		DataType       string
		AllowNull      bool
	}

	Index struct {
		*schema.Index
		Fields []*Field
	}

	generator struct {
		Object       *Object
		Fields       []*Field
		PrimaryField *Field
		Indexes      []*Index
	}
)

func GenCode(tables []*schema.Table, withTest bool) {
	genCommonCode()
	if withTest {
		genTestMainCode()
	}
	for _, table := range tables {
		g := newTableGenerator(table)
		g.genDAOBaseCode()
		g.genDAOCode()
		if withTest {
			g.genTestCode()
		}
	}
}

func newTableGenerator(table *schema.Table) *generator {
	object := &Object{
		Name:           table.Name,
		CamelCaseName:  util.ToCamelCase(table.Name),
		PascalCaseName: util.ToPascalCase(table.Name),
	}

	fields := make([]*Field, 0, len(table.Columns))
	fieldMap := make(map[string]*Field)
	for _, c := range table.Columns {
		f := &Field{
			Column:         c,
			CamelCaseName:  util.ToCamelCase(c.Name),
			PascalCaseName: util.ToPascalCase(c.Name),
			DataType:       convertType(c),
		}
		fields = append(fields, f)
		fieldMap[c.Name] = f
	}

	indexes := make([]*Index, 0, len(table.Indexes))
	var primayField *Field
	for _, index := range table.Indexes {
		indexFields := lo.Map(index.Columns, func(c string, i int) *Field {
			return fieldMap[c]
		})
		i := &Index{
			Index:  index,
			Fields: indexFields,
		}
		indexes = append(indexes, i)
		if index.Name == "PRIMARY" {
			if len(index.Columns) != 1 {
				panic("primary key length should be 1")
			}
			primayField = fieldMap[index.Columns[0]]
		}
	}

	return &generator{
		Object:       object,
		Fields:       fields,
		PrimaryField: primayField,
		Indexes:      indexes,
	}
}

func writeCode(path string, code string) {
	sourceCode := []byte(code)
	if formatCode {
		var err error
		sourceCode, err = format.Source(sourceCode)
		util.AssertNotNil(err)
	}
	os.WriteFile(path, sourceCode, os.ModePerm)
}

func genCommonCode() {
	code := tmpl.Exec("common", map[string]any{
		"primaryField": PrimaryField,
	})
	writeCode("dal/dao/common_gen.go", code)
}

func (g *generator) genDAOBaseCode() {
	code := strings.Join([]string{
		g.genBasic(),
		g.genQuery(),
		g.genExec(),
	}, "\n\n")
	writeCode(fmt.Sprintf("dal/dao/%s_dao_gen.go", g.Object.Name), code)
}

func (g *generator) genDAOCode() {
	code := tmpl.Exec("dao", map[string]any{
		"tableObject":      g.Object.Name,
		"camelCaseObject":  g.Object.CamelCaseName,
		"pascalCaseObject": g.Object.PascalCaseName,
	})
	writeCode(fmt.Sprintf("dal/dao/%s_dao.go", g.Object.Name), code)
}

func genTestMainCode() {
	code := tmpl.Exec("main_test", map[string]any{})
	writeCode("dal/dao/main_test.go", code)
}

func (g *generator) genTestCode() {
	code := g.genTest()
	writeCode(fmt.Sprintf("dal/dao/%s_test.go", g.Object.Name), code)
}
