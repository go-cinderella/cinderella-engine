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

func (defineManager VariableDataManager) SelectByExecutionIdAndName(name string, executionId string) (variable.Variable, bool, error) {
	var result []variable.Variable

	variableQuery := contextutil.GetQuery().ActRuVariable
	if err := variableQuery.Where(variableQuery.ExecutionID.Eq(executionId)).Where(variableQuery.Name.Eq(name)).Fetch(&result); err != nil {
		return variable.Variable{}, false, err
	}

	if len(result) > 0 {
		return result[0], true, nil
	}

	return variable.Variable{}, false, nil
}

func (defineManager VariableDataManager) SelectByExecutionIdsAndName(executionIds []string, name string) ([]variable.Variable, error) {
	var result []variable.Variable

	variableQuery := contextutil.GetQuery().ActRuVariable
	if err := variableQuery.Where(variableQuery.ExecutionID.In(executionIds...)).Where(variableQuery.Name.Eq(name)).Fetch(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (variableManager VariableDataManager) SelectByTaskIdAndName(name string, taskId string) (variable.Variable, bool, error) {
	var result []variable.Variable

	variableQuery := contextutil.GetQuery().ActRuVariable
	if err := variableQuery.Where(variableQuery.TaskID.Eq(taskId)).Where(variableQuery.Name.Eq(name)).Fetch(&result); err != nil {
		return variable.Variable{}, false, err
	}

	if len(result) > 0 {
		return result[0], true, nil
	}

	return variable.Variable{}, false, nil
}

func (variableManager VariableDataManager) SelectByExecutionId(executionId string) ([]variable.Variable, error) {
	var variables []variable.Variable

	variableQuery := contextutil.GetQuery().ActRuVariable
	err := variableQuery.Where(variableQuery.ExecutionID.Eq(executionId)).Fetch(&variables)

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

func (variableManager VariableDataManager) DeleteByExecutionIdAndNames(executionId string, variableNames []string) error {
	variableQuery := contextutil.GetQuery().ActRuVariable
	_, err := variableQuery.Where(variableQuery.ExecutionID.Eq(executionId), variableQuery.Name.In(variableNames...)).Delete()
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
		Columns:   []clause.Column{{Name: "name_"}, {Name: "execution_id_"}},
		DoUpdates: clause.AssignmentColumns([]string{"double_", "long_", "text_"}),
	}).Create(value).Error
	return err
}
