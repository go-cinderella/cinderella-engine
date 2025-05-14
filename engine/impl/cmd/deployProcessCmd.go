package cmd

import (
	"context"
	"fmt"
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/go-cinderella/cinderella-engine/engine/impl/converter"
	"github.com/go-cinderella/cinderella-engine/engine/internal/datamanager"
	"github.com/go-cinderella/cinderella-engine/engine/internal/model"
	"github.com/samber/lo"
	"github.com/unionj-cloud/toolkit/stringutils"
	"strings"
	"time"
)

var _ engine.Command = (*DeploymentCmd)(nil)

type DeploymentCmd struct {
	Name          string
	Key           string
	TenantId      string
	Bytes         []byte
	Description   *string
	Ctx           context.Context
	Transactional bool
}

func (receiver DeploymentCmd) IsTransactional() bool {
	return receiver.Transactional
}

func (receiver DeploymentCmd) Context() context.Context {
	return receiver.Ctx
}

func (receiver DeploymentCmd) Execute(ctx engine.Context) (interface{}, error) {
	if stringutils.IsEmpty(receiver.Key) {
		return nil, fmt.Errorf("deployment key required")
	}

	bpmnXMLConverter := converter.BpmnXMLConverter{}
	bpmnModel := bpmnXMLConverter.ConvertToBpmnModel(receiver.Bytes)
	process := bpmnModel.GetProcess()
	processId := process.GetId()

	deployment := model.ActReDeployment{Name_: &receiver.Name, Key_: &receiver.Key, DeployTime_: lo.ToPtr(time.Now().UTC()), ProcessID_: processId}
	deploymentDataManager := datamanager.GetDeploymentDataManager()
	if err := deploymentDataManager.Insert(&deployment); err != nil {
		return nil, err
	}

	processDefinitionEntityManager := entitymanager.GetProcessDefinitionEntityManager()
	resourceName := strings.TrimSuffix(receiver.Name, ".bpmn20.xml") + ".bpmn20.xml"

	bytearray := model.ActGeBytearray{Name_: &resourceName, DeploymentID_: &deployment.ID_, Bytes_: &receiver.Bytes}
	if err := datamanager.GetResourceDataManager().Insert(&bytearray); err != nil {
		return nil, err
	}

	processDefinition := model.ActReProcdef{Name_: &receiver.Name, Key_: receiver.Key, DeploymentID_: &deployment.ID_, ResourceName_: &resourceName, ProcessID_: processId, Description_: receiver.Description}
	result, err := processDefinitionEntityManager.Insert(&processDefinition)
	if err != nil {
		return nil, err
	}

	return result, nil
}
