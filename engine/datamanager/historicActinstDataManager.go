package datamanager

import (
	"github.com/go-cinderella/cinderella-engine/engine/constant"
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/dto/historicactinst"
	"github.com/go-cinderella/cinderella-engine/engine/internal/datamanager/abstract"
	"github.com/go-cinderella/cinderella-engine/engine/internal/model"
	"github.com/samber/lo"
	"github.com/unionj-cloud/toolkit/stringutils"
	"github.com/wubin1989/gen/field"
	"github.com/wubin1989/gorm/clause"
	"strings"
	"time"
)

type HistoricActinstDataManager struct {
	abstract.DataManager
}

func (historicActinstManager HistoricActinstDataManager) RecordActEndByExecutionId(executionId string, actId string, deleteReason *string) error {
	var historicActs []*model.ActHiActinst
	hiActInstQuery := contextutil.GetQuery().ActHiActinst
	if err := hiActInstQuery.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where(hiActInstQuery.ActID.Eq(actId)).
		Where(hiActInstQuery.ExecutionID.Eq(executionId)).
		Where(hiActInstQuery.EndTime.IsNull()).
		Fetch(&historicActs); err != nil {
		return err
	}
	if len(historicActs) == 0 {
		return nil
	}

	historicAct := historicActs[0]

	now := time.Now().UTC()
	start := historicAct.StartTime_.UTC()
	duration := int64(now.Sub(start)) / constant.DurationUnit

	historicAct.EndTime_ = &now
	historicAct.Duration_ = &duration
	historicAct.DeleteReason_ = deleteReason

	_, err := hiActInstQuery.Where(hiActInstQuery.ID.Eq(historicAct.ID_)).Updates(historicAct)
	return err
}

func (historicActinstManager HistoricActinstDataManager) RecordTaskId(task model.ActRuTask) error {
	hiActInstQuery := contextutil.GetQuery().ActHiActinst

	update := make(map[string]interface{})
	update["assignee_"] = task.Assignee_
	update["task_id_"] = &task.ID_

	_, err := hiActInstQuery.Where(hiActInstQuery.ActID.Eq(*task.TaskDefKey_)).Where(hiActInstQuery.ExecutionID.Eq(*task.ExecutionID_)).Where(hiActInstQuery.EndTime.IsNull()).Updates(update)
	return err
}

func (historicActinstManager HistoricActinstDataManager) MigrateProcDefID(oldProcDefId, newProcDefId string) error {
	hiActInstQuery := contextutil.GetQuery().ActHiActinst
	updateExample := model.ActHiActinst{
		ProcDefID_: newProcDefId,
	}

	_, err := hiActInstQuery.Where(hiActInstQuery.ProcDefID.Eq(oldProcDefId)).Updates(&updateExample)
	return err
}

func (historicActinstManager HistoricActinstDataManager) ChangeTaskAssignee(taskId string, userId *string) (err error) {
	update := make(map[string]interface{})
	update["assignee_"] = userId

	hiActInstQuery := contextutil.GetQuery().ActHiActinst
	_, err = hiActInstQuery.Where(hiActInstQuery.TaskID.Eq(taskId)).Updates(update)
	return err
}

func (historicActinstManager HistoricActinstDataManager) RecordActEndByProcessInstanceId(processInstanceId string, deleteReason *string) error {
	var historicActs []*model.ActHiActinst
	hiActInstQuery := contextutil.GetQuery().ActHiActinst
	if err := hiActInstQuery.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where(hiActInstQuery.ProcInstID.Eq(processInstanceId)).
		Where(hiActInstQuery.EndTime.IsNull()).
		Fetch(&historicActs); err != nil {
		return err
	}

	now := time.Now().UTC()
	var err error

	lo.ForEachWhile(historicActs, func(historicAct *model.ActHiActinst, index int) (goon bool) {
		start := historicAct.StartTime_.UTC()
		duration := int64(now.Sub(start)) / constant.DurationUnit

		historicAct.EndTime_ = &now
		historicAct.Duration_ = &duration
		historicAct.DeleteReason_ = deleteReason

		_, err = hiActInstQuery.Where(hiActInstQuery.ID.Eq(historicAct.ID_)).Updates(historicAct)
		if err != nil {
			return false
		}
		return true
	})

	return err
}

