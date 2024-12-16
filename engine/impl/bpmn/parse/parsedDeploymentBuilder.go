package parse

import (
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
)

type ParsedDeploymentBuilder struct {
	deployment         entitymanager.DeploymentEntity
	bpmnParser         *BpmnParser
	deploymentSettings map[string]interface{}
}

func NewParsedDeploymentBuilder(deployment entitymanager.DeploymentEntity, bpmnParser *BpmnParser, deploymentSettings map[string]interface{}) ParsedDeploymentBuilder {
	return ParsedDeploymentBuilder{deployment, bpmnParser, deploymentSettings}
}

func (parsedDeploymentBuilder ParsedDeploymentBuilder) Build() ParsedDeployment {
	resources := parsedDeploymentBuilder.deployment.GetResources()
	bpmnParse := parsedDeploymentBuilder.createBpmnParseFromResource(resources)
	return ParsedDeployment{BpmnParse: bpmnParse}
}

func (parsedDeploymentBuilder ParsedDeploymentBuilder) createBpmnParseFromResource(resource entitymanager.ResourceEntity) BpmnParse {
	name := resource.GetName()
	bytes := resource.GetBytes()
	bpmnParse := parsedDeploymentBuilder.bpmnParser.CreateParse().SourceInputStream(bytes).Deployment(parsedDeploymentBuilder.deployment).SourceName(name)
	bpmnParse.Execute()
	return bpmnParse
}
