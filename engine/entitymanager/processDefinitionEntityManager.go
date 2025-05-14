package entitymanager

import (
	"errors"

	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/dto/procdef"
	"github.com/go-cinderella/cinderella-engine/engine/internal/datamanager"
	"github.com/go-cinderella/cinderella-engine/engine/internal/model"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/wubin1989/gorm"
	"github.com/wubin1989/gorm/clause"
)

type ProcessDefinitionEntityManager struct {
}

func (processDefinitionEntityManager ProcessDefinitionEntityManager) getProcessDefinitionEntity(processDefinition model.ActReProcdef) ProcessDefinitionEntity {
	processDefinitionEntity := ProcessDefinitionEntity{}
	processDefinitionEntity.SetId(processDefinition.ID_)
	processDefinitionEntity.SetName(cast.ToString(processDefinition.Name_))
	processDefinitionEntity.SetDescription(cast.ToString(processDefinition.Description_))
	processDefinitionEntity.SetKey(processDefinition.Key_)
	processDefinitionEntity.SetVersion(int(processDefinition.Version_))
	processDefinitionEntity.SetCreatedBy(cast.ToString(processDefinition.CreatedBy_))
	processDefinitionEntity.SetCreatedByName(cast.ToString(processDefinition.CreatedByName_))
	processDefinitionEntity.SetCategory(cast.ToString(processDefinition.Category_))
	processDefinitionEntity.SetDeploymentId(cast.ToString(processDefinition.DeploymentID_))
	processDefinitionEntity.SetResourceName(cast.ToString(processDefinition.ResourceName_))
	return processDefinitionEntity
}

func (processDefinitionEntityManager ProcessDefinitionEntityManager) FindProcessDefinitionById(processDefinitionId string) (ProcessDefinitionEntity, error) {
	processDefinition := model.ActReProcdef{}
	processDefinitionDataManager := datamanager.GetProcessDefinitionDataManager()
	if err := processDefinitionDataManager.FindById(processDefinitionId, &processDefinition); err != nil {
		return ProcessDefinitionEntity{}, err
	}
	processDefinitionEntity := processDefinitionEntityManager.getProcessDefinitionEntity(processDefinition)
	return processDefinitionEntity, nil
}

func (processDefinitionEntityManager ProcessDefinitionEntityManager) FindLatestProcessDefinitionByKey(processDefinitionKey string) (ProcessDefinitionEntity, error) {
	processDefinitionDataManager := datamanager.GetProcessDefinitionDataManager()
	processDefinition, err := processDefinitionDataManager.FindDeployedProcessDefinitionByKey(processDefinitionKey)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error("FindDeployedProcessDefinitionByKey err :", err)
		return ProcessDefinitionEntity{}, err
	}
	processDefinitionEntity := processDefinitionEntityManager.getProcessDefinitionEntity(processDefinition)
	return processDefinitionEntity, nil
}

func (processDefinitionEntityManager ProcessDefinitionEntityManager) FindResourceEntityByProcessDefinitionById(processDefinitionId string) (ResourceEntity, error) {
	processDefinitionEntity, err := processDefinitionEntityManager.FindProcessDefinitionById(processDefinitionId)
	if err != nil {
		return ResourceEntity{}, err
	}

	deploymentId := processDefinitionEntity.GetDeploymentId()
	resourceName := processDefinitionEntity.GetResourceName()
	resourceEntity, err := GetResourceEntityManager().FindResourceByDeploymentIdAndResourceName(deploymentId, resourceName)
	if err != nil {
		return ResourceEntity{}, err
	}

	return resourceEntity, nil
}

func (processDefinitionEntityManager ProcessDefinitionEntityManager) FindByDeploymentId(deploymentId string) (ProcessDefinitionEntity, error) {
	processDefinitionDataManager := datamanager.GetProcessDefinitionDataManager()
	processDefinition, err := processDefinitionDataManager.FindDeployedProcessDefinitionByDeploymentId(deploymentId)
	if err != nil {
		return ProcessDefinitionEntity{}, err
	}

	processDefinitionEntity := processDefinitionEntityManager.getProcessDefinitionEntity(processDefinition)

	deploymentEntity, err := GetDeploymentEntityManager().FindById(deploymentId)
	if err != nil {
		return ProcessDefinitionEntity{}, err
	}
	processDefinitionEntity.ResourceContent = deploymentEntity.GetResources().GetBytes()

	return processDefinitionEntity, nil
}

func (processDefinitionEntityManager ProcessDefinitionEntityManager) Insert(processDefinition *model.ActReProcdef) (ProcessDefinitionEntity, error) {
	var oldProcessDefinitions []model.ActReProcdef
	procDefQuery := contextutil.GetQuery().ActReProcdef
	if err := procDefQuery.Clauses(clause.Locking{Strength: "UPDATE"}).Where(procDefQuery.Key.Eq(processDefinition.Key_)).Where(procDefQuery.DeploymentID.IsNotNull()).Order(procDefQuery.Version.Desc()).Limit(1).Fetch(&oldProcessDefinitions); err != nil {
		return ProcessDefinitionEntity{}, err
	}

	var version int32

	if len(oldProcessDefinitions) > 0 {
		version = oldProcessDefinitions[0].Version_
	}

	version++

	processDefinition.Version_ = version

	processDefinitionDataManager := datamanager.GetProcessDefinitionDataManager()
	if err := processDefinitionDataManager.Insert(processDefinition); err != nil {
		return ProcessDefinitionEntity{}, err
	}

	processDefinitionEntity := processDefinitionEntityManager.getProcessDefinitionEntity(*processDefinition)
	return processDefinitionEntity, nil
}

func (processDefinitionEntityManager ProcessDefinitionEntityManager) List(listRequest procdef.ListRequest) (result []ProcessDefinitionEntity, total int32, err error) {
	processDefinitionDataManager := datamanager.GetProcessDefinitionDataManager()
	processDefinitions, total, err := processDefinitionDataManager.List(listRequest)
	if err != nil {
		return nil, 0, err
	}

	result = lo.Map(processDefinitions, func(processDefinition datamanager.ProcdefDTO, _ int) ProcessDefinitionEntity {
		entity := processDefinitionEntityManager.getProcessDefinitionEntity(processDefinition.ActReProcdef)
		entity.DeployTime = processDefinition.DeployTime_
		return entity
	})
	return result, total, nil
}
