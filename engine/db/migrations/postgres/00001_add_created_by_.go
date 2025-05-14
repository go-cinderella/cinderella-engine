package postgres

import (
	"context"
	"database/sql"
	"github.com/go-cinderella/cinderella-engine/engine/config"
	"github.com/go-cinderella/cinderella-engine/engine/db"
)

func Up00001(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `alter table `+db.GetTableName(config.G_Config.Db.Name, "act_re_procdef")+` add created_by_ varchar;`)
	return err
}

func Down00001(ctx context.Context, tx *sql.Tx) error {
	return nil
}
