package entitymanager

type ResourceEntity struct {
	Name         string
	Bytes        []byte
	DeploymentId string
}

func (resourceEntity *ResourceEntity) GetName() string {
	return resourceEntity.Name
}

func (resourceEntity *ResourceEntity) SetName(name string) {
	resourceEntity.Name = name
}

func (resourceEntity ResourceEntity) GetBytes() []byte {
	return resourceEntity.Bytes
}

func (resourceEntity *ResourceEntity) SetBytes(bytes []byte) {
	resourceEntity.Bytes = bytes
}

func (resourceEntity *ResourceEntity) GetDeploymentId() string {
	return resourceEntity.DeploymentId
}

func (resourceEntity *ResourceEntity) SetDeploymentId(deploymentId string) {
	resourceEntity.DeploymentId = deploymentId
}
