package dbops

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

const (
	USERNAME = "root"
	PASSWORD = "root"
	NETWORK  = "tcp"
	SERVER   = "localhost"
	PORT     = 3306
	DATABASE = "go_video"
)

var (
	dbConn *sql.DB
	err    error
)

func init() {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	dbConn, err = sql.Open("mysql", dsn)
	//dbConn, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/go_video?charset=utf8")
	if err != nil {
		panic(err.Error())
	}
}
