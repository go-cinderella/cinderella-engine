package entitymanager

import (
	"errors"
	"github.com/go-cinderella/cinderella-engine/engine/dto/historicprocess"
	"github.com/go-cinderella/cinderella-engine/engine/internal/datamanager"
	"github.com/go-cinderella/cinderella-engine/engine/internal/model"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/unionj-cloud/toolkit/stringutils"
	"time"
)

type HistoricProcessInstanceEntityManager struct {
}

func (historicProcessInstanceEntityManager HistoricProcessInstanceEntityManager) FindById(historicProcessInstanceId string) (HistoricProcessInstanceEntity, error) {
	processDataManager := datamanager.GetHistoricProcessDataManager()
	historicProcessInst := &model.ActHiProcinst{}
	if err := processDataManager.FindById(historicProcessInstanceId, historicProcessInst); err != nil {
		return HistoricProcessInstanceEntity{}, err
	}

	if stringutils.IsEmpty(historicProcessInst.ID_) {
		return HistoricProcessInstanceEntity{}, errors.New("historic process instance not found")
	}

	historicProcessInstanceEntity := HistoricProcessInstanceEntity{}
	historicProcessInstanceEntity.SetId(historicProcessInst.ID_)
	historicProcessInstanceEntity.SetProcessDefinitionId(historicProcessInst.ProcDefID_)
	historicProcessInstanceEntity.EndTime = historicProcessInst.EndTime_
	return historicProcessInstanceEntity, nil
}

func (historicProcessInstanceEntityManager HistoricProcessInstanceEntityManager) List(listRequest historicprocess.ListRequest) ([]HistoricProcessInstanceEntity, error) {
	historicProcessDataManager := datamanager.GetHistoricProcessDataManager()
	hiProcinsts, err := historicProcessDataManager.List(listRequest)
	if err != nil {
		return nil, err
	}
	result := lo.Map[datamanager.ActHiProcinstDTO, HistoricProcessInstanceEntity](hiProcinsts, func(item datamanager.ActHiProcinstDTO, index int) HistoricProcessInstanceEntity {
		return HistoricProcessInstanceEntity{
			AbstractEntity: AbstractEntity{
				Id: item.ID_,
			},
			BusinessKey:         cast.ToString(item.BusinessKey_),
			BusinessStatus:      cast.ToString(item.BusinessStatus_),
			ProcessDefinitionId: cast.ToString(item.ProcDefID_),
			StartTime:           lo.ToPtr(item.StartTime_),
			EndTime: lo.TernaryF[*time.Time](item.EndTime_ != nil, func() *time.Time {
				return lo.ToPtr(*item.EndTime_)
			}, func() *time.Time {
				return nil
			}),
			DurationInMillis:       cast.ToInt(item.Duration_),
			StartUserId:            cast.ToString(item.StartUserID_),
			StartActivityId:        cast.ToString(item.StartActID_),
			EndActivityId:          cast.ToString(item.EndActID_),
			DeleteReason:           cast.ToString(item.DeleteReason_),
			SuperProcessInstanceId: cast.ToString(item.SuperProcessInstanceID_),
			TenantId:               cast.ToString(item.TenantID_),
			ProcessDefinitionName:  item.ProcessDefinitionName,
		}
	})
	return result, nil
}
