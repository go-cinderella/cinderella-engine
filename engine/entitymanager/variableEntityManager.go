package entitymanager

import (
	"github.com/go-cinderella/cinderella-engine/engine/internal/datamanager"
	"github.com/go-cinderella/cinderella-engine/engine/variable"
	log "github.com/sirupsen/logrus"
	"github.com/unionj-cloud/toolkit/copier"
)

type VariableEntityManager struct {
}

func (variableEntityManager VariableEntityManager) UpsertVariable(v variable.Variable) error {
	variableDataManager := datamanager.GetVariableDataManager()
	if err := variableDataManager.Upsert(v); err != nil {
		return err
	}
	if err := variableEntityManager.UpsertHistoricVariable(v); err != nil {
		return err
	}
	return nil
}

func (variableEntityManager VariableEntityManager) UpsertHistoricVariable(v variable.Variable) (err error) {
	var historicVariable variable.HistoricVariable
	copier.DeepCopy(v, &historicVariable)

	historicVariableManager := datamanager.GetHistoricVariableDataManager()
	return historicVariableManager.Upsert(historicVariable)
}

func (variableEntityManager VariableEntityManager) DeleteVariableInstanceByTask(taskId string) {
	variableDataManager := datamanager.GetVariableDataManager()
	variables, err := variableDataManager.SelectByTaskId(taskId)
	if err != nil {
		log.Error("select by taskId err：", err)
		return
	}
	for _, variable := range variables {
		variableDataManager.Delete(variable.ID_)
	}
}

func (variableEntityManager VariableEntityManager) DeleteByProcessInstanceId(processInstanceId string) error {
	variableDataManager := datamanager.GetVariableDataManager()
	err := variableDataManager.DeleteByProcessInstanceId(processInstanceId)
	return err
}
