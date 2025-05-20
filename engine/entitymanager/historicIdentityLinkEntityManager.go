package entitymanager

import (
	"github.com/go-cinderella/cinderella-engine/engine/datamanager"
	"github.com/go-cinderella/cinderella-engine/engine/model"
	"github.com/samber/lo"
)

type HistoricIdentityLinkEntityManager struct {
}

func (historicIdentityLinkEntityManager HistoricIdentityLinkEntityManager) SelectByProcessInstanceId(processInstanceId string) ([]*HistoricIdentityLinkEntity, error) {
	historicIdentityLinkDataManager := datamanager.GetHistoricIdentityLinkDataManager()

	hiIdentitylinks, err := historicIdentityLinkDataManager.SelectByProcessInstanceId(processInstanceId)
	if err != nil {
		return nil, err
	}

	result := lo.Map[model.ActHiIdentitylink, *HistoricIdentityLinkEntity](hiIdentitylinks, func(item model.ActHiIdentitylink, index int) *HistoricIdentityLinkEntity {
		return &HistoricIdentityLinkEntity{
			AbstractEntity: AbstractEntity{
				Id: item.ID_,
			},
			GroupID:    item.GroupID_,
			Type:       item.Type_,
			UserID:     item.UserID_,
			TaskID:     item.TaskID_,
			ProcInstID: item.ProcInstID_,
		}
	})

	return result, nil
}
