package deployer

import (
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/parse"
)

var _ engine.Deployer = (*BpmnDeployer)(nil)

type BpmnDeployer struct {
	ParsedDeployment               parse.ParsedDeployment
	ParsedDeploymentBuilderFactory parse.ParsedDeploymentBuilderFactory
}

func (bpmnDeployer *BpmnDeployer) Deploy(entity any, deploymentSettings map[string]interface{}) {
	deploymentEntity := entity.(entitymanager.DeploymentEntity)
	parsedDeployment := bpmnDeployer.ParsedDeploymentBuilderFactory.GetBuilderForDeploymentAndSettings(deploymentEntity, deploymentSettings).Build()
	bpmnDeployer.ParsedDeployment = parsedDeployment
}
func (bpmnDeployer BpmnDeployer) GetProcess(key string) model.Process {
	return *bpmnDeployer.ParsedDeployment.BpmnParse.BpmnModel.GetProcessById(key)
}
func (bpmnDeployer BpmnDeployer) SetParsedDeploymentBuilderFactory(parsedDeploymentBuilderFactory parse.ParsedDeploymentBuilderFactory) {
	bpmnDeployer.ParsedDeploymentBuilderFactory = parsedDeploymentBuilderFactory
}
