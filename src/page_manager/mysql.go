
package kweb_manager

import (
	os "os"
	fmt "fmt"
	time "time"

	ufile "github.com/KpnmServer/go-util/file"
	json "github.com/KpnmServer/go-util/json"

	kpsql "github.com/KpnmServer/go-kpsql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	_SQL_HOST string // = "127.0.0.1"
	_SQL_PORT uint16 // = 3306
	_SQL_USER string // = "user"
	_SQL_PWD string // = "password"
	_SQL_BASE string // = "database"
	_SQL_CSET string // = "charset"
)

var SQLDB kpsql.SqlDatabase

func init(){
	var err error
	{// Read config
		var fd *os.File
		var err error
		fd, err = os.Open(ufile.JoinPath("config", "sql.json"))
		if err != nil {
			panic(err)
		}
		defer fd.Close()
		var obj = make(json.JsonObj)
		err = json.ReadJson(fd, &obj)
		if err != nil {
			panic(err)
		}
		_SQL_HOST = obj.GetString("host")
		_SQL_PORT = obj.GetUInt16("port")
		_SQL_USER = obj.GetString("user")
		_SQL_PWD = obj.GetString("password")
		_SQL_BASE = obj.GetString("database")
		_SQL_CSET = "utf8"
	}
	{// Connect SQL
		dbDSN := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true",
			_SQL_USER, _SQL_PWD, _SQL_HOST, _SQL_PORT, _SQL_BASE, _SQL_CSET)

		SQLDB, err = kpsql.Open("mysql", dbDSN)
		if err != nil {
			panic(err)
			return
		}
		SQLDB.DB().SetMaxOpenConns(128)
		SQLDB.DB().SetConnMaxLifetime(120 * time.Second)
	}
}

