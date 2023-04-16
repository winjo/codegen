package gen

import (
	"fmt"

	"github.com/winjo/codegen/dal/schema"
)

var typeMapString = map[string]string{
	// bool
	"bool":    "bool",
	"boolean": "bool",
	// number
	"tinyint":   "int64",
	"smallint":  "int64",
	"mediumint": "int64",
	"int":       "int64",
	"int1":      "int64",
	"int2":      "int64",
	"int3":      "int64",
	"int4":      "int64",
	"int8":      "int64",
	"integer":   "int64",
	"bigint":    "int64",
	"float":     "float64",
	"float4":    "float64",
	"float8":    "float64",
	"double":    "float64",
	"decimal":   "float64",
	"dec":       "float64",
	"fixed":     "float64",
	"real":      "float64",
	"bit":       "int64",
	// date & time
	"date":      "time.Time",
	"datetime":  "time.Time",
	"timestamp": "time.Time",
	"time":      "string",
	"year":      "int64",
	// string
	"linestring":      "string",
	"multilinestring": "string",
	"nvarchar":        "string",
	"nchar":           "string",
	"char":            "string",
	"character":       "string",
	"varchar":         "string",
	"binary":          "string",
	"bytea":           "string",
	"longvarbinary":   "string",
	"varbinary":       "string",
	"tinytext":        "string",
	"text":            "string",
	"mediumtext":      "string",
	"longtext":        "string",
	"enum":            "string",
	"set":             "string",
	"json":            "string",
	"jsonb":           "string",
	"blob":            "string",
	"longblob":        "string",
	"mediumblob":      "string",
	"tinyblob":        "string",
	"ltree":           "[]byte",
}

func convertType(c *schema.Column) string {
	goType, ok := typeMapString[c.Type]
	if !ok {
		panic(fmt.Sprintf("unsupport database type: %v", c.Type))
	}

	if c.Nullable {
		switch goType {
		case "int64":
			return "null.Int"
		case "float64":
			return "null.Float"
		case "bool":
			return "null.Bool"
		case "string":
			return "null.String"
		case "time.Time":
			return "null.Time"
		}
	}

	return goType
}
