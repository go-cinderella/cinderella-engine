package entitymanager

import (
	"github.com/go-cinderella/cinderella-engine/engine/datamanager"
	"github.com/go-cinderella/cinderella-engine/engine/dto/historicactinst"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/go-cinderella/cinderella-engine/engine/model"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"time"
)

type HistoricActivityInstanceEntityManager struct {
	DefaultHistoryManager
}

func (historicActivityInstanceEntityManager HistoricActivityInstanceEntityManager) RecordActivityStart(entity delegate.DelegateExecution) error {
	now := time.Now().UTC()
	historicActinst := model.ActHiActinst{}
	historicActinst.ExecutionID_ = entity.GetExecutionId()
	historicActinst.ProcDefID_ = entity.GetProcessDefinitionId()
	historicActinst.ProcInstID_ = entity.GetProcessInstanceId()
	historicActinst.ActID_ = entity.GetCurrentActivityId()
	historicActinst.StartTime_ = now
	if entity.GetCurrentFlowElement() != nil {
		historicActinst.ActName_ = lo.ToPtr(entity.GetCurrentFlowElement().GetName())
		historicActinst.ActType_ = parseActivityType(entity.GetCurrentFlowElement())
	}
	actinstDataManager := datamanager.GetHistoricActinstDataManager()
	return actinstDataManager.Insert(&historicActinst)
}

func (historicActivityInstanceEntityManager HistoricActivityInstanceEntityManager) RecordSequenceFlow(entity delegate.DelegateExecution) error {
	now := time.Now().UTC()
	historicActinst := model.ActHiActinst{}
	historicActinst.ProcDefID_ = entity.GetProcessDefinitionId()
	historicActinst.ProcInstID_ = entity.GetProcessInstanceId()
	historicActinst.ActID_ = entity.GetCurrentActivityId()
	historicActinst.StartTime_ = now
	historicActinst.EndTime_ = &now
	historicActinst.Duration_ = lo.ToPtr(int64(0))
	if entity.GetCurrentFlowElement() != nil {
		historicActinst.ActName_ = lo.ToPtr(entity.GetCurrentFlowElement().GetName())
		historicActinst.ActType_ = parseActivityType(entity.GetCurrentFlowElement())
	}
	actinstDataManager := datamanager.GetHistoricActinstDataManager()
	return actinstDataManager.Insert(&historicActinst)
}

func (historicActivityInstanceEntityManager HistoricActivityInstanceEntityManager) RecordActEndByExecutionId(element delegate.FlowElement, entity delegate.DelegateExecution, deleteReason *string) error {
	actinstDataManager := datamanager.GetHistoricActinstDataManager()
	err := actinstDataManager.RecordActEndByExecutionId(entity.GetExecutionId(), element.GetId(), deleteReason)
	return err
}

func (historicActivityInstanceEntityManager HistoricActivityInstanceEntityManager) MigrateProcDefID(oldProcessDefinitionEntity, newProcessDefinitionEntity ProcessDefinitionEntity) error {
	actinstDataManager := datamanager.GetHistoricActinstDataManager()
	err := actinstDataManager.MigrateProcDefID(oldProcessDefinitionEntity.GetId(), newProcessDefinitionEntity.GetId())
	return err
}

func (historicActivityInstanceEntityManager HistoricActivityInstanceEntityManager) MigrateAct(processDefinitionEntity ProcessDefinitionEntity, oldActivityId string, newActivityId string, newName string, newType string) error {
	actinstDataManager := datamanager.GetHistoricActinstDataManager()
	err := actinstDataManager.MigrateAct(processDefinitionEntity.GetId(), oldActivityId, newActivityId, newName, newType)
	return err
}

func (historicActivityInstanceEntityManager HistoricActivityInstanceEntityManager) RecordActEndByExecutionIdAndActId(executionId, activityId string, deleteReason *string) error {
	actinstDataManager := datamanager.GetHistoricActinstDataManager()
	err := actinstDataManager.RecordActEndByExecutionId(executionId, activityId, deleteReason)
	return err
}

func parseActivityType(element delegate.FlowElement) string {
	return element.GetHandlerType()
}

func (historicActivityInstanceEntityManager HistoricActivityInstanceEntityManager) RecordTaskId(task model.ActRuTask) error {
	actinstDataManager := datamanager.GetHistoricActinstDataManager()
	err := actinstDataManager.RecordTaskId(task)
	return err
}

func (historicActivityInstanceEntityManager HistoricActivityInstanceEntityManager) RecordActEndByProcessInstanceId(processInstanceId string, deleteReason *string) error {
	actinstDataManager := datamanager.GetHistoricActinstDataManager()
	err := actinstDataManager.RecordActEndByProcessInstanceId(processInstanceId, deleteReason)
	return err
}

func (historicActivityInstanceEntityManager HistoricActivityInstanceEntityManager) List(listRequest historicactinst.ListRequest) ([]HistoricActivityInstanceEntity, error) {
	historicActinstDataManager := datamanager.GetHistoricActinstDataManager()
	historicActinsts, err := historicActinstDataManager.List(listRequest)
	if err != nil {
		return nil, err
	}
	result := lo.Map[model.ActHiActinst, HistoricActivityInstanceEntity](historicActinsts, func(item model.ActHiActinst, index int) HistoricActivityInstanceEntity {
		return HistoricActivityInstanceEntity{
			AbstractEntity: AbstractEntity{
				Id: item.ID_,
			},
			ActivityId:              item.ActID_,
			ActivityName:            cast.ToString(item.ActName_),
			ActivityType:            item.ActType_,
			ProcessDefinitionId:     item.ProcDefID_,
			ProcessInstanceId:       item.ProcInstID_,
			ExecutionId:             item.ExecutionID_,
			TaskId:                  cast.ToString(item.TaskID_),
			CalledProcessInstanceId: cast.ToString(item.CallProcInstID_),
			Assignee:                cast.ToString(item.Assignee_),
			StartTime:               item.StartTime_,
			EndTime:                 item.EndTime_,
			DurationInMillis:        cast.ToInt(item.Duration_),
			TenantId:                item.TenantID_,
			DeleteReason:            item.DeleteReason_,
			BusinessResult:          item.BusinessResult_,
			BusinessParameter:       item.BusinessParameter_,
		}
	})
	return result, nil
}

func (historicActivityInstanceEntityManager HistoricActivityInstanceEntityManager) DeleteByExecutionId(entity delegate.DelegateExecution) error {
	actinstDataManager := datamanager.GetHistoricActinstDataManager()
	return actinstDataManager.DeleteByExecutionId(entity.GetExecutionId())
}

func (historicActivityInstanceEntityManager HistoricActivityInstanceEntityManager) DeleteByProcessInstanceId(processInstanceId, actId string) error {
	actinstDataManager := datamanager.GetHistoricActinstDataManager()
	return actinstDataManager.DeleteByProcessInstanceId(processInstanceId, actId)
}

func (historicActivityInstanceEntityManager HistoricActivityInstanceEntityManager) RecordBusinessDataByExecutionId(entity delegate.DelegateExecution, businessParameter, businessResult string) error {
	actinstDataManager := datamanager.GetHistoricActinstDataManager()
	err := actinstDataManager.RecordBusinessDataByExecutionId(entity.GetExecutionId(), &businessParameter, &businessResult)
	return err
}
