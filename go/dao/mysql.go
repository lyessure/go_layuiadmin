package dao

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var Db *sqlx.DB

var (
	mysqluser     = "root"
	mysqlpass     = "123456"
	mysqlhost     = "localhost"
	mysqldb       = "admin"
	maxthread     = 20
	maxidlethread = 10
)

func Initmysql() {
	var err error

	if Db != nil {
		return
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True", mysqluser, mysqlpass, mysqlhost, mysqldb)
	Db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		panic("db init fail")
	}
	Db.SetMaxOpenConns(maxthread)
	Db.SetMaxIdleConns(maxidlethread)
	Db.SetConnMaxLifetime(10 * time.Minute)
}
