.PHONY: gen-examples test-examples

gen-examples:
	cd examples && rm -rf ./dal && go run .. dal -ds 'root:abc123@tcp(127.0.0.1:33060)/codegen_test?charset=utf8mb4' --test

test-examples:
	MYSQL_HOST=127.0.0.1 MYSQL_PORT=33060 MYSQL_USER=root MYSQL_PASS=abc123 go test -v ./examples/dal/dao
