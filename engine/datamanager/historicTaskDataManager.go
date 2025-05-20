package datamanager

import (
	"github.com/go-cinderella/cinderella-engine/engine/constant"
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/dto/historictask"
	"github.com/go-cinderella/cinderella-engine/engine/internal/datamanager/abstract"
	"github.com/go-cinderella/cinderella-engine/engine/internal/model"
	"github.com/unionj-cloud/toolkit/stringutils"
	"github.com/wubin1989/gen/field"
	"github.com/wubin1989/gorm/clause"
	"strings"
	"time"
)

type HistoricTaskDataManager struct {
	abstract.DataManager
}

func (historicTaskManager HistoricTaskDataManager) RecordTaskEnd(id string, deleteReason *string) error {
	var historicTasks []*model.ActHiTaskinst
	actHiTaskinstQ := contextutil.GetQuery().ActHiTaskinst
	if err := actHiTaskinstQ.Clauses(clause.Locking{Strength: "UPDATE"}).Where(actHiTaskinstQ.ID.Eq(id)).Where(actHiTaskinstQ.EndTime.IsNull()).Fetch(&historicTasks); err != nil {
		return err
	}
	if len(historicTasks) == 0 {
		return nil
	}

	historicTask := historicTasks[0]

	now := time.Now().UTC()
	start := historicTask.StartTime_.UTC()
	duration := int64(now.Sub(start)) / constant.DurationUnit

	historicTask.EndTime_ = &now
	historicTask.DeleteReason_ = deleteReason
	historicTask.Duration_ = &duration

	_, err := actHiTaskinstQ.Where(actHiTaskinstQ.ID.Eq(id)).Updates(&historicTask)
	return err
}

func (historicTaskManager HistoricTaskDataManager) List(listRequest historictask.ListRequest) ([]model.ActHiTaskinst, error) {
	actHiTaskinstQ := contextutil.GetQuery().ActHiTaskinst
	do := actHiTaskinstQ.Where()

	if stringutils.IsNotEmpty(listRequest.ProcessInstanceId) {
		do = do.Where(actHiTaskinstQ.ProcInstID.Eq(listRequest.ProcessInstanceId))
	}

	if stringutils.IsNotEmpty(listRequest.TaskDefinitionKey) {
		do = do.Where(actHiTaskinstQ.TaskDefKey.Eq(listRequest.TaskDefinitionKey))
	}

	if listRequest.Finished != nil {
		if *listRequest.Finished {
			do = do.Where(actHiTaskinstQ.EndTime.IsNotNull())
		} else {
			do = do.Where(actHiTaskinstQ.EndTime.IsNull())
		}
	}

	commonRequest := listRequest.ListCommonRequest
	if stringutils.IsNotEmpty(commonRequest.Sort) {
		var sortField field.Field
		switch commonRequest.Sort {
		case "start":
			sortField = field.Field(actHiTaskinstQ.StartTime)
		default:
			sortField = field.NewField((&model.ActHiTaskinst{}).TableName(), commonRequest.Sort)
		}

		if stringutils.IsNotEmpty(commonRequest.Order) && strings.ToLower(commonRequest.Order) == "desc" {
			do = do.Order(sortField.Desc())
		} else {
			do = do.Order(sortField)
		}
	}

	var result []model.ActHiTaskinst
	if err := do.Offset(commonRequest.Start).Limit(commonRequest.Size).Fetch(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (historicTaskManager HistoricTaskDataManager) ChangeTaskAssignee(task model.ActRuTask) error {
	taskQuery := contextutil.GetQuery().ActHiTaskinst

	update := make(map[string]interface{})
	update["assignee_"] = task.Assignee_
	update["claim_time_"] = task.ClaimTime_

	_, err := taskQuery.Where(taskQuery.ID.Eq(task.ID_)).Updates(update)
	return err
}

func (historicTaskManager HistoricTaskDataManager) MigrateNameAndTaskDefKey(procDefId, oldActivityId, newActivityId, newName string) error {
	taskQuery := contextutil.GetQuery().ActHiTaskinst
	updateExample := model.ActHiTaskinst{
		Name_:       &newName,
		TaskDefKey_: &newActivityId,
	}

	_, err := taskQuery.Where(taskQuery.ProcDefID.Eq(procDefId), taskQuery.TaskDefKey.Eq(oldActivityId)).Updates(&updateExample)
	return err
}

func (historicTaskManager HistoricTaskDataManager) MigrateProcDefID(oldProcDefId, newProcDefId string) error {
	taskQuery := contextutil.GetQuery().ActHiTaskinst
	updateExample := model.ActHiTaskinst{
		ProcDefID_: &newProcDefId,
	}

	_, err := taskQuery.Where(taskQuery.ProcDefID.Eq(oldProcDefId)).Updates(&updateExample)
	return err
}
