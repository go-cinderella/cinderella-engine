package datamanager

import (
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/db"
	"github.com/go-cinderella/cinderella-engine/engine/dto/task"
	"github.com/go-cinderella/cinderella-engine/engine/internal/datamanager/abstract"
	"github.com/go-cinderella/cinderella-engine/engine/internal/model"
	log "github.com/sirupsen/logrus"
	"github.com/unionj-cloud/toolkit/stringutils"
	gormgen "github.com/wubin1989/gen"
	"github.com/wubin1989/gen/field"
	"strings"
)

type TaskDataManager struct {
	abstract.DataManager
}

func (taskManager TaskDataManager) GetById(id int64) (model.ActRuTask, error) {
	task := model.ActRuTask{}
	err := db.DB().Where(`"id_" = ?`, id).First(&task).Error
	if err != nil {
		log.Error("find bu id err:", err)
	}
	return task, err
}

func (taskManager TaskDataManager) FindByProcessInstanceId(processInstanceId string) ([]model.ActRuTask, error) {
	var tasks []model.ActRuTask
	taskQuery := contextutil.GetQuery().ActRuTask
	err := taskQuery.Where(taskQuery.ProcInstID.Eq(processInstanceId)).Fetch(&tasks)
	return tasks, err
}

func (taskManager TaskDataManager) FindByExecutionId(executionId string) ([]model.ActRuTask, error) {
	var tasks []model.ActRuTask
	taskQuery := contextutil.GetQuery().ActRuTask
	err := taskQuery.Where(taskQuery.ExecutionID.Eq(executionId)).Fetch(&tasks)
	return tasks, err
}

func (taskManager TaskDataManager) MigrateProcDefID(oldProcDefId, newProcDefId string) error {
	taskQuery := contextutil.GetQuery().ActRuTask
	updateExample := model.ActRuTask{
		ProcDefID_: &newProcDefId,
	}

	_, err := taskQuery.Where(taskQuery.ProcDefID.Eq(oldProcDefId)).Updates(&updateExample)
	return err
}

func (taskManager TaskDataManager) MigrateNameAndTaskDefKey(procDefId, oldActivityId, newActivityId, newName string) error {
	taskQuery := contextutil.GetQuery().ActRuTask
	updateExample := model.ActRuTask{
		Name_:       &newName,
		TaskDefKey_: &newActivityId,
	}

	_, err := taskQuery.Where(taskQuery.ProcDefID.Eq(procDefId), taskQuery.TaskDefKey.Eq(oldActivityId)).Updates(&updateExample)
	return err
}

func (taskManager TaskDataManager) QueryUndoTask(userId, groupId string) (taskResult []task.TaskDTO, err error) {
	taskResult = make([]task.TaskDTO, 0)
	taskQuery := contextutil.GetQuery().ActRuTask
	identityLinkQuery := contextutil.GetQuery().ActRuIdentitylink
	queryBuilder := taskQuery.Select(taskQuery.ALL, identityLinkQuery.UserID, identityLinkQuery.GroupID).
		LeftJoin(identityLinkQuery, identityLinkQuery.TaskID.EqCol(taskQuery.ID))
	if stringutils.IsNotEmpty(userId) {
		queryBuilder = queryBuilder.Where(identityLinkQuery.UserID.Eq(userId))
	}
	if stringutils.IsNotEmpty(groupId) {
		queryBuilder = queryBuilder.Where(identityLinkQuery.GroupID.Eq(groupId))
	}
	err = queryBuilder.Fetch(&taskResult)
	if err != nil {
		return taskResult, err
	}
	return taskResult, nil
}

func (taskManager TaskDataManager) ChangeTaskAssignee(task model.ActRuTask) error {
	taskQuery := contextutil.GetQuery().ActRuTask

	update := make(map[string]any)
	update["assignee_"] = task.Assignee_
	update["claim_time_"] = task.ClaimTime_

	if _, err := taskQuery.Where(taskQuery.ID.Eq(task.ID_)).Updates(update); err != nil {
		return err
	}
	return nil
}

