package entitymanager

import (
	"errors"
	"github.com/go-cinderella/cinderella-engine/engine/datamanager"
	"github.com/go-cinderella/cinderella-engine/engine/errs"
	"github.com/go-cinderella/cinderella-engine/engine/model"
	"github.com/spf13/cast"
	"github.com/unionj-cloud/toolkit/zlogger"
	"github.com/wubin1989/gorm"
)

type DeploymentEntityManager struct {
}

func (deploymentEntityManager DeploymentEntityManager) FindById(deploymentId string) (DeploymentEntity, error) {
	dataManager := datamanager.GetDeploymentDataManager()
	deployment := model.ActReDeployment{}
	err := dataManager.FirstById(deploymentId, &deployment)
	if err != nil {
		zlogger.Error().Err(err).Msgf("get deployment err: %s, deploymentId: %s", err, deploymentId)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return DeploymentEntity{}, errs.ErrDeploymentNotFound
		}
		return DeploymentEntity{}, errs.ErrInternalError
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
		return deploymentEntity, err
	}
	resourceEntity := ResourceEntity{}
	resourceEntity.SetName(cast.ToString(resource.Name_))
	resourceEntity.SetDeploymentId(cast.ToString(resource.DeploymentID_))
	resourceEntity.SetBytes(*resource.Bytes_)
	deploymentEntity.SetResources(resourceEntity)
	return deploymentEntity, nil
}
