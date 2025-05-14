package db

import (
	"fmt"
	"github.com/go-cinderella/cinderella-engine/engine/runtime"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"
	"github.com/samber/do"
	"github.com/unionj-cloud/toolkit/stringutils"
	"github.com/wubin1989/gorm"
	"strings"
	"sync"
)

const GooseTableName = goose.DefaultTablename + "_act"

var txHolder *sync.Map

func init() {
	txHolder = new(sync.Map)
}

func InitTXDB(db *gorm.DB) {
	txHolder.Store(runtime.GoroutineId(), db)
}

func ClearTXDB() {
	txHolder.Delete(runtime.GoroutineId())
}

func DB() *gorm.DB {
	db, ok := txHolder.Load(runtime.GoroutineId())
	if !ok {
		return nil
	}
	return db.(*gorm.DB)
}

func GetTableName(database, table string) string {
	var tableName string

	if stringutils.IsNotEmpty(database) {
		tableName = fmt.Sprintf("%s.%s", database, table)
	} else {
		tableName = table
	}

	db := do.MustInvoke[*gorm.DB](nil)

	var builder strings.Builder
	db.Dialector.QuoteTo(&builder, tableName)
	return builder.String()
}
