package utils

import (
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
)

// ProcessDefinitionUtil must be used in a command execute function
type ProcessDefinitionUtil struct {
}

func (processDefinitionUtil ProcessDefinitionUtil) GetProcess(processDefinitionId string) (model.Process, error) {
	processDefinitionManager := entitymanager.GetProcessDefinitionEntityManager()
	processDefinitionEntity, err := processDefinitionManager.FindProcessDefinitionById(processDefinitionId)
	if err != nil {
		return model.Process{}, err
	}
	process, err := processDefinitionUtil.resolveProcess(processDefinitionEntity)
	if err != nil {
		return model.Process{}, err
	}
	return process, nil
}

func (processDefinitionUtil ProcessDefinitionUtil) deploy(deployment entitymanager.DeploymentEntity) model.Process {
	deployer := contextutil.GetBpmnDeployer()
	deploymentSettings := contextutil.GetDeploymentSettings()
	deployer.Deploy(deployment, deploymentSettings)
	return deployer.GetProcess(deployment.ProcessId)
}

func (processDefinitionUtil ProcessDefinitionUtil) resolveProcess(processDefinitionEntity entitymanager.ProcessDefinitionEntity) (model.Process, error) {
	deploymentEntityManager := entitymanager.GetDeploymentEntityManager()

	deploymentEntity, err := deploymentEntityManager.FindById(processDefinitionEntity.GetDeploymentId())
	if err != nil {
		return model.Process{}, err
	}

	process := processDefinitionUtil.deploy(deploymentEntity)
	return process, nil
}
