package interceptor

import (
	"database/sql"
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/db"
	"github.com/wubin1989/gorm"
)

type TransactionContextInterceptor struct {
	Next engine.Interceptor
}

func (transactionContextInterceptor TransactionContextInterceptor) Execute(command engine.Command) (value interface{}, err error) {

	sqlSession := db.DB()

	if _, ok := sqlSession.Statement.ConnPool.(*sql.Tx); ok {
		// already in transaction which is managed by outside
		return transactionContextInterceptor.Next.Execute(command)
	}

	defer db.ClearTXDB()

	ctx := command.Context()
	if ctx != nil {
		sqlSession = sqlSession.WithContext(ctx)
	}

	if command.IsTransactional() {
		sqlSession.Transaction(func(tx *gorm.DB) error {
			db.InitTXDB(tx)
			value, err = transactionContextInterceptor.Next.Execute(command)
			return err
		})

		return value, err
	}

	db.InitTXDB(sqlSession)

	return transactionContextInterceptor.Next.Execute(command)
}

func (transactionContextInterceptor *TransactionContextInterceptor) SetNext(next engine.Interceptor) {
	transactionContextInterceptor.Next = next
}
