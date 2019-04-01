package dbops

import (
	_ "github.com/go-sql-driver/mysql"
	"go_log/collect_log/defs"
)

// 增加log记录
func AddLog(log *defs.Log) error  {
	stmtIns, err := dbConn.Prepare("INSERT INTO log (url, url_id, type, ip) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmtIns.Exec(log.Url, log.UrlId, log.Type, log.Type)
	if err != nil {
		return err
	}

	defer stmtIns.Close()
	return nil

}
