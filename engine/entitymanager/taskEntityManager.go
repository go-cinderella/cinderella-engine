package entitymanager

import (
	"errors"
	"github.com/go-cinderella/cinderella-engine/engine/datamanager"
	"github.com/go-cinderella/cinderella-engine/engine/dto/task"
	"github.com/go-cinderella/cinderella-engine/engine/errs"
	"github.com/go-cinderella/cinderella-engine/engine/model"
	"github.com/jinzhu/copier"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/unionj-cloud/toolkit/zlogger"
	"github.com/wubin1989/gorm"
)

type TaskEntityManager struct {
}

func (taskEntityManager TaskEntityManager) FindById(id string) (TaskEntity, error) {
	task := model.ActRuTask{}
	taskDataManager := datamanager.GetTaskDataManager()
	err := taskDataManager.FirstById(id, &task)
	if err != nil {
		zlogger.Error().Err(err).Msgf("get task err: %s, taskId: %s", err, id)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return TaskEntity{}, errs.ErrTaskNotFound
		}
		return TaskEntity{}, errs.ErrInternalError
	}
	var taskEntity TaskEntity
	taskEntity.SetAssignee(task.Assignee_)
	taskEntity.SetStartTime(cast.ToTime(task.CreateTime_))
	taskEntity.SetProcessInstanceId(cast.ToString(task.ProcInstID_))
	taskEntity.SetId(task.ID_)
	taskEntity.SetTaskDefineKey(cast.ToString(task.TaskDefKey_))
	taskEntity.SetTaskDefineName(cast.ToString(task.Name_))
	taskEntity.SetExecutionId(cast.ToString(task.ExecutionID_))
	taskEntity.ProcessDefinitionId = cast.ToString(task.ProcDefID_)

	executionDataManager := datamanager.GetExecutionDataManager()
	execution := &model.ActRuExecution{}
	if err = executionDataManager.FirstById(*task.ExecutionID_, execution); err != nil {
		zlogger.Error().Err(err).Msgf("get execution err: %s, executionId: %s", err, *task.ExecutionID_)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return TaskEntity{}, errs.ErrExecutionNotFound
		}
		return TaskEntity{}, errs.ErrInternalError
	}
	executionEntity := ExecutionEntity{}
	executionEntity.SetId(execution.ID_)
	executionEntity.SetProcessDefinitionId(*execution.ProcDefID_)
	executionEntity.SetProcessInstanceId(*execution.ProcInstID_)

	taskEntity.ExecutionEntity = executionEntity
	return taskEntity, nil
}
func (taskEntityManager TaskEntityManager) DeleteTask(task TaskEntity, deleteReason *string) error {
	identityLinkManager := GetIdentityLinkManager()
	if err := identityLinkManager.DeleteIdentityLinksByTaskId(task.GetId()); err != nil {
		return err
	}

	taskDataManager := datamanager.GetTaskDataManager()
	if err := taskDataManager.Delete(task.GetId()); err != nil {
		return err
	}

	historicTaskManager := GetHistoricTaskInstanceEntityManager()
	err := historicTaskManager.RecordTaskEnd(task.GetId(), deleteReason)
	return err
}

func (taskEntityManager TaskEntityManager) FindByProcessInstanceId(processInstanceId string) ([]TaskEntity, error) {
	taskDataManager := datamanager.GetTaskDataManager()
	tasks, err := taskDataManager.FindByProcessInstanceId(processInstanceId)
	if err != nil {
		return nil, err
	}
	var taskEntitys []TaskEntity
	for _, task := range tasks {
		taskEntity := TaskEntity{}
		taskEntity.SetId(task.ID_)
		taskEntity.SetTaskDefineKey(cast.ToString(task.TaskDefKey_))
		taskEntity.SetExecutionId(cast.ToString(task.ExecutionID_))
		taskEntitys = append(taskEntitys, taskEntity)
	}
	return taskEntitys, nil
}

func (taskEntityManager TaskEntityManager) FindByExecutionId(executionId string) ([]TaskEntity, error) {
	taskDataManager := datamanager.GetTaskDataManager()
	tasks, err := taskDataManager.FindByExecutionId(executionId)
	if err != nil {
		return nil, err
	}
	var taskEntitys []TaskEntity
	for _, task := range tasks {
		taskEntity := TaskEntity{}
		taskEntity.SetId(task.ID_)
		taskEntity.SetTaskDefineKey(cast.ToString(task.TaskDefKey_))
		taskEntity.SetExecutionId(cast.ToString(task.ExecutionID_))
		taskEntitys = append(taskEntitys, taskEntity)
	}
	return taskEntitys, nil
}

func (taskEntityManager TaskEntityManager) InsertTask(taskEntity *TaskEntity) error {
	task := model.ActRuTask{
		Rev_:             lo.ToPtr(int32(1)),
		ExecutionID_:     &taskEntity.ExecutionId,
		ProcInstID_:      &taskEntity.ProcessInstanceId,
		ProcDefID_:       &taskEntity.ProcessDefinitionId,
		Name_:            &taskEntity.TaskDefineName,
		TaskDefKey_:      &taskEntity.TaskDefinitionKey,
		Assignee_:        taskEntity.Assignee,
		CreateTime_:      &taskEntity.StartTime,
		DueDate_:         taskEntity.DueDate,
		SuspensionState_: lo.ToPtr(int32(1)),
		TenantID_:        taskEntity.TenantId,
		FormKey_:         taskEntity.FormKey,
		ClaimTime_:       taskEntity.ClaimTime,
		IsCountEnabled_:  lo.ToPtr(true),
		VarCount_:        lo.ToPtr(int32(0)),
		IDLinkCount_:     lo.ToPtr(int32(0)),
		SubTaskCount_:    lo.ToPtr(int32(0)),
	}
	taskDataManager := datamanager.GetTaskDataManager()
	err := taskDataManager.Insert(&task)
	taskEntity.SetId(task.ID_)
	if err != nil {
		log.Error("create task err:", err)
		return err
	}
	err = taskEntityManager.recordTaskCreated(task)
	return err
}

