package cmd

import (
	"context"
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	bpmn_model "github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/go-cinderella/cinderella-engine/engine/internal/datamanager"
	model "github.com/go-cinderella/cinderella-engine/engine/internal/model"
	"github.com/go-cinderella/cinderella-engine/engine/utils"
	"github.com/samber/lo"
	"github.com/unionj-cloud/toolkit/stringutils"
	"time"
)

var _ engine.Command = (*StartProcessInstanceByKeyCmd)(nil)

type StartProcessInstanceByKeyCmd struct {
	ProcessInstanceId    string
	ProcessDefinitionKey string
	Variables            map[string]interface{}
	BusinessKey          string
	TenantId             string
	UserId               string
	StartActivityId      string
	StartActivityName    string
	Ctx                  context.Context
	Transactional        bool
}

func (receiver StartProcessInstanceByKeyCmd) IsTransactional() bool {
	return receiver.Transactional
}

func (receiver StartProcessInstanceByKeyCmd) Context() context.Context {
	return receiver.Ctx
}

func (receiver StartProcessInstanceByKeyCmd) Start(ctx engine.Context) (entitymanager.ExecutionEntity, error) {
	processDefinitionEntityManager := entitymanager.GetProcessDefinitionEntityManager()
	definitionEntity, err := processDefinitionEntityManager.FindLatestProcessDefinitionByKey(receiver.ProcessDefinitionKey)
	if err != nil {
		return entitymanager.ExecutionEntity{}, err
	}

	processUtils := utils.ProcessDefinitionUtil{}
	process, err := processUtils.GetProcess(definitionEntity.GetId())
	if err != nil {
		return entitymanager.ExecutionEntity{}, err
	}

	var element delegate.FlowElement
	if stringutils.IsNotEmpty(receiver.StartActivityId) {
		// 从指定节点开始流程流转
		element = process.GetFlowElement(receiver.StartActivityId)
	} else if stringutils.IsNotEmpty(receiver.StartActivityName) {
		flowElementNameMap := make(map[string]delegate.FlowElement)
		lo.ForEach(process.FlowElementList, func(item delegate.FlowElement, index int) {
			flowElementNameMap[item.GetName()] = item
		})

		element = flowElementNameMap[receiver.StartActivityName]
	} else {
		// 默认从开始节点开始流程流转
		flowElement := process.InitialFlowElement
		element = flowElement.(*bpmn_model.StartEvent)
	}

	processInstance := model.ActRuExecution{}
	processInstance.Rev_ = lo.ToPtr(int32(1))
	if stringutils.IsNotEmpty(receiver.ProcessInstanceId) {
		processInstance.ID_ = receiver.ProcessInstanceId
	} else {
		processInstance.ID_, _ = contextutil.GetIDGenerator().NextID()
	}
	processInstance.BusinessKey_ = &receiver.BusinessKey
	processInstance.TenantID_ = &receiver.TenantId
	processInstance.StartTime_ = lo.ToPtr(time.Now().UTC())
	processInstance.ProcDefID_ = lo.ToPtr(definitionEntity.GetId())
	processInstance.IsActive_ = lo.ToPtr(true)
	processInstance.StartUserID_ = &receiver.UserId
	processInstance.StartActID_ = lo.ToPtr(element.GetId())
	processInstance.ProcInstID_ = &processInstance.ID_
	processInstance.RootProcInstID_ = &processInstance.ID_
	processInstance.IsConcurrent_ = lo.ToPtr(false)
	processInstance.IsScope_ = lo.ToPtr(true)
	processInstance.IsEventScope_ = lo.ToPtr(false)
	processInstance.IsMiRoot_ = lo.ToPtr(false)
	processInstance.SuspensionState_ = lo.ToPtr(int32(1))
	processInstance.IsCountEnabled_ = lo.ToPtr(true)

	// 生成流程实例
	executionDataManager := datamanager.GetExecutionDataManager()
	if err = executionDataManager.CreateProcessInstance(&processInstance); err != nil {
		return entitymanager.ExecutionEntity{}, err
	}

	processInstanceEntity := entitymanager.ExecutionEntity{
		VariableScopeImpl: &entitymanager.VariableScopeImpl{},
	}
	processInstanceEntity.SetId(processInstance.ID_)
	processInstanceEntity.SetProcessDefinitionId(*processInstance.ProcDefID_)
	processInstanceEntity.SetProcessInstanceId(*processInstance.ProcInstID_)

	//保存流程变量
	if err = processInstanceEntity.SetVariable(&processInstanceEntity, receiver.Variables); err != nil {
		return entitymanager.ExecutionEntity{}, err
	}

	if processInstance.StartUserID_ != nil {
		identityLinkManager := entitymanager.GetIdentityLinkManager()
		link := model.ActRuIdentitylink{
			Rev_:        lo.ToPtr(int32(1)),
			Type_:       lo.ToPtr("starter"),
			UserID_:     processInstance.StartUserID_,
			ProcInstID_: processInstance.ProcInstID_,
		}
		err = identityLinkManager.CreateIdentityLink(link)
		if err != nil {
			return entitymanager.ExecutionEntity{}, err
		}
	}

	// 生成执行实例
	execution := entitymanager.ExecutionEntity{
		VariableScopeImpl:   &entitymanager.VariableScopeImpl{},
		ProcessInstanceId:   processInstance.ID_,
		ProcessDefinitionId: definitionEntity.GetId(),
		StartTime:           time.Now().UTC(),
	}
	execution.SetCurrentFlowElement(element)

	executionEntityManager := entitymanager.GetExecutionEntityManager()
	if err = executionEntityManager.CreateExecution(&execution); err != nil {
		return entitymanager.ExecutionEntity{}, err
	}

	return execution, nil
}

func (receiver StartProcessInstanceByKeyCmd) Execute(ctx engine.Context) (interface{}, error) {
	execution, err := receiver.Start(ctx)
	if err != nil {
		return nil, err
	}
	contextutil.GetAgendaFromContext(ctx).PlanContinueProcessOperation(&execution)
	return execution, nil
}
