# codegen
> 代码生成

# 使用
## 安装
Go >= 1.18

```bash
go install github.com/winjo/codegen
```

## 数据库 dal 代码生成
```bash
codegen dal -ds '<user>:<password>@tcp(127.0.0.1:3306)/<db>?charset=utf8mb4'
```

自动拉取库所有表，并根据索引生成代码，生成方法包含
- Find 查询数据
- Page 分页查询
- Count 查询记录数
- GetBy 唯一索引查询
- ExistBy 唯一索引查询数据是否存在
- FindBy 普通索引查询
- PageBy 普通索引分页查询
- FindInBy 唯一索引 in 查询
- PageInBy 唯一索引分页 in 查询
- Insert 插入数据
- Update 更新数据
- UpdatePartial 更新部分数据
- Delete 删除数据

具体代码参考 [examples](./examples/dal/dao)

### 约定
- 创建时间字段使用 gmt_create
- 修改时间字段使用 gmt_modified
- dsn 增加配置 parseTime=true
- 默认 go-sql-driver/mysql 使用的是 prepare，如果 server 不支持，可使用本地插值的方式，dsn 配置 interpolateParams=true