func (taskEntityManager TaskEntityManager) recordTaskCreated(task model.ActRuTask) error {
	historicTaskManager := datamanager.GetHistoricTaskDataManager()
	historicTask := taskEntityManager.newHistoricTask(task)
	if err := historicTaskManager.Insert(&historicTask); err != nil {
		return err
	}

	entityManager := GetHistoricActivityInstanceEntityManager()
	err := entityManager.RecordTaskId(task)
	return err
}

func (taskEntityManager TaskEntityManager) newHistoricTask(task model.ActRuTask) model.ActHiTaskinst {
	var historicTask model.ActHiTaskinst
	copier.Copy(&historicTask, &task)
	historicTask.StartTime_ = *task.CreateTime_
	return historicTask
}

func (taskEntityManager TaskEntityManager) ChangeTaskAssignee(taskEntity TaskEntity, userId *string) error {
	taskModel := model.ActRuTask{
		ID_:        taskEntity.Id,
		Assignee_:  userId,
		ClaimTime_: taskEntity.ClaimTime,
	}

	taskDataManager := datamanager.GetTaskDataManager()
	if err := taskDataManager.ChangeTaskAssignee(taskModel); err != nil {
		return err
	}

	historicTaskDataManager := datamanager.GetHistoricTaskDataManager()
	if err := historicTaskDataManager.ChangeTaskAssignee(taskModel); err != nil {
		return err
	}

	historicActinstDataManager := datamanager.GetHistoricActinstDataManager()
	if err := historicActinstDataManager.ChangeTaskAssignee(taskEntity.Id, userId); err != nil {
		return err
	}

	// 暂时不用act_ru_actinst表，只用act_hi_actinst

	return nil
}

func (taskEntityManager TaskEntityManager) MigrateProcDefID(oldProcessDefinitionEntity, newProcessDefinitionEntity ProcessDefinitionEntity) error {
	taskDataManager := datamanager.GetTaskDataManager()
	err := taskDataManager.MigrateProcDefID(oldProcessDefinitionEntity.GetId(), newProcessDefinitionEntity.GetId())
	return err
}

func (taskEntityManager TaskEntityManager) MigrateNameAndTaskDefKey(processDefinitionEntity ProcessDefinitionEntity, oldActivityId string, newActivityId string, newName string) error {
	taskDataManager := datamanager.GetTaskDataManager()
	err := taskDataManager.MigrateNameAndTaskDefKey(processDefinitionEntity.GetId(), oldActivityId, newActivityId, newName)
	if err != nil {
		return err
	}

	historicTaskDataManager := datamanager.GetHistoricTaskDataManager()
	err = historicTaskDataManager.MigrateNameAndTaskDefKey(processDefinitionEntity.GetId(), oldActivityId, newActivityId, newName)
	return err
}

func (taskEntityManager TaskEntityManager) List(listRequest task.ListRequest) ([]TaskEntity, error) {
	taskDataManager := datamanager.GetTaskDataManager()
	taskinsts, err := taskDataManager.List(listRequest)
	if err != nil {
		return nil, err
	}
	result := lo.Map[model.ActRuTask, TaskEntity](taskinsts, func(item model.ActRuTask, index int) TaskEntity {
		return TaskEntity{
			AbstractEntity: AbstractEntity{
				Id: item.ID_,
			},
			Assignee:                  item.Assignee_,
			StartTime:                 *item.CreateTime_,
			TaskDefinitionKey:         *item.TaskDefKey_,
			TaskDefineName:            cast.ToString(item.Name_),
			ClaimTime:                 item.ClaimTime_,
			Owner:                     cast.ToString(item.Owner_),
			DelegationState:           cast.ToString(item.Delegation_),
			Description:               cast.ToString(item.Description_),
			DueDate:                   item.DueDate_,
			Priority:                  cast.ToInt(item.Priority_),
			Suspended:                 false,
			ScopeDefinitionId:         cast.ToString(item.ScopeDefinitionID_),
			ScopeId:                   cast.ToString(item.ScopeID_),
			SubScopeId:                cast.ToString(item.SubScopeID_),
			ScopeType:                 cast.ToString(item.ScopeType_),
			PropagatedStageInstanceId: cast.ToString(item.PropagatedStageInstID_),
			TenantId:                  item.TenantID_,
			Category:                  cast.ToString(item.Category_),
			FormKey:                   item.FormKey_,
			ParentTaskId:              cast.ToString(item.ParentTaskID_),
			ExecutionId:               cast.ToString(item.ExecutionID_),
			ProcessInstanceId:         cast.ToString(item.ProcInstID_),
			ProcessDefinitionId:       cast.ToString(item.ProcDefID_),
		}
	})
	return result, nil
}
