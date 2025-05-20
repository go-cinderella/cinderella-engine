package entitymanager

import (
	"github.com/go-cinderella/cinderella-engine/engine/datamanager"
	"github.com/go-cinderella/cinderella-engine/engine/dto/historictask"
	"github.com/go-cinderella/cinderella-engine/engine/model"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type HistoricTaskInstanceEntityManager struct {
	DefaultHistoryManager
}

func (historicTaskInstanceEntityManager HistoricTaskInstanceEntityManager) RecordTaskEnd(taskId string, deleteReason *string) error {
	historicTaskDataManager := datamanager.GetHistoricTaskDataManager()
	err := historicTaskDataManager.RecordTaskEnd(taskId, deleteReason)
	return err
}

func (historicTaskInstanceEntityManager HistoricTaskInstanceEntityManager) MigrateProcDefID(oldProcessDefinitionEntity, newProcessDefinitionEntity ProcessDefinitionEntity) error {
	historicTaskDataManager := datamanager.GetHistoricTaskDataManager()
	err := historicTaskDataManager.MigrateProcDefID(oldProcessDefinitionEntity.GetId(), newProcessDefinitionEntity.GetId())
	return err
}

func (historicTaskInstanceEntityManager HistoricTaskInstanceEntityManager) List(listRequest historictask.ListRequest) ([]HistoricTaskInstanceEntity, error) {
	historicTaskDataManager := datamanager.GetHistoricTaskDataManager()
	hiTaskinsts, err := historicTaskDataManager.List(listRequest)
	if err != nil {
		return nil, err
	}
	result := lo.Map[model.ActHiTaskinst, HistoricTaskInstanceEntity](hiTaskinsts, func(item model.ActHiTaskinst, index int) HistoricTaskInstanceEntity {
		return HistoricTaskInstanceEntity{
			AbstractEntity: AbstractEntity{
				Id: item.ID_,
			},
			ProcessDefinitionId: cast.ToString(item.ProcDefID_),
			ProcessInstanceId:   cast.ToString(item.ProcInstID_),
			ExecutionId:         cast.ToString(item.ExecutionID_),
			Name:                cast.ToString(item.Name_),
			Description:         cast.ToString(item.Description_),
			DeleteReason:        cast.ToString(item.DeleteReason_),
			Owner:               cast.ToString(item.Owner_),
			Assignee:            cast.ToString(item.Assignee_),
			TaskDefinitionKey:   cast.ToString(item.TaskDefKey_),
			FormKey:             cast.ToString(item.FormKey_),
			Priority:            cast.ToInt(item.Priority_),
			ParentTaskId:        cast.ToString(item.ParentTaskID_),
			TenantId:            cast.ToString(item.TenantID_),
			Category:            cast.ToString(item.Category_),
			DurationInMillis:    cast.ToInt(item.Duration_),
			StartTime:           item.StartTime_,
			EndTime:             item.EndTime_,
			ClaimTime:           item.ClaimTime_,
			DueDate:             item.DueDate_,
		}
	})
	return result, nil
}