func (historicActinstManager HistoricActinstDataManager) List(listRequest historicactinst.ListRequest) ([]model.ActHiActinst, error) {
	actHiActinstQ := contextutil.GetQuery().ActHiActinst
	do := actHiActinstQ.Where()

	if stringutils.IsNotEmpty(listRequest.ProcessInstanceId) {
		do = do.Where(actHiActinstQ.ProcInstID.Eq(listRequest.ProcessInstanceId))
	}

	if listRequest.Finished != nil {
		if *listRequest.Finished {
			do = do.Where(actHiActinstQ.EndTime.IsNotNull())
		} else {
			do = do.Where(actHiActinstQ.EndTime.IsNull())
		}
	}

	if stringutils.IsNotEmpty(listRequest.ActivityType) {
		if strings.Contains(listRequest.ActivityType, ",") {
			do = do.Where(actHiActinstQ.ActType.In(stringutils.Split(listRequest.ActivityType, ",")...))
		} else {
			do = do.Where(actHiActinstQ.ActType.Eq(listRequest.ActivityType))
		}
	}

	if stringutils.IsNotEmpty(listRequest.ActivityId) {
		if strings.Contains(listRequest.ActivityId, ",") {
			do = do.Where(actHiActinstQ.ActID.In(stringutils.Split(listRequest.ActivityId, ",")...))
		} else {
			do = do.Where(actHiActinstQ.ActID.Eq(listRequest.ActivityId))
		}
	}

	commonRequest := listRequest.ListCommonRequest
	if stringutils.IsNotEmpty(commonRequest.Sort) {
		var sortField field.Field
		switch commonRequest.Sort {
		case "start":
			sortField = field.Field(actHiActinstQ.StartTime)
		default:
			sortField = field.NewField((&model.ActHiActinst{}).TableName(), commonRequest.Sort)
		}

		if stringutils.IsNotEmpty(commonRequest.Order) && strings.ToLower(commonRequest.Order) == "desc" {
			do = do.Order(sortField.Desc())
		} else {
			do = do.Order(sortField)
		}
	}

	var result []model.ActHiActinst
	if err := do.Offset(commonRequest.Start).Limit(commonRequest.Size).Fetch(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (historicActinstManager HistoricActinstDataManager) DeleteByExecutionId(executionId string) error {
	hiActInstQuery := contextutil.GetQuery().ActHiActinst
	_, err := hiActInstQuery.Where(hiActInstQuery.ExecutionID.Eq(executionId)).Delete()
	return err
}

func (historicActinstManager HistoricActinstDataManager) DeleteByProcessInstanceId(processInstanceId, actId string) error {
	hiActInstQuery := contextutil.GetQuery().ActHiActinst
	_, err := hiActInstQuery.Where(hiActInstQuery.ProcInstID.Eq(processInstanceId), hiActInstQuery.ActID.Eq(actId)).Delete()
	return err
}

func (historicActinstManager HistoricActinstDataManager) MigrateAct(procDefId, oldActivityId, newActivityId, newName, newType string) error {
	hiActInstQuery := contextutil.GetQuery().ActHiActinst
	updateExample := model.ActHiActinst{
		ActID_:   newActivityId,
		ActName_: &newName,
		ActType_: newType,
	}

	_, err := hiActInstQuery.Where(hiActInstQuery.ProcDefID.Eq(procDefId), hiActInstQuery.ActID.Eq(oldActivityId)).Updates(&updateExample)
	return err
}

func (historicActinstManager HistoricActinstDataManager) RecordBusinessDataByExecutionId(executionId string, businessParameter, businessResult *string) error {
	hiActInstQuery := contextutil.GetQuery().ActHiActinst
	updateExample := model.ActHiActinst{
		BusinessResult_:    businessResult,
		BusinessParameter_: businessParameter,
	}

	_, err := hiActInstQuery.Where(hiActInstQuery.ExecutionID.Eq(executionId)).Updates(&updateExample)
	return err
}
