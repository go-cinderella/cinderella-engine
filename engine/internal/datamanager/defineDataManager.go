package datamanager

import (
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/internal/datamanager/abstract"
	"github.com/go-cinderella/cinderella-engine/engine/internal/model"
)

type ProcessDefinitionDataManager struct {
	abstract.DataManager
}

func (define ProcessDefinitionDataManager) FindDeployedProcessDefinitionByKey(key string) (model.ActReProcdef, error) {
	processDefinition := model.ActReProcdef{}
	procDefQuery := contextutil.GetQuery().ActReProcdef
	err := procDefQuery.Where(procDefQuery.Key.Eq(key)).Where(procDefQuery.DeploymentID.IsNotNull()).Order(procDefQuery.Version.Desc()).Fetch(&processDefinition)
	return processDefinition, err
}

func (define ProcessDefinitionDataManager) FindDeployedProcessDefinitionByDeploymentId(deploymentId string) (model.ActReProcdef, error) {
	processDefinition := model.ActReProcdef{}
	procDefQuery := contextutil.GetQuery().ActReProcdef
	err := procDefQuery.Where(procDefQuery.DeploymentID.Eq(deploymentId)).Fetch(&processDefinition)
	return processDefinition, err
}
