package entitymanager

import (
	"errors"
	"math"
	"strings"

	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/unionj-cloud/toolkit/stringutils"

	"github.com/go-cinderella/cinderella-engine/engine/constant"
	"github.com/go-cinderella/cinderella-engine/engine/datamanager"
	"github.com/go-cinderella/cinderella-engine/engine/dto/execution"
	"github.com/go-cinderella/cinderella-engine/engine/dto/historicactinst"
	"github.com/go-cinderella/cinderella-engine/engine/dto/request"
	"github.com/go-cinderella/cinderella-engine/engine/impl/converter"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/go-cinderella/cinderella-engine/engine/model"
	"github.com/go-cinderella/cinderella-engine/engine/variable"
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

	entityImpl := ExecutionEntity{}
	entityImpl.SetId(execution.ID_)
	entityImpl.SetProcessDefinitionId(*execution.ProcDefID_)
	entityImpl.SetProcessInstanceId(*execution.ProcInstID_)
	entityImpl.SetCurrentActivityId(cast.ToString(execution.ActID_))
	entityImpl.SetParentId(cast.ToString(execution.ParentID_))
	entityImpl.SetMultiInstanceRoot(cast.ToBool(execution.IsMiRoot_))
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
			AbstractEntity: AbstractEntity{
				Id: item.ID_,
			},
			BusinessKey:         cast.ToString(item.BusinessKey_),
			ProcessInstanceId:   cast.ToString(item.ProcInstID_),
			CurrentActivityId:   cast.ToString(item.ActID_),
			BusinessStatus:      cast.ToString(item.BusinessStatus_),
			Suspended:           cast.ToBool(item.SuspensionState_),
			ProcessDefinitionId: cast.ToString(item.ProcDefID_),
			StartUserId:         cast.ToString(item.StartUserID_),
			StartTime:           *item.StartTime_,
			CallbackId:          cast.ToString(item.CallbackID_),
			CallbackType:        cast.ToString(item.CallbackType_),
			ReferenceId:         cast.ToString(item.ReferenceID_),
			ReferenceType:       cast.ToString(item.ReferenceType_),
			TenantId:            item.TenantID_,
			ParentId:            cast.ToString(item.ParentID_),
			isMultiInstanceRoot: cast.ToBool(item.IsMiRoot_),
		}
	})
	return result, nil
}

func (executionEntityManager ExecutionEntityManager) CreateExecution(executionEntity *ExecutionEntity) error {
	actRuExecution := model.ActRuExecution{
		Rev_:                    lo.ToPtr(cast.ToInt32(1)),
		ProcInstID_:             &executionEntity.ProcessInstanceId,
		ParentID_:               &executionEntity.ParentId,
		ProcDefID_:              &executionEntity.ProcessDefinitionId,
		RootProcInstID_:         &executionEntity.ProcessInstanceId,
		ActID_:                  &executionEntity.CurrentActivityId,
		IsActive_:               lo.ToPtr(true),
		IsConcurrent_:           lo.ToPtr(false),
		IsScope_:                lo.ToPtr(false),
		IsEventScope_:           lo.ToPtr(false),
		IsMiRoot_:               &executionEntity.isMultiInstanceRoot,
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

	executionDataManager := datamanager.GetExecutionDataManager()
	if err := executionDataManager.Insert(&actRuExecution); err != nil {
		return err
	}

	executionEntity.SetId(actRuExecution.ID_)
	return nil
}

func (executionEntityManager ExecutionEntityManager) DeleteExecution(executionId string) error {
	executionDataManager := datamanager.GetExecutionDataManager()
	return executionDataManager.Delete(executionId)
}

func (executionEntityManager ExecutionEntityManager) DeleteRelatedDataForExecution(executionId string, deleteReason *string) error {
	variableDataManager := datamanager.GetVariableDataManager()
	if err := variableDataManager.DeleteByExecutionId(executionId); err != nil {
		return err
	}

	taskEntities, err := taskEntityManager.FindByExecutionId(executionId)
	if err != nil {
		return err
	}

	lo.ForEachWhile(taskEntities, func(item TaskEntity, index int) (goon bool) {
		if err = taskEntityManager.DeleteTask(item, deleteReason); err != nil {
			return false
		}
		return true
	})

	if err != nil {
		return err
	}

	return nil
}

func (executionEntityManager ExecutionEntityManager) DeleteChildExecutions(children []ExecutionEntity, deleteReason *string) error {
	var err error
	length := len(children)

	for i := length - 1; i >= 0; i-- {
		execution := children[i]

		if err = historicActivityInstanceEntityManager.RecordActEndByExecutionIdAndActId(execution.GetExecutionId(), execution.GetCurrentActivityId(), deleteReason); err != nil {
			return err
		}

		if err = executionEntityManager.DeleteRelatedDataForExecution(execution.GetExecutionId(), deleteReason); err != nil {
			return err
		}

		if err = executionEntityManager.DeleteExecution(execution.GetExecutionId()); err != nil {
			return err
		}
	}

	return nil
}

func (executionEntityManager ExecutionEntityManager) GetTopKValueFromChildExecutions(children []ExecutionEntity, variableName string, k int) ([]interface{}, error) {
	executionIds := lo.Map[ExecutionEntity, string](children, func(item ExecutionEntity, index int) string {
		return item.GetExecutionId()
	})

	variables, err := variableEntityManager.GetVariablesByParentIdAndVariableName(executionIds, variableName)
	if err != nil {
		return nil, err
	}

	if len(variables) == 0 {
		return nil, nil
	}

	variableValues := lo.Map[variable.Variable, interface{}](variables, func(item variable.Variable, index int) interface{} {
		return item.GetValue()
	})

	return GetTopKValues(variableValues, k), nil
}

func (executionEntityManager ExecutionEntityManager) CollectChildren(executionId string) ([]ExecutionEntity, error) {
	var result []ExecutionEntity

	children, err := executionEntityManager.List(execution.ListRequest{
		ListCommonRequest: request.ListCommonRequest{
			Start: 0,
			Size:  math.MaxInt32,
		},
		ParentId: executionId,
	})
	if err != nil {
		return nil, err
	}
	if len(children) == 0 {
		return nil, nil
	}
	result = append(result, children...)

	var grandChildren []ExecutionEntity

	lo.ForEachWhile(children, func(item ExecutionEntity, index int) (goon bool) {
		grandChildren, err = executionEntityManager.CollectChildren(item.GetExecutionId())
		if err != nil {
			return false
		}
		if len(grandChildren) == 0 {
			return true
		}
		result = append(result, grandChildren...)
		return true
	})
	if err != nil {
		return nil, err
	}
	return result, nil
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
