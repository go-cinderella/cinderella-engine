package entitymanager

import "time"

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

func (processDefinitionEntityImpl *ProcessDefinitionEntity) SetVersion(version int) {
	processDefinitionEntityImpl.Version = version
}

func (processDefinitionEntityImpl *ProcessDefinitionEntity) SetCategory(category string) {
	processDefinitionEntityImpl.Category = category
}

func (processDefinitionEntityImpl *ProcessDefinitionEntity) SetResourceName(resourceName string) {
	processDefinitionEntityImpl.ResourceName = resourceName
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
