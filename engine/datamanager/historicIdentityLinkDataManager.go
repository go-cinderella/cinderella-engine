package datamanager

import (
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/datamanager/abstract"
	"github.com/go-cinderella/cinderella-engine/engine/model"
)

type HistoricIdentityLinkDataManager struct {
	abstract.DataManager
}

func (historicIdentityLinkDataManager HistoricIdentityLinkDataManager) SelectByProcessInstanceId(processInstanceId string) ([]model.ActHiIdentitylink, error) {
	identityLink := make([]model.ActHiIdentitylink, 0)
	identityLinkQuery := contextutil.GetQuery().ActHiIdentitylink
	err := identityLinkQuery.Where(identityLinkQuery.ProcInstID.Eq(processInstanceId)).Fetch(&identityLink)
	return identityLink, err
}