func (taskManager TaskDataManager) List(listRequest task.ListRequest) ([]model.ActRuTask, error) {
	taskQuery := contextutil.GetQuery().ActRuTask
	do := taskQuery.Where()

	if stringutils.IsNotEmpty(listRequest.ProcessInstanceId) {
		do = do.Where(taskQuery.ProcInstID.Eq(listRequest.ProcessInstanceId))
	}

	if stringutils.IsNotEmpty(listRequest.TaskDefinitionKey) {
		do = do.Where(taskQuery.TaskDefKey.Eq(listRequest.TaskDefinitionKey))
	}

	if len(listRequest.TaskDefinitionKeys) > 0 {
		do = do.Where(taskQuery.TaskDefKey.In(listRequest.TaskDefinitionKeys...))
	}

	if stringutils.IsNotEmpty(listRequest.CandidateOrAssigned) || len(listRequest.CandidateGroupIn) > 0 {
		linkQuery := contextutil.GetQuery().ActRuIdentitylink

		if stringutils.IsNotEmpty(listRequest.CandidateOrAssigned) {
			subQuery := linkQuery.Select(linkQuery.ID).Where(
				linkQuery.TaskID.EqCol(taskQuery.ID),
				linkQuery.Type.Eq("candidate"),
				linkQuery.UserID.Eq(listRequest.CandidateOrAssigned),
			)
			do = do.Where(taskQuery.Where(taskQuery.Assignee.Eq(listRequest.CandidateOrAssigned)).Or(taskQuery.Where(gormgen.Exists(subQuery))))
		}

		if len(listRequest.CandidateGroupIn) > 0 {
			subQuery := linkQuery.Select(linkQuery.ID).Where(
				linkQuery.TaskID.EqCol(taskQuery.ID),
				linkQuery.Type.Eq("candidate"),
				linkQuery.GroupID.In(listRequest.CandidateGroupIn...),
			)
			do = do.Where(gormgen.Exists(subQuery))
		}
	}

	commonRequest := listRequest.ListCommonRequest
	if stringutils.IsNotEmpty(commonRequest.Sort) {
		var sortField field.Field
		switch commonRequest.Sort {
		case "start":
			sortField = field.Field(taskQuery.CreateTime)
		default:
			sortField = field.NewField((&model.ActRuTask{}).TableName(), commonRequest.Sort)
		}

		if stringutils.IsNotEmpty(commonRequest.Order) && strings.ToLower(commonRequest.Order) == "desc" {
			do = do.Order(sortField.Desc())
		} else {
			do = do.Order(sortField)
		}
	}

	var result []model.ActRuTask
	if err := do.Offset(commonRequest.Start).Limit(commonRequest.Size).Fetch(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (taskManager TaskDataManager) GetUniqTaskDefKeys(listRequest task.ListRequest) ([]model.ActRuTask, error) {
	taskQuery := contextutil.GetQuery().ActRuTask
	do := taskQuery.Where()

	if stringutils.IsNotEmpty(listRequest.ProcessInstanceId) {
		do = do.Where(taskQuery.ProcInstID.Eq(listRequest.ProcessInstanceId))
	}

	if stringutils.IsNotEmpty(listRequest.TaskDefinitionKey) {
		do = do.Where(taskQuery.TaskDefKey.Eq(listRequest.TaskDefinitionKey))
	}

	if stringutils.IsNotEmpty(listRequest.CandidateOrAssigned) || len(listRequest.CandidateGroupIn) > 0 {
		linkQuery := contextutil.GetQuery().ActRuIdentitylink

		if stringutils.IsNotEmpty(listRequest.CandidateOrAssigned) {
			subQuery := linkQuery.Select(linkQuery.ID).Where(
				linkQuery.TaskID.EqCol(taskQuery.ID),
				linkQuery.Type.Eq("candidate"),
				linkQuery.UserID.Eq(listRequest.CandidateOrAssigned),
			)
			do = do.Where(taskQuery.Where(taskQuery.Assignee.Eq(listRequest.CandidateOrAssigned)).Or(taskQuery.Where(gormgen.Exists(subQuery))))
		}

		if len(listRequest.CandidateGroupIn) > 0 {
			subQuery := linkQuery.Select(linkQuery.ID).Where(
				linkQuery.TaskID.EqCol(taskQuery.ID),
				linkQuery.Type.Eq("candidate"),
				linkQuery.GroupID.In(listRequest.CandidateGroupIn...),
			)
			do = do.Where(gormgen.Exists(subQuery))
		}
	}

	commonRequest := listRequest.ListCommonRequest
	if stringutils.IsNotEmpty(commonRequest.Sort) {
		var sortField field.Field
		switch commonRequest.Sort {
		case "start":
			sortField = field.Field(taskQuery.CreateTime)
		default:
			sortField = field.NewField((&model.ActRuTask{}).TableName(), commonRequest.Sort)
		}

		if stringutils.IsNotEmpty(commonRequest.Order) && strings.ToLower(commonRequest.Order) == "desc" {
			do = do.Order(sortField.Desc())
		} else {
			do = do.Order(sortField)
		}
	}

	var result []model.ActRuTask
	if err := do.Distinct(taskQuery.TaskDefKey, taskQuery.ProcDefID).Offset(commonRequest.Start).Limit(commonRequest.Size).Fetch(&result); err != nil {
		return nil, err
	}

	return result, nil
}
