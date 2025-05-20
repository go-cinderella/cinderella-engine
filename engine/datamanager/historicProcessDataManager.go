package datamanager

import (
	"github.com/go-cinderella/cinderella-engine/engine/constant"
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/datamanager/abstract"
	"github.com/go-cinderella/cinderella-engine/engine/dto/historicprocess"
	"github.com/go-cinderella/cinderella-engine/engine/model"
	"github.com/unionj-cloud/toolkit/stringutils"
	"github.com/wubin1989/gen/field"
	"github.com/wubin1989/gorm/clause"
	"strings"
	"time"
)

type HistoricProcessDataManager struct {
	abstract.DataManager
}

func (historicProcessManager HistoricProcessDataManager) MarkEnded(processInstanceId string, endActID *string, deleteReason *string) (err error) {
	var historicProcesses []*model.ActHiProcinst
	actHiProcinstQ := contextutil.GetQuery().ActHiProcinst
	if err := actHiProcinstQ.Clauses(clause.Locking{Strength: "UPDATE"}).Where(actHiProcinstQ.ID.Eq(processInstanceId)).
		Where(actHiProcinstQ.EndTime.IsNull()).Fetch(&historicProcesses); err != nil {
		return err
	}
	if len(historicProcesses) == 0 {
		return nil
	}

	historicProcess := historicProcesses[0]
	now := time.Now().UTC()
	start := historicProcess.StartTime_.UTC()
	duration := int64(now.Sub(start)) / constant.DurationUnit

	historicProcess.EndTime_ = &now
	historicProcess.Duration_ = &duration
	historicProcess.DeleteReason_ = deleteReason
	historicProcess.EndActID_ = endActID

	_, err = actHiProcinstQ.Where(actHiProcinstQ.ID.Eq(processInstanceId)).Updates(&historicProcess)
	return err
}

func (historicProcessManager HistoricProcessDataManager) RecordBusinessStatus(processInstanceId, businessStatus string) error {
	actHiProcinstQ := contextutil.GetQuery().ActHiProcinst
	updateExample := model.ActHiProcinst{
		BusinessStatus_: &businessStatus,
	}

	_, err := actHiProcinstQ.Where(actHiProcinstQ.ID.Eq(processInstanceId)).Updates(updateExample)
	return err
}

func (historicProcessManager HistoricProcessDataManager) Migrate(oldProcDefId, newProcDefId, startActId string) error {
	hiProcinstQ := contextutil.GetQuery().ActHiProcinst
	updateExample := model.ActHiProcinst{
		ProcDefID_:  newProcDefId,
		StartActID_: &startActId,
	}
	_, err := hiProcinstQ.Where(hiProcinstQ.ProcDefID.Eq(oldProcDefId)).Updates(updateExample)
	return err
}

func (historicProcessManager HistoricProcessDataManager) MigrateBusinessStatus(procDefId, oldActivityId, newActivityId string) error {
	if oldActivityId == newActivityId {
		return nil
	}
	hiProcinstQ := contextutil.GetQuery().ActHiProcinst
	updateExample := model.ActHiProcinst{
		BusinessStatus_: &newActivityId,
	}
	_, err := hiProcinstQ.Where(hiProcinstQ.ProcDefID.Eq(procDefId), hiProcinstQ.BusinessStatus.Eq(oldActivityId)).Updates(updateExample)
	return err
}

type ActHiProcinstDTO struct {
	model.ActHiProcinst
	ProcessDefinitionName string
}

func (historicProcessManager HistoricProcessDataManager) List(listRequest historicprocess.ListRequest) ([]ActHiProcinstDTO, error) {
	hiProcinstQ := contextutil.GetQuery().ActHiProcinst
	procdefQ := contextutil.GetQuery().ActReProcdef
	do := hiProcinstQ.Select(hiProcinstQ.ALL, procdefQ.Name.As("process_definition_name")).LeftJoin(procdefQ, procdefQ.ID.EqCol(hiProcinstQ.ProcDefID)).Where()

	if stringutils.IsNotEmpty(listRequest.ProcessInstanceId) {
		do = do.Where(hiProcinstQ.ID.Eq(listRequest.ProcessInstanceId))
	}

	commonRequest := listRequest.ListCommonRequest
	if stringutils.IsNotEmpty(commonRequest.Sort) {
		var sortField field.Field
		switch commonRequest.Sort {
		case "start":
			sortField = field.Field(hiProcinstQ.StartTime)
		default:
			sortField = field.NewField((&model.ActHiProcinst{}).TableName(), commonRequest.Sort)
		}

		if stringutils.IsNotEmpty(commonRequest.Order) && strings.ToLower(commonRequest.Order) == "desc" {
			do = do.Order(sortField.Desc())
		} else {
			do = do.Order(sortField)
		}
	}

	var result []ActHiProcinstDTO
	if err := do.Offset(commonRequest.Start).Limit(commonRequest.Size).Fetch(&result); err != nil {
		return nil, err
	}

	return result, nil
}
