package main

import (
    "database/sql"
    "fmt"
    "io/ioutil"
    "log"
    "strings"

    _ "github.com/go-sql-driver/mysql"
    "github.com/spf13/viper"
)

type Table struct {
    Name   string
    Fields []Field
}

type Field struct {
    Name   string
    Type   string
}

func main() {
    // 读取配置文件(read config file)
    viper.SetConfigName("config")
    viper.AddConfigPath(".")
    err := viper.ReadInConfig()
    if err != nil {
        log.Fatal(err)
    }

    // 连接数据库1(connect to db1)
    db1, err := sql.Open(viper.GetString("database1.driver"), viper.GetString("database1.dsn"))
    if err != nil {
        log.Fatal(err)
    }
    defer db1.Close()

    // 连接数据库2(connect to db2)
    db2, err := sql.Open(viper.GetString("database2.driver"), viper.GetString("database2.dsn"))
    if err != nil {
        log.Fatal(err)
    }
    defer db2.Close()

    // 获取数据库1中所有表的结构(get table schema from db1)
    tables1, err := getTables(db1)
    if err != nil {
        log.Fatal(err)
    }

    // 获取数据库2中所有表的结构(get table schema from db2)
    tables2, err := getTables(db2)
    if err != nil {
        log.Fatal(err)
    }

    // 对比两个数据库中所有表的结构(diff table schema)
    var diffTables []string
    for _, table1 := range tables1 {
        table2, ok := tables2[table1.Name]
        if !ok {
            diffTables = append(diffTables, table1.Name)
            continue
        }

        if !compareFields(table1.Fields, table2.Fields) {
            diffTables = append(diffTables, table1.Name)
        }
    }

    // 输出不一致的表到文件中
    if len(diffTables) > 0 {
        output := strings.Join(diffTables, "\n")
        err = ioutil.WriteFile(viper.GetString("output"), []byte(output), 0644)
        if err != nil {
            log.Fatal(err)
        }
    }
}

// 获取数据库中所有表的结构
func getTables(db *sql.DB) (map[string]Table, error) {
    tables := make(map[string]Table)

    rows, err := db.Query("SHOW TABLES")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var tableName string
        if err := rows.Scan(&tableName); err != nil {
            return nil, err
        }

        table := Table{Name: tableName}

        fields, err := getFields(db, tableName)
        if err != nil {
            return nil, err
        }
        table.Fields = fields

        tables[tableName] = table
    }

    return tables, nil
}

// 获取表中所有字段的结构(get table columns)
func getFields(db *sql.DB, tableName string) ([]Field, error) {
    fields := make([]Field, 0)

    rows, err := db.Query(fmt.Sprintf("SHOW COLUMNS FROM %s", tableName))
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var fieldName, fieldType string
        var fieldNull string
        var key,defaultVal,extra interface{}

        if err := rows.Scan(&fieldName, &fieldType, &fieldNull, &key, &defaultVal, &extra); err != nil {
            return nil, err
        }

        field := Field{Name: fieldName, Type: fieldType}

        fields = append(fields, field)
    }

    return fields, nil
}

// 对比两个表的字段结构是否一致(diff table schema)
func compareFields(fields1, fields2 []Field) bool {
    if len(fields1) != len(fields2) {
        return false
    }

    for i := range fields1 {
        if fields1[i].Name != fields2[i].Name ||
            fields1[i].Type != fields2[i].Type {
            return false
        }
    }

    return true
}
