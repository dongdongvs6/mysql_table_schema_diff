# mysql_table_schema_diff(MySQL简单易用的表结构对比工具)
use simple way to diff mysql table schema, include table_name, column_name, column_type and column_length

# Use method(使用方法）
1. create one config.yaml file like: (创建１个config.yaml的配置文件)
`
database1:
  driver: mysql
  dsn: user:password@tcp(host:port)/database1

database2:
  driver: mysql
  dsn: user:password@tcp(host:port)/database2

output: diff_tables.txt
`
2. build　go file

go build table_diff.go

3. put the binary file and config file in the same path, run schema diff (将二进制与配置文件放在１个目录下，运行对比)

./table_diff

4. check the diff result in output file (在配置文件中的输出文件中查看不一致的表)
