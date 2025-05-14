package postgres

import (
	"github.com/pressly/goose/v3"
)

var Migrations = []*goose.Migration{
	goose.NewGoMigration(
		1,
		&goose.GoFunc{RunTx: Up00001},
		&goose.GoFunc{RunTx: Down00001},
	),
	goose.NewGoMigration(
		2,
		&goose.GoFunc{RunTx: Up00002},
		&goose.GoFunc{RunTx: Down00002},
	),
}
