package datamanager

import (
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/dto/execution"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/go-cinderella/cinderella-engine/engine/internal/datamanager/abstract"
	"github.com/go-cinderella/cinderella-engine/engine/internal/model"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/unionj-cloud/toolkit/stringutils"
	"github.com/wubin1989/gen/field"
	"strings"
)

type ExecutionDataManager struct {
	abstract.DataManager
}

// 创建流程实例
func (executionDataManager *ExecutionDataManager) CreateProcessInstance(processInstance *model.ActRuExecution) error {
	if err := executionDataManager.Insert(processInstance); err != nil {
		log.Infoln("create processInstance err", err)
		return err
	}

	err := executionDataManager.createHistoricProcessInstance(processInstance)
	return err
}

// RecordBusinessStatus records business status of a process instance
func (executionDataManager *ExecutionDataManager) RecordBusinessStatus(processInstanceId, businessStatus string) error {
	executionQuery := contextutil.GetQuery().ActRuExecution
	updateExample := model.ActRuExecution{
		BusinessStatus_: &businessStatus,
	}
	_, err := executionQuery.Where(executionQuery.ID.Eq(processInstanceId)).Updates(updateExample)
	if err != nil {
		return err
	}

	historicProcessManager := GetHistoricProcessDataManager()
	err = historicProcessManager.RecordBusinessStatus(processInstanceId, businessStatus)
	return err
}

// 查询流程实例
func (executionDataManager *ExecutionDataManager) GetProcessInstance(processInstanceId string) (model.ActRuExecution, error) {
	instance := model.ActRuExecution{}
	if err := executionDataManager.FindById(processInstanceId, &instance); err != nil {
		log.Infoln("create processInstance err", err)
		return model.ActRuExecution{}, nil
	}
	return instance, nil
}

// DeleteProcessInstance 删除流程实例
func (executionDataManager ExecutionDataManager) DeleteProcessInstance(execution delegate.DelegateExecution, deleteReason *string) error {
	processInstanceId := execution.GetProcessInstanceId()
	executionQuery := contextutil.GetQuery().ActRuExecution
	_, err := executionQuery.Where(executionQuery.ID.Eq(processInstanceId)).Or(executionQuery.ParentID.Eq(processInstanceId)).Delete()
	if err != nil {
		return err
	}

	currentActivityId := execution.GetCurrentActivityId()
	historicProcessManager := GetHistoricProcessDataManager()
	err = historicProcessManager.MarkEnded(processInstanceId, &currentActivityId, deleteReason)
	return err
}

func (executionDataManager *ExecutionDataManager) createHistoricProcessInstance(processInstance *model.ActRuExecution) (err error) {
	historicProcess := model.ActHiProcinst{}
	historicProcess.ID_ = processInstance.ID_
	historicProcess.ProcInstID_ = processInstance.ID_
	historicProcess.ProcDefID_ = cast.ToString(processInstance.ProcDefID_)
	historicProcess.TenantID_ = processInstance.TenantID_
	historicProcess.StartTime_ = cast.ToTime(processInstance.StartTime_)
	historicProcess.Name_ = processInstance.Name_
	historicProcess.BusinessKey_ = processInstance.BusinessKey_
	historicProcess.StartUserID_ = processInstance.StartUserID_
	historicProcess.StartActID_ = processInstance.StartActID_
	historicProcess.SuperProcessInstanceID_ = processInstance.SuperExec_
	historicProcessManager := GetHistoricProcessDataManager()
	return historicProcessManager.Insert(&historicProcess)
}

