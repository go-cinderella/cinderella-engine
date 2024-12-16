package entitymanager

import (
	"github.com/go-cinderella/cinderella-engine/engine/internal/datamanager"
	"github.com/spf13/cast"
)

type ResourceEntityManager struct {
}

func (resourceEntityManager ResourceEntityManager) FindResourceByDeploymentIdAndResourceName(deploymentId, resourceName string) (ResourceEntity, error) {
	resourceDataManager := datamanager.GetResourceDataManager()
	resource, err := resourceDataManager.FindResourceByDeploymentIdAndResourceName(deploymentId, resourceName)
	if err != nil {
		return ResourceEntity{}, err
	}

	resourceEntity := ResourceEntity{}
	resourceEntity.SetName(cast.ToString(resource.Name_))
	resourceEntity.SetDeploymentId(cast.ToString(resource.DeploymentID_))
	resourceEntity.SetBytes(*resource.Bytes_)

	return resourceEntity, nil
}
