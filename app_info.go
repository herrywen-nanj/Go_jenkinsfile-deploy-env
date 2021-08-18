package main

import (
    "fmt"
    "flag"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

const (
    mysql_user   = "root"
    mysql_passwd = "aqn2AO3N8HDrHMFG"
    mysql_ip     = "mysql-svc.test-midd"
    mysql_port   = "3306"
    mysql_db     = "jenkins"
)

var name string
var group string
var key string
var value string

func main() {


    flag.StringVar(&name, "name", "", "应用名称")
    flag.StringVar(&group, "group", "", "应用项目组")
    flag.StringVar(&key, "key", "", "所需查询的key")
    flag.Parse()

    db, err := sql.Open("mysql", mysql_user+":"+mysql_passwd+"@tcp("+mysql_ip+":"+mysql_port+")/"+mysql_db+"?charset=utf8")
    if err != nil {
        //fmt.Println("连接数据库失败:",err)
    }
    sql := "SELECT " + key + " FROM app_info WHERE name=? and project_group=?"
    err = db.QueryRow(sql,name,group).Scan(&value)
    if err != nil {
        //fmt.Println(err)
    }

    if len(value) == 0 {
        fmt.Println("data not found")
    }else{
        fmt.Println(value)
    }

    defer db.Close()
}