func (executionDataManager ExecutionDataManager) List(listRequest execution.ListRequest) ([]model.ActRuExecution, error) {
	executionQ := contextutil.GetQuery().ActRuExecution
	do := executionQ.Where()

	if stringutils.IsNotEmpty(listRequest.ProcessInstanceId) {
		do = do.Where(executionQ.ProcInstID.Eq(listRequest.ProcessInstanceId))
	}

	if stringutils.IsNotEmpty(listRequest.ParentId) {
		do = do.Where(executionQ.ParentID.Eq(listRequest.ParentId))
	}

	if listRequest.ChildOnly != nil && *listRequest.ChildOnly {
		do = do.Where(executionQ.ParentID.IsNotNull())
	}

	commonRequest := listRequest.ListCommonRequest
	if stringutils.IsNotEmpty(commonRequest.Sort) {
		var sortField field.Field
		switch commonRequest.Sort {
		case "start":
			sortField = field.Field(executionQ.StartTime)
		default:
			sortField = field.NewField((&model.ActRuExecution{}).TableName(), commonRequest.Sort)
		}

		if stringutils.IsNotEmpty(commonRequest.Order) && strings.ToLower(commonRequest.Order) == "desc" {
			do = do.Order(sortField.Desc())
		} else {
			do = do.Order(sortField)
		}
	}

	var result []model.ActRuExecution
	if err := do.Offset(commonRequest.Start).Limit(commonRequest.Size).Fetch(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (executionDataManager *ExecutionDataManager) MigrateProcessInstance(oldProcDefId, newProcDefId, startActId string) error {
	executionQuery := contextutil.GetQuery().ActRuExecution
	updateExample := model.ActRuExecution{
		ProcDefID_:  &newProcDefId,
		StartActID_: &startActId,
	}
	_, err := executionQuery.Where(executionQuery.ProcDefID.Eq(oldProcDefId), executionQuery.ParentID.IsNull()).Updates(updateExample)
	return err
}

func (executionDataManager *ExecutionDataManager) MigrateExecutionProcDefID(oldProcDefId, newProcDefId string) error {
	executionQuery := contextutil.GetQuery().ActRuExecution
	updateExample := model.ActRuExecution{
		ProcDefID_: &newProcDefId,
	}
	_, err := executionQuery.Where(executionQuery.ProcDefID.Eq(oldProcDefId), executionQuery.ParentID.IsNotNull()).Updates(updateExample)
	return err
}

func (executionDataManager *ExecutionDataManager) MigrateExecutionActID(procDefId, oldActivityId, newActivityId string) error {
	executionQuery := contextutil.GetQuery().ActRuExecution
	updateExample := model.ActRuExecution{
		ActID_: &newActivityId,
	}
	_, err := executionQuery.Where(executionQuery.ProcDefID.Eq(procDefId), executionQuery.ActID.Eq(oldActivityId), executionQuery.ParentID.IsNotNull()).Updates(updateExample)
	return err
}

func (executionDataManager *ExecutionDataManager) MigrateProcessInstanceBusinessStatus(procDefId, oldActivityId, newActivityId string) error {
	if oldActivityId == newActivityId {
		return nil
	}
	executionQuery := contextutil.GetQuery().ActRuExecution
	updateExample := model.ActRuExecution{
		BusinessStatus_: &newActivityId,
	}
	_, err := executionQuery.Where(executionQuery.ProcDefID.Eq(procDefId), executionQuery.BusinessStatus.Eq(oldActivityId), executionQuery.ParentID.IsNull()).Updates(updateExample)
	return err
}

func (executionDataManager *ExecutionDataManager) UpdateActive(id string, isActive bool) error {
	executionQuery := contextutil.GetQuery().ActRuExecution
	updateExample := model.ActRuExecution{
		IsActive_: &isActive,
	}
	_, err := executionQuery.Where(executionQuery.ID.Eq(id)).Updates(updateExample)
	return err
}

func (executionDataManager *ExecutionDataManager) IsActive(id string) (bool, error) {
	executionQuery := contextutil.GetQuery().ActRuExecution
	instance := model.ActRuExecution{}
	if err := executionQuery.Where(executionQuery.ID.Eq(id)).Fetch(&instance); err != nil {
		return false, err
	}

	if instance.IsActive_ != nil {
		return *instance.IsActive_, nil
	}
	return false, nil
}
