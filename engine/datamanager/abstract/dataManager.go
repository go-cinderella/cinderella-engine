package abstract

import (
	"github.com/go-cinderella/cinderella-engine/engine/db"
	"github.com/go-cinderella/cinderella-engine/engine/model"
)

type DataManager struct {
	model.AbstractModel
}

func (dataManagers DataManager) Insert(data interface{}) error {
	err := db.DB().Create(data).Error
	return err
}

func (dataManagers DataManager) FindById(id string, data interface{}) error {
	err := db.DB().Where(`"id_" = ?`, id).First(data).Error
	return err
}
func (dataManagers DataManager) Delete(id string) error {
	tableName := dataManagers.TableName()
	err := db.DB().Where(`"id_" = ?`, id).Table(tableName).Delete(nil).Error
	return err
}
