package datamanager

import (
	"strings"
	"time"

	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/datamanager/abstract"
	"github.com/go-cinderella/cinderella-engine/engine/dto/procdef"
	"github.com/go-cinderella/cinderella-engine/engine/model"
	"github.com/unionj-cloud/toolkit/stringutils"
	"github.com/wubin1989/gen/field"
)

type ProcessDefinitionDataManager struct {
	abstract.DataManager
}

func (define ProcessDefinitionDataManager) FindDeployedProcessDefinitionByKey(key string) (model.ActReProcdef, error) {
	processDefinition := model.ActReProcdef{}
	procDefQuery := contextutil.GetQuery().ActReProcdef
	err := procDefQuery.Where(procDefQuery.Key.Eq(key)).Where(procDefQuery.DeploymentID.IsNotNull()).Order(procDefQuery.Version.Desc()).Fetch(&processDefinition)
	return processDefinition, err
}

func (define ProcessDefinitionDataManager) FindDeployedProcessDefinitionByDeploymentId(deploymentId string) (model.ActReProcdef, error) {
	processDefinition := model.ActReProcdef{}
	procDefQuery := contextutil.GetQuery().ActReProcdef
	err := procDefQuery.Where(procDefQuery.DeploymentID.Eq(deploymentId)).Fetch(&processDefinition)
	return processDefinition, err
}

type ProcdefDTO struct {
	model.ActReProcdef
	DeployTime_ time.Time `gorm:"column:deploy_time_;type:timestamp without time zone" json:"deploy_time_"`
}

func (define ProcessDefinitionDataManager) List(listRequest procdef.ListRequest) (result []ProcdefDTO, total int32, err error) {
	deploymentQ := contextutil.GetQuery().ActReDeployment
	procdefQ := contextutil.GetQuery().ActReProcdef
	do := procdefQ.Select(procdefQ.ALL, deploymentQ.DeployTime).
		LeftJoin(deploymentQ, deploymentQ.ID.EqCol(procdefQ.DeploymentID)).
		Where()

	if stringutils.IsNotEmpty(listRequest.ProcessDefinitionId) {
		do = do.Where(procdefQ.ID.Eq(listRequest.ProcessDefinitionId))
	}

	if stringutils.IsNotEmpty(listRequest.ProcessDefinitionKey) {
		do = do.Where(procdefQ.Key.Eq(listRequest.ProcessDefinitionKey))
	}

	if len(listRequest.ProcessDefinitionKeyIn) > 0 {
		do = do.Where(procdefQ.Key.In(listRequest.ProcessDefinitionKeyIn...))
	}

	t, err := do.Count()
	if err != nil {
		return nil, 0, err
	}

	commonRequest := listRequest.ListCommonRequest
	if stringutils.IsNotEmpty(commonRequest.Sort) {
		var sortField field.Field
		switch commonRequest.Sort {
		case "deploy_time":
			sortField = field.Field(deploymentQ.DeployTime)
		case "version":
			sortField = field.Field(procdefQ.Version)
		default:
			sortField = field.NewField((&model.ActReProcdef{}).TableName(), commonRequest.Sort)
		}

		if stringutils.IsNotEmpty(commonRequest.Order) && strings.ToLower(commonRequest.Order) == "desc" {
			do = do.Order(sortField.Desc())
		} else {
			do = do.Order(sortField)
		}
	}

	if err := do.Offset(commonRequest.Start).Limit(commonRequest.Size).Fetch(&result); err != nil {
		return nil, 0, err
	}

	return result, int32(t), nil
}
