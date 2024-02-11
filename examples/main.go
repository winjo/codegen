package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/winjo/codegen/examples/dal/dao"
)

type dal struct {
	Sample *dao.SampleDAO
}

var DAL *dal

func main() {
	db, err := sql.Open("mysql", "root:abc123@tcp(127.0.0.1:33060)/codegen_test_examples?charset=utf8mb4&parseTime=true&loc=Local&interpolateParams=true")
	if err != nil {
		panic(fmt.Errorf("connect mysql failed %s", err))
	}
	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("ping failed %s", err))
	}
	conn, err := db.Conn(context.Background())
	if err != nil {
		panic(fmt.Errorf("ping failed %s", err))
	}
	// DAL = &dal{
	// 	Sample: dao.NewSampleDAO(db),
	// }
	sample := dao.NewSampleDAO(conn)
	data, err := sample.PageInById(context.Background(), []int64{1}, 1, 1)
	if err != nil {
		panic(fmt.Errorf("ping failed %s", err))
	}
	fmt.Println(data.Data[0])
	// tx.Commit()
}
