package entitymanager

import (
	"errors"
	"github.com/go-cinderella/cinderella-engine/engine/constant"
	"github.com/go-cinderella/cinderella-engine/engine/dto/execution"
	"github.com/go-cinderella/cinderella-engine/engine/dto/historicactinst"
	"github.com/go-cinderella/cinderella-engine/engine/dto/request"
	"github.com/go-cinderella/cinderella-engine/engine/impl/converter"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/go-cinderella/cinderella-engine/engine/internal/datamanager"
	"github.com/go-cinderella/cinderella-engine/engine/internal/model"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/unionj-cloud/toolkit/stringutils"
	"math"
	"strings"
)

type ExecutionEntityManager struct {
}

func (executionEntityManager ExecutionEntityManager) FindById(executionId string) (ExecutionEntity, error) {
	executionDataManager := datamanager.GetExecutionDataManager()
	var execution model.ActRuExecution
	if err := executionDataManager.FindById(executionId, &execution); err != nil {
		return ExecutionEntity{}, err
	}
	if stringutils.IsEmpty(execution.ID_) {
		return ExecutionEntity{}, errors.New("execution not found")
	}

	entityImpl := ExecutionEntity{
		VariableScopeImpl: &VariableScopeImpl{},
	}
	entityImpl.SetId(execution.ID_)
	entityImpl.SetProcessDefinitionId(*execution.ProcDefID_)
	entityImpl.SetProcessInstanceId(*execution.ProcInstID_)
	if execution.ActID_ != nil {
		entityImpl.SetCurrentActivityId(*execution.ActID_)
	}
	return entityImpl, nil
}

func (executionEntityManager ExecutionEntityManager) List(listRequest execution.ListRequest) ([]ExecutionEntity, error) {
	executionDataManager := datamanager.GetExecutionDataManager()
	executions, err := executionDataManager.List(listRequest)
	if err != nil {
		return nil, err
	}
	result := lo.Map[model.ActRuExecution, ExecutionEntity](executions, func(item model.ActRuExecution, index int) ExecutionEntity {
		return ExecutionEntity{
			VariableScopeImpl: &VariableScopeImpl{},
			AbstractEntity: AbstractEntity{
				Id: item.ID_,
			},
			BusinessKey:         cast.ToString(item.BusinessKey_),
			ProcessInstanceId:   cast.ToString(item.ProcInstID_),
			CurrentActivityId:   cast.ToString(item.ActID_),
			BusinessStatus:      cast.ToString(item.BusinessStatus_),
			Suspended:           cast.ToBool(item.SuspensionState_),
			ProcessDefinitionId: cast.ToString(item.ProcDefID_),
			ActivityId:          cast.ToString(item.ActID_),
			StartUserId:         cast.ToString(item.StartUserID_),
			StartTime:           *item.StartTime_,
			CallbackId:          cast.ToString(item.CallbackID_),
			CallbackType:        cast.ToString(item.CallbackType_),
			ReferenceId:         cast.ToString(item.ReferenceID_),
			ReferenceType:       cast.ToString(item.ReferenceType_),
			TenantId:            item.TenantID_,
		}
	})
	return result, nil
}

func (executionEntityManager ExecutionEntityManager) CreateExecution(executionEntity *ExecutionEntity) error {
	executionDataManager := datamanager.GetExecutionDataManager()
	actRuExecution := model.ActRuExecution{
		Rev_:                    lo.ToPtr(cast.ToInt32(1)),
		ProcInstID_:             &executionEntity.ProcessInstanceId,
		ParentID_:               &executionEntity.ProcessInstanceId,
		ProcDefID_:              &executionEntity.ProcessDefinitionId,
		RootProcInstID_:         &executionEntity.ProcessInstanceId,
		ActID_:                  &executionEntity.CurrentActivityId,
		IsActive_:               lo.ToPtr(true),
		IsConcurrent_:           lo.ToPtr(false),
		IsScope_:                lo.ToPtr(false),
		IsEventScope_:           lo.ToPtr(false),
		IsMiRoot_:               lo.ToPtr(false),
		SuspensionState_:        lo.ToPtr(cast.ToInt32(1)),
		StartTime_:              &executionEntity.StartTime,
		IsCountEnabled_:         lo.ToPtr(true),
		EvtSubscrCount_:         lo.ToPtr(cast.ToInt32(0)),
		TaskCount_:              lo.ToPtr(cast.ToInt32(0)),
		JobCount_:               lo.ToPtr(cast.ToInt32(0)),
		TimerJobCount_:          lo.ToPtr(cast.ToInt32(0)),
		SuspJobCount_:           lo.ToPtr(cast.ToInt32(0)),
		DeadletterJobCount_:     lo.ToPtr(cast.ToInt32(0)),
		ExternalWorkerJobCount_: lo.ToPtr(cast.ToInt32(0)),
		VarCount_:               lo.ToPtr(cast.ToInt32(0)),
		IDLinkCount_:            lo.ToPtr(cast.ToInt32(0)),
	}
	if err := executionDataManager.Insert(&actRuExecution); err != nil {
		return err
	}

	executionEntity.SetId(actRuExecution.ID_)
	return nil
}

