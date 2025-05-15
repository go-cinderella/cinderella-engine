package entitymanager

import (
	"time"
)

type ProcessDefinitionEntity struct {
	AbstractEntity
	Name            string
	Description     string
	Key             string
	Version         int
	Category        string
	DeploymentId    string
	ResourceName    string
	ResourceContent []byte
	DeployTime      time.Time
	CreatedBy       string
	CreatedByName   string
}

func (processDefinitionEntityImpl *ProcessDefinitionEntity) SetKey(key string) {
	processDefinitionEntityImpl.Key = key
}

func (processDefinitionEntityImpl *ProcessDefinitionEntity) SetName(name string) {
	processDefinitionEntityImpl.Name = name
}

func (processDefinitionEntityImpl *ProcessDefinitionEntity) SetDescription(description string) {
	processDefinitionEntityImpl.Description = description
}

func (processDefinitionEntityImpl *ProcessDefinitionEntity) SetDeploymentId(deploymentId string) {
	processDefinitionEntityImpl.DeploymentId = deploymentId
}

func (processDefinitionEntityImpl *ProcessDefinitionEntity) SetDeployTime(deployTime time.Time) {
	processDefinitionEntityImpl.DeployTime = deployTime
}

func (processDefinitionEntityImpl *ProcessDefinitionEntity) SetVersion(version int) {
	processDefinitionEntityImpl.Version = version
}

func (processDefinitionEntityImpl *ProcessDefinitionEntity) SetCreatedBy(createdBy string) {
	processDefinitionEntityImpl.CreatedBy = createdBy
}

func (processDefinitionEntityImpl *ProcessDefinitionEntity) SetCreatedByName(createdByName string) {
	processDefinitionEntityImpl.CreatedByName = createdByName
}

func (processDefinitionEntityImpl *ProcessDefinitionEntity) SetCategory(category string) {
	processDefinitionEntityImpl.Category = category
}

func (processDefinitionEntityImpl *ProcessDefinitionEntity) SetResourceName(resourceName string) {
	processDefinitionEntityImpl.ResourceName = resourceName
}

func (processDefinitionEntityImpl *ProcessDefinitionEntity) SetResourceContent(resourceContent []byte) {
	processDefinitionEntityImpl.ResourceContent = resourceContent
}

func (processDefinitionEntityImpl *ProcessDefinitionEntity) GetKey() string {
	return processDefinitionEntityImpl.Key
}

func (processDefinitionEntityImpl *ProcessDefinitionEntity) GetName() string {
	return processDefinitionEntityImpl.Name
}

func (processDefinitionEntityImpl *ProcessDefinitionEntity) GetDescription() string {
	return processDefinitionEntityImpl.Description
}
func (processDefinitionEntityImpl *ProcessDefinitionEntity) GetDeploymentId() string {
	return processDefinitionEntityImpl.DeploymentId
}

func (processDefinitionEntityImpl *ProcessDefinitionEntity) GetVersion() int {
	return processDefinitionEntityImpl.Version
}

func (processDefinitionEntityImpl *ProcessDefinitionEntity) GetResourceName() string {
	return processDefinitionEntityImpl.ResourceName
}
