package cmd

import (
	"context"
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	bpmn_model "github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/internal/datamanager"
	model "github.com/go-cinderella/cinderella-engine/engine/internal/model"
	"github.com/go-cinderella/cinderella-engine/engine/utils"
	"github.com/samber/lo"
	"time"
)

var _ engine.Command = (*StartProcessInstanceByKeyCmd)(nil)

type StartProcessInstanceByKeyCmd struct {
	ProcessDefinitionKey string
	Variables            map[string]interface{}
	BusinessKey          string
	TenantId             string
	UserId               string
	Ctx                  context.Context
	Transactional        bool
}

func (startProcessInstanceByKeyCmd StartProcessInstanceByKeyCmd) IsTransactional() bool {
	return startProcessInstanceByKeyCmd.Transactional
}

func (startProcessInstanceByKeyCmd StartProcessInstanceByKeyCmd) Context() context.Context {
	return startProcessInstanceByKeyCmd.Ctx
}

func (startProcessInstanceByKeyCmd StartProcessInstanceByKeyCmd) Start(ctx engine.Context) (entitymanager.ExecutionEntity, error) {
	processDefinitionEntityManager := entitymanager.GetProcessDefinitionEntityManager()
	definitionEntity, err := processDefinitionEntityManager.FindLatestProcessDefinitionByKey(startProcessInstanceByKeyCmd.ProcessDefinitionKey)
	if err != nil {
		return entitymanager.ExecutionEntity{}, err
	}

	processUtils := utils.ProcessDefinitionUtil{}
	process, err := processUtils.GetProcess(definitionEntity.GetId())
	if err != nil {
		return entitymanager.ExecutionEntity{}, err
	}

	//获取开始节点
	flowElement := process.InitialFlowElement
	element := flowElement.(*bpmn_model.StartEvent)

	processInstance := model.ActRuExecution{}
	processInstance.Rev_ = lo.ToPtr(int32(1))
	processInstance.ID_, _ = contextutil.GetIDGenerator().NextID()
	processInstance.BusinessKey_ = &startProcessInstanceByKeyCmd.BusinessKey
	processInstance.TenantID_ = &startProcessInstanceByKeyCmd.TenantId
	processInstance.StartTime_ = lo.ToPtr(time.Now().UTC())
	processInstance.ProcDefID_ = lo.ToPtr(definitionEntity.GetId())
	processInstance.IsActive_ = lo.ToPtr(true)
	processInstance.StartUserID_ = &startProcessInstanceByKeyCmd.UserId
	processInstance.StartActID_ = &element.Id
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
	if err = processInstanceEntity.SetVariable(&processInstanceEntity, startProcessInstanceByKeyCmd.Variables); err != nil {
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

func (startProcessInstanceByKeyCmd StartProcessInstanceByKeyCmd) Execute(ctx engine.Context) (interface{}, error) {
	execution, err := startProcessInstanceByKeyCmd.Start(ctx)
	if err != nil {
		return nil, err
	}
	contextutil.GetAgendaFromContext(ctx).PlanContinueProcessOperation(&execution)
	return execution, nil
}
