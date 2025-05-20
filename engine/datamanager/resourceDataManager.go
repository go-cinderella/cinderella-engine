package datamanager

import (
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/internal/datamanager/abstract"
	. "github.com/go-cinderella/cinderella-engine/engine/internal/model"
	"github.com/unionj-cloud/toolkit/stringutils"
)

type ResourceDataManager struct {
	abstract.DataManager
}

func (define ResourceDataManager) FindResourceByDeploymentIdAndResourceName(deploymentId, resourceName string) (ActGeBytearray, error) {
	bytearray := ActGeBytearray{}
	bytearrayQuery := contextutil.GetQuery().ActGeBytearray

	queryBuilder := bytearrayQuery.Where(bytearrayQuery.DeploymentID.Eq(deploymentId))
	if stringutils.IsNotEmpty(resourceName) {
		queryBuilder = queryBuilder.Where(bytearrayQuery.Name.Eq(resourceName))
	}

	err := queryBuilder.Fetch(&bytearray)
	return bytearray, err
}
