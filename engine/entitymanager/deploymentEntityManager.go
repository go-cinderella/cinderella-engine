package entitymanager

import (
	"github.com/go-cinderella/cinderella-engine/engine/datamanager"
	"github.com/go-cinderella/cinderella-engine/engine/model"
	"github.com/spf13/cast"
)

type DeploymentEntityManager struct {
}

func (deploymentEntityManager DeploymentEntityManager) FindById(deploymentId string) (DeploymentEntity, error) {
	dataManager := datamanager.GetDeploymentDataManager()
	deployment := model.ActReDeployment{}
	err := dataManager.FindById(deploymentId, &deployment)
	if err != nil {
		return DeploymentEntity{}, err
	}
	deploymentEntity := DeploymentEntity{}
	deploymentEntity.ProcessId = deployment.ProcessID_
	deploymentEntity.SetName(cast.ToString(deployment.Name_))
	if deployment.Key_ != nil {
		deploymentEntity.SetKey(*deployment.Key_)
	} else {
		processDefinitionDataManager := datamanager.GetProcessDefinitionDataManager()
		procdef, err := processDefinitionDataManager.FindDeployedProcessDefinitionByDeploymentId(deploymentId)
		if err != nil {
			return DeploymentEntity{}, err
		}
		deploymentEntity.SetKey(procdef.Key_)
	}
	resourceDataManager := datamanager.GetResourceDataManager()
	resource, err := resourceDataManager.FindResourceByDeploymentIdAndResourceName(deployment.ID_, cast.ToString(deployment.Name_)+".bpmn20.xml")
	if err != nil {
		panic(err)
	}
	resourceEntity := ResourceEntity{}
	resourceEntity.SetName(cast.ToString(resource.Name_))
	resourceEntity.SetDeploymentId(cast.ToString(resource.DeploymentID_))
	resourceEntity.SetBytes(*resource.Bytes_)
	deploymentEntity.SetResources(resourceEntity)
	return deploymentEntity, nil
}
