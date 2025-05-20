package entitymanager

import (
	"errors"
	"github.com/go-cinderella/cinderella-engine/engine/datamanager"
	"github.com/go-cinderella/cinderella-engine/engine/dto/historicprocess"
	"github.com/go-cinderella/cinderella-engine/engine/model"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/unionj-cloud/toolkit/stringutils"
)

type HistoricProcessInstanceEntityManager struct {
}

func (historicProcessInstanceEntityManager HistoricProcessInstanceEntityManager) FindById(historicProcessInstanceId string) (HistoricProcessInstanceEntity, error) {
	processDataManager := datamanager.GetHistoricProcessDataManager()
	historicProcessInst := model.ActHiProcinst{}
	if err := processDataManager.FindById(historicProcessInstanceId, &historicProcessInst); err != nil {
		return HistoricProcessInstanceEntity{}, err
	}

	if stringutils.IsEmpty(historicProcessInst.ID_) {
		return HistoricProcessInstanceEntity{}, errors.New("historic process instance not found")
	}

	historicProcessInstanceEntity := toHistoricProcessInstanceEntity(datamanager.ActHiProcinstDTO{
		ActHiProcinst: historicProcessInst,
	})
	return historicProcessInstanceEntity, nil
}

func (historicProcessInstanceEntityManager HistoricProcessInstanceEntityManager) List(listRequest historicprocess.ListRequest) ([]HistoricProcessInstanceEntity, error) {
	historicProcessDataManager := datamanager.GetHistoricProcessDataManager()
	hiProcinsts, err := historicProcessDataManager.List(listRequest)
	if err != nil {
		return nil, err
	}
	result := lo.Map[datamanager.ActHiProcinstDTO, HistoricProcessInstanceEntity](hiProcinsts, func(item datamanager.ActHiProcinstDTO, index int) HistoricProcessInstanceEntity {
		return toHistoricProcessInstanceEntity(item)
	})
	return result, nil
}

func toHistoricProcessInstanceEntity(item datamanager.ActHiProcinstDTO) HistoricProcessInstanceEntity {
	return HistoricProcessInstanceEntity{
		AbstractEntity: AbstractEntity{
			Id: item.ID_,
		},
		BusinessKey:            cast.ToString(item.BusinessKey_),
		BusinessStatus:         cast.ToString(item.BusinessStatus_),
		ProcessDefinitionId:    cast.ToString(item.ProcDefID_),
		StartTime:              lo.ToPtr(item.StartTime_),
		EndTime:                item.EndTime_,
		DurationInMillis:       cast.ToInt(item.Duration_),
		StartUserId:            cast.ToString(item.StartUserID_),
		StartActivityId:        cast.ToString(item.StartActID_),
		EndActivityId:          cast.ToString(item.EndActID_),
		DeleteReason:           cast.ToString(item.DeleteReason_),
		SuperProcessInstanceId: cast.ToString(item.SuperProcessInstanceID_),
		TenantId:               cast.ToString(item.TenantID_),
		ProcessDefinitionName:  item.ProcessDefinitionName,
	}
}
