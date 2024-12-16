package entitymanager

import (
	"github.com/go-cinderella/cinderella-engine/engine/internal/datamanager"
	"github.com/go-cinderella/cinderella-engine/engine/internal/model"
	"github.com/jinzhu/copier"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"time"
)

type IdentityLinkEntityManager struct {
}

func (identityLinkEntityManager IdentityLinkEntityManager) CreateIdentityLink(identityLink model.ActRuIdentitylink) (err error) {
	linkDataManager := datamanager.GetIdentityLinkDataManager()
	linkDataManager.Insert(&identityLink)
	if err != nil {
		log.Infoln("Create IdentityLink Err ", err)
		return err
	}
	err = identityLinkEntityManager.createHistoricIdentityLink(identityLink)
	return err
}

func (identityLinkEntityManager IdentityLinkEntityManager) createHistoricIdentityLink(identityLink model.ActRuIdentitylink) (err error) {
	var historicIdentityLink model.ActHiIdentitylink
	copier.Copy(&historicIdentityLink, &identityLink)
	historicIdentityLink.CreateTime_ = lo.ToPtr(time.Now().UTC())
	historicIdentityLinkDataManager := datamanager.GetHistoricIdentityLinkDataManager()
	err = historicIdentityLinkDataManager.Insert(&historicIdentityLink)
	return err
}

func (identityLinkEntityManager IdentityLinkEntityManager) DeleteIdentityLinksByTaskId(taskId string) error {
	linkDataManager := datamanager.GetIdentityLinkDataManager()
	links, err := linkDataManager.SelectByTaskId(taskId)
	if err != nil {
		return err
	}
	for _, link := range links {
		if err = linkDataManager.Delete(link.ID_); err != nil {
			return err
		}
	}
	return nil
}

func (identityLinkEntityManager IdentityLinkEntityManager) DeleteByProcessInstanceId(processInstanceId string) error {
	linkDataManager := datamanager.GetIdentityLinkDataManager()
	err := linkDataManager.DeleteByProcessInstanceId(processInstanceId)
	return err
}

func (identityLinkEntityManager IdentityLinkEntityManager) SelectByProcessInstanceId(processInstanceId string) ([]*IdentityLinkEntity, error) {
	identityLinkDataManager := datamanager.GetIdentityLinkDataManager()

	ruIdentitylinks, err := identityLinkDataManager.SelectByProcessInstanceId(processInstanceId)
	if err != nil {
		return nil, err
	}

	result := lo.Map[model.ActRuIdentitylink, *IdentityLinkEntity](ruIdentitylinks, func(item model.ActRuIdentitylink, index int) *IdentityLinkEntity {
		return &IdentityLinkEntity{
			AbstractEntity: AbstractEntity{
				Id: item.ID_,
			},
			Rev:        item.Rev_,
			GroupID:    item.GroupID_,
			Type:       item.Type_,
			UserID:     item.UserID_,
			TaskID:     item.TaskID_,
			ProcInstID: item.ProcInstID_,
			ProcDefID:  item.ProcDefID_,
		}
	})

	return result, nil
}