func (executionEntityManager ExecutionEntityManager) DeleteExecution(executionId string) error {
	// TODO 目前所有变量都是流程实例的变量
	//variableDataManager := datamanager.GetVariableDataManager()
	//err := variableDataManager.DeleteByExecutionId(executionId)
	//if err != nil {
	//	return err
	//}

	executionDataManager := datamanager.GetExecutionDataManager()
	err := executionDataManager.Delete(executionId)
	return err
}

func (executionEntityManager ExecutionEntityManager) RecordBusinessStatus(delegateExecution delegate.DelegateExecution) error {
	historicActinstDataManager := datamanager.GetHistoricActinstDataManager()
	historicActinsts, err := historicActinstDataManager.List(historicactinst.ListRequest{
		ListCommonRequest: request.ListCommonRequest{
			Size: math.MaxInt32,
		},
		Finished:          lo.ToPtr(false),
		ProcessInstanceId: delegateExecution.GetProcessInstanceId(),
		ActivityType:      strings.Join(constant.ELEMENT_TASK_LIST, ","),
	})
	if err != nil {
		return err
	}
	if len(historicActinsts) == 0 {
		return nil
	}

	businessStatus := historicActinsts[0].ActID_
	executionDataManager := datamanager.GetExecutionDataManager()
	err = executionDataManager.RecordBusinessStatus(delegateExecution.GetProcessInstanceId(), businessStatus)
	return err
}

func (executionEntityManager ExecutionEntityManager) MigrateProcessInstanceProcDefIdAndStartActId(oldProcessDefinitionEntity, newProcessDefinitionEntity ProcessDefinitionEntity) error {
	executionDataManager := datamanager.GetExecutionDataManager()

	bpmnXMLConverter := converter.BpmnXMLConverter{}
	bpmnModel := bpmnXMLConverter.ConvertToBpmnModel(newProcessDefinitionEntity.ResourceContent)
	process := bpmnModel.GetProcess()
	startActId := process.InitialFlowElement.GetId()

	if err := executionDataManager.MigrateProcessInstance(oldProcessDefinitionEntity.GetId(), newProcessDefinitionEntity.GetId(), startActId); err != nil {
		return err
	}

	historicProcessDataManager := datamanager.GetHistoricProcessDataManager()
	if err := historicProcessDataManager.Migrate(oldProcessDefinitionEntity.GetId(), newProcessDefinitionEntity.GetId(), startActId); err != nil {
		return err
	}

	return nil
}

func (executionEntityManager ExecutionEntityManager) MigrateProcessInstanceBusinessStatus(processDefinitionEntity ProcessDefinitionEntity, oldActivityId string, newActivityId string) error {
	executionDataManager := datamanager.GetExecutionDataManager()

	if err := executionDataManager.MigrateProcessInstanceBusinessStatus(processDefinitionEntity.GetId(), oldActivityId, newActivityId); err != nil {
		return err
	}

	historicProcessDataManager := datamanager.GetHistoricProcessDataManager()
	if err := historicProcessDataManager.MigrateBusinessStatus(processDefinitionEntity.GetId(), oldActivityId, newActivityId); err != nil {
		return err
	}

	return nil
}

func (executionEntityManager ExecutionEntityManager) MigrateExecutionProcDefID(oldProcessDefinitionEntity, newProcessDefinitionEntity ProcessDefinitionEntity) error {
	executionDataManager := datamanager.GetExecutionDataManager()
	err := executionDataManager.MigrateExecutionProcDefID(oldProcessDefinitionEntity.GetId(), newProcessDefinitionEntity.GetId())
	return err
}

func (executionEntityManager ExecutionEntityManager) MigrateExecutionActID(processDefinitionEntity ProcessDefinitionEntity, oldActivityId string, newActivityId string) error {
	executionDataManager := datamanager.GetExecutionDataManager()
	err := executionDataManager.MigrateExecutionActID(processDefinitionEntity.GetId(), oldActivityId, newActivityId)
	return err
}
