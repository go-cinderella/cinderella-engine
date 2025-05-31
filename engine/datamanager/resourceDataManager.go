package datamanager

import (
	"errors"
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/datamanager/abstract"
	"github.com/go-cinderella/cinderella-engine/engine/errs"
	"github.com/go-cinderella/cinderella-engine/engine/model"
	"github.com/unionj-cloud/toolkit/stringutils"
	"github.com/unionj-cloud/toolkit/zlogger"
	"github.com/wubin1989/gorm"
)

type ResourceDataManager struct {
	abstract.DataManager
}

// FindResourceByDeploymentIdAndResourceName may raise errs.ErrResourceNotFound or errs.ErrInternalError error
func (define ResourceDataManager) FindResourceByDeploymentIdAndResourceName(deploymentId, resourceName string) (model.ActGeBytearray, error) {
	bytearrayQuery := contextutil.GetQuery().ActGeBytearray

	queryBuilder := bytearrayQuery.Where(bytearrayQuery.DeploymentID.Eq(deploymentId))
	if stringutils.IsNotEmpty(resourceName) {
		queryBuilder = queryBuilder.Where(bytearrayQuery.Name.Eq(resourceName))
	}

	result, err := queryBuilder.First()
	if err != nil {
		zlogger.Error().Err(err).Msgf("get resource err: %s, deploymentId: %s, resourceName: %s", err, deploymentId, resourceName)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ActGeBytearray{}, errs.ErrResourceNotFound
		}
		return model.ActGeBytearray{}, errs.ErrInternalError
	}
	return *result, nil
}
