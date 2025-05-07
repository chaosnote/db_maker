package main

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"

	"github.com/chaosnote/db_maker/utils"
)

type DBStore interface{}

type db_store struct {
	utils.LogStore

	db *sql.DB
}

func NewDBStore() DBStore {
	ref := db_store{
		LogStore: utils.NewConsoleLogger(1),
	}

	defer ref.Flush()

	// 例 : "user:password@tcp(ip)?parseTime=true/dbname"
	cmd := fmt.Sprintf(`%s:%s@tcp(%s)/%s?parseTime=true`, "chris", "123456", "192.168.0.236:3306", "game_dev")
	ref.Debug(utils.LogFields{"cmd": cmd})

	var e error
	ref.db, e = sql.Open("mysql", cmd)
	if e != nil {
		ref.Panic(e)
	}
	e = ref.db.Ping()
	if e != nil {
		ref.Panic(e)
	}

	ref.addSPDropTables()
	ref.addSPDropTable()
	ref.addSPSearchTables()
	ref.addSPUpsertUser()
	ref.addSPWallet()

	e = filepath.Walk("./asset/db/_sql", func(file_path string, info fs.FileInfo, e error) error {
		if e != nil {
			return e
		}
		if info.IsDir() {
			return nil
		}

		ref.Debug(utils.LogFields{"path": info.Name()})
		ref.execSQLFile(file_path)

		return e
	})
	if e != nil {
		ref.Panic(e)
	}

	ref.addBaseData()

	return ref
}

//-----------------------------------------------

func (ds *db_store) execSQLFile(file_path string) {
	var e error
	var content []byte

	defer func() {
		if e != nil {
			ds.Panic(e)
		}
	}()

	content, e = os.ReadFile(file_path)
	if e != nil {
		return
	}

	ds.execSQLText(string(content))
	if e != nil {
		return
	}
}

func (ds *db_store) execSQLText(content string) (e error) {
	_, e = ds.db.Exec(content)
	if e != nil {
		ds.Panic(e)
		return
	}
	return
}
