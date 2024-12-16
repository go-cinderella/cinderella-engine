package parse

import (
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
)

type ParsedDeploymentBuilderFactory struct {
	BpmnParser *BpmnParser
}

func (parsedDeploymentBuilderFactory ParsedDeploymentBuilderFactory) GetBuilderForDeploymentAndSettings(deployment entitymanager.DeploymentEntity, deploymentSettings map[string]interface{}) ParsedDeploymentBuilder {
	return ParsedDeploymentBuilder{deployment, parsedDeploymentBuilderFactory.BpmnParser, deploymentSettings}
}
