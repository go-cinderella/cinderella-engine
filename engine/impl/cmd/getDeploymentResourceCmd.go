package cmd

import (
	"context"
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/go-cinderella/cinderella-engine/engine/internal/errs"
	"github.com/unionj-cloud/toolkit/stringutils"
)

var _ engine.Command = (*GetDeploymentResourceCmd)(nil)

type GetDeploymentResourceCmd struct {
	DeploymentId        string
	ProcessDefinitionId string
	ResourceName        string
	Ctx                 context.Context
	Transactional       bool
}

func (getDeploymentResourceCmd GetDeploymentResourceCmd) IsTransactional() bool {
	return getDeploymentResourceCmd.Transactional
}

func (getDeploymentResourceCmd GetDeploymentResourceCmd) Context() context.Context {
	return getDeploymentResourceCmd.Ctx
}

func (getDeploymentResourceCmd GetDeploymentResourceCmd) Execute(ctx engine.Context) (interface{}, error) {
	if stringutils.IsEmpty(getDeploymentResourceCmd.DeploymentId) && stringutils.IsEmpty(getDeploymentResourceCmd.ProcessDefinitionId) {
		return nil, errs.NewCinderellaIllegalArgumentError("One of deploymentId and processDefinitionId is required")
	}

	var resourceEntity entitymanager.ResourceEntity
	var err error

	deploymentId := getDeploymentResourceCmd.DeploymentId
	resourceName := getDeploymentResourceCmd.ResourceName

	if stringutils.IsNotEmpty(deploymentId) {
		resourceEntityManager := entitymanager.GetResourceEntityManager()
		resourceEntity, err = resourceEntityManager.FindResourceByDeploymentIdAndResourceName(deploymentId, resourceName)
	} else {
		processDefinitionEntityManager := entitymanager.GetProcessDefinitionEntityManager()
		resourceEntity, err = processDefinitionEntityManager.FindResourceEntityByProcessDefinitionById(getDeploymentResourceCmd.ProcessDefinitionId)
	}

	if err != nil {
		return nil, err
	}

	return resourceEntity, nil
}
