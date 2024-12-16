package engine

import (
	"github.com/go-cinderella/cinderella-engine/engine/idgenerator"
)

type ProcessEngineConfiguration interface {
	GetCommandExecutor() Executor
	GetIDGenerator() idgenerator.IDGenerator
	GetDeploymentSettings() map[string]interface{}
	GetBpmnDeployer() Deployer
	GetExpressionManager() ExpressionManager
}
