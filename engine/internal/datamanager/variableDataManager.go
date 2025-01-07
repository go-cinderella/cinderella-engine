package datamanager

import (
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/db"
	"github.com/go-cinderella/cinderella-engine/engine/internal/datamanager/abstract"
	"github.com/go-cinderella/cinderella-engine/engine/internal/errs"
	"github.com/go-cinderella/cinderella-engine/engine/variable"
	log "github.com/sirupsen/logrus"
	"github.com/unionj-cloud/toolkit/copier"
	"github.com/wubin1989/gorm/clause"
)

type VariableDataManager struct {
	abstract.DataManager
}

func (defineManager VariableDataManager) SelectByProcessInstanceIdAndName(name string, processInstanceId string) (variable.Variable, error) {
	vari := variable.Variable{}
	variableQuery := contextutil.GetQuery().ActRuVariable
	err := variableQuery.Where(variableQuery.ProcInstID.Eq(processInstanceId)).Where(variableQuery.Name.Eq(name)).Fetch(&vari)
	if err != nil {
		return variable.Variable{}, err
	}
	return vari, nil
}

func (variableManager VariableDataManager) SelectTaskId(name string, taskId string) (variable.Variable, error) {
	vari := variable.Variable{}
	variableQuery := contextutil.GetQuery().ActRuVariable
	err := variableQuery.Where(variableQuery.TaskID.Eq(taskId)).Where(variableQuery.Name.Eq(name)).Fetch(&vari)
	if err != nil {
		return variable.Variable{}, err
	}
	return vari, nil
}

func (variableManager VariableDataManager) SelectByProcessInstanceId(processInstanceId string) ([]variable.Variable, error) {
	var variables []variable.Variable

	variableQuery := contextutil.GetQuery().ActRuVariable
	err := variableQuery.Where(variableQuery.ProcInstID.Eq(processInstanceId)).Fetch(&variables)

	return variables, err
}

func (variableManager VariableDataManager) DeleteByProcessInstanceId(processInstanceId string) error {
	variableQuery := contextutil.GetQuery().ActRuVariable
	_, err := variableQuery.Where(variableQuery.ProcInstID.Eq(processInstanceId)).Delete()
	return err
}

func (variableManager VariableDataManager) DeleteByExecutionId(executionId string) error {
	variableQuery := contextutil.GetQuery().ActRuVariable
	_, err := variableQuery.Where(variableQuery.ExecutionID.Eq(executionId)).Delete()
	return err
}

func (variableManager VariableDataManager) SelectByTaskId(taskId string) ([]variable.Variable, error) {
	variables := make([]variable.Variable, 0)
	variableQuery := contextutil.GetQuery().ActRuVariable
	err := variableQuery.Where(variableQuery.TaskID.Eq(taskId)).Fetch(&variables)
	if err != nil {
		log.Infoln("Select Variable err: ", err)
		return variables, err
	}
	if variables != nil && len(variables) > 0 {
		return variables, nil
	}
	return variables, errs.CinderellaError{Code: "1001", Msg: "Not Find"}
}

func (variableManager VariableDataManager) Delete(variableId string) error {
	variableQuery := contextutil.GetQuery().ActRuVariable
	_, err := variableQuery.Where(variableQuery.ID.Eq(variableId)).Delete()
	if err != nil {
		log.Infoln("delete Variable err: ", err)
	}
	return err
}

func (variableManager VariableDataManager) Upsert(varinst variable.Variable) error {
	var value map[string]interface{}
	copier.DeepCopy(varinst, &value)

	err := db.DB().Model(&variable.Variable{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name_"}, {Name: "proc_inst_id_"}},
		DoUpdates: clause.AssignmentColumns([]string{"double_", "long_", "text_"}),
	}).Create(value).Error
	return err
}
