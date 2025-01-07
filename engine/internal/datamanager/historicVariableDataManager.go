package datamanager

import (
	"github.com/go-cinderella/cinderella-engine/engine/db"
	"github.com/go-cinderella/cinderella-engine/engine/internal/datamanager/abstract"
	"github.com/go-cinderella/cinderella-engine/engine/variable"
	"github.com/unionj-cloud/toolkit/copier"
	"github.com/wubin1989/gorm/clause"
)

type HistoricVariableDataManager struct {
	abstract.DataManager
}

func (historicVariableDataManager HistoricVariableDataManager) Upsert(varinst variable.HistoricVariable) error {
	var value map[string]interface{}
	copier.DeepCopy(varinst, &value)

	err := db.DB().Model(&variable.HistoricVariable{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name_"}, {Name: "proc_inst_id_"}},
		DoUpdates: clause.AssignmentColumns([]string{"double_", "long_", "text_"}),
	}).Create(value).Error
	return err
}
