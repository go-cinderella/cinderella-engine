package engine

import (
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
)

type Deployer interface {
	Deploy(deploymentEntity any, deploymentSettings map[string]interface{})
	GetProcess(key string) model.Process
}
