package db

import (
	"github.com/go-cinderella/cinderella-engine/engine/runtime"
	_ "github.com/go-sql-driver/mysql"
	"github.com/samber/do"
	"github.com/wubin1989/gorm"
	"sync"
)

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
		return do.MustInvoke[*gorm.DB](nil)
	}
	return db.(*gorm.DB)
}
