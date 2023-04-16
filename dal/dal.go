package dal

import (
	"fmt"
	"os"

	"github.com/winjo/codegen/base"
	"github.com/winjo/codegen/dal/gen"
	"github.com/winjo/codegen/dal/schema"
	"github.com/winjo/codegen/dal/util"
)

var CmdDal = &base.Command{
	Name: "dal",
}

var (
	datasource = CmdDal.Flag.String("ds", "", "mysql datasource")
	withTest   = CmdDal.Flag.Bool("test", false, "gen test code")
)

func init() {
	CmdDal.Run = run
}

func run(cmd *base.Command, args []string) {
	if *datasource == "" {
		fmt.Printf("Run codegen %s -h for usage.\n", cmd.Flag.Name())
		os.Exit(2)
	}

	im := schema.NewInformationSchemaModel(*datasource)
	tables := im.GetAllTables()

	util.EnsureDir("dal/structure")
	util.EnsureDir("dal/dao")

	schema.GenStructure(tables)

	gen.GenCode(tables, *withTest)
}
