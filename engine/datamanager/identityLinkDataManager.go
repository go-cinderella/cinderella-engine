package datamanager

import (
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/datamanager/abstract"
	"github.com/go-cinderella/cinderella-engine/engine/errs"
	. "github.com/go-cinderella/cinderella-engine/engine/model"
	log "github.com/sirupsen/logrus"
)

type IdentityLinkDataManager struct {
	abstract.DataManager
}

func (identityLinkManager IdentityLinkDataManager) SelectByProcessInstanceId(processInstanceId string) ([]ActRuIdentitylink, error) {
	identityLink := make([]ActRuIdentitylink, 0)
	identityLinkQuery := contextutil.GetQuery().ActRuIdentitylink
	err := identityLinkQuery.Where(identityLinkQuery.ProcInstID.Eq(processInstanceId)).Fetch(&identityLink)
	return identityLink, err
}

func (identityLinkManager IdentityLinkDataManager) DeleteByProcessInstanceId(processInstanceId string) error {
	identityLinkQuery := contextutil.GetQuery().ActRuIdentitylink
	_, err := identityLinkQuery.Where(identityLinkQuery.ProcInstID.Eq(processInstanceId)).Delete()
	return err
}

func (identityLinkManager IdentityLinkDataManager) SelectByTaskId(taskId string) ([]ActRuIdentitylink, error) {
	identityLink := make([]ActRuIdentitylink, 0)
	identityLinkQuery := contextutil.GetQuery().ActRuIdentitylink
	err := identityLinkQuery.Where(identityLinkQuery.TaskID.Eq(taskId)).Fetch(&identityLink)
	if err != nil {
		log.Infoln("Select Variable err: ", err)
	}
	if identityLink != nil && len(identityLink) > 0 {
		return identityLink, nil
	}
	return identityLink, errs.CinderellaError{Code: "1001", Msg: "Not Find"}
}
