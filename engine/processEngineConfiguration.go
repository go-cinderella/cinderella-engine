package engine

import (
	"github.com/go-cinderella/cinderella-engine/engine/idgenerator"
	"github.com/go-resty/resty/v2"
)

type ProcessEngineConfiguration interface {
	GetCommandExecutor() Executor
	GetIDGenerator() idgenerator.IDGenerator
	GetDeploymentSettings() map[string]interface{}
	GetBpmnDeployer() Deployer
	GetHttpClient() *resty.Client
}
