package mysql

import (
	"database/sql"
	"fmt"
	"os"
	_ "github.com/go-sql-driver/mysql"
	"log"
	cfg "filestore-hsz/service/dbproxy/config"
)

var db *sql.DB

func init() {
	log.Println("db conn address" + cfg.MySQLSource)
	db, _ = sql.Open("mysql", cfg.MySQLSource)
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		fmt.Println("Failed to connect to mysql, err" + err.Error())
		os.Exit(1)
	}
}

// 返回数据连接对象
func DBConn() *sql.DB {
	return db
}

func ParseRows(rows *sql.Rows) []map[string]interface{} {
	// 获取记录列（名）
	columns, _ := rows.Columns()
	// 创建列值的slice（values），并为每一列初始化一个指针
	// scanArgs用作rows.Scan中的传入参数
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	// record为每次迭代中存储记录的临时变量
	record := make(map[string]interface{})
	// records为函数最终返回的数据（列表）
	records := make([]map[string]interface{}, 0)
	// 迭代记录
	for rows.Next() {
		//每Scan一次，将一行数据保存到record字典
		err := rows.Scan(scanArgs...)
		checkErr(err)

		for i, col := range values {
			if col != nil {
				record[columns[i]] = col
			}
		}
		records = append(records, record)
	}
	return records
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

