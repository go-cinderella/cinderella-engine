package entitymanager

import (
	"time"
)

type DeploymentEntity struct {
	Name           string
	Key            string
	ProcessId      string
	TenantId       string
	DeploymentTime time.Time
	IsNew          bool
	ResourceEntity ResourceEntity
}

func (deploymentEntity *DeploymentEntity) AddResource(resource ResourceEntity) {
	deploymentEntity.ResourceEntity = resource
}

func (deploymentEntity *DeploymentEntity) GetResources() ResourceEntity {
	return deploymentEntity.ResourceEntity
}

func (deploymentEntity *DeploymentEntity) GetName() string {
	return deploymentEntity.Name
}

func (deploymentEntity *DeploymentEntity) SetName(name string) {

}

func (deploymentEntity *DeploymentEntity) SetKey(key string) {
	deploymentEntity.Key = key
}
func (deploymentEntity *DeploymentEntity) GetKey() string {
	return deploymentEntity.Key
}
func (deploymentEntity *DeploymentEntity) SetTenantId(tenantId string) {

}

func (deploymentEntity *DeploymentEntity) SetResources(resourceEntity ResourceEntity) {
	deploymentEntity.ResourceEntity = resourceEntity
}

func (deploymentEntity *DeploymentEntity) SetDeploymentTime(deploymentTime time.Time) {

}
