package cmd

import (
	"context"
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
)

var _ engine.Command = (*GetIdentityLinksForProcessInstanceCmd)(nil)

type GetIdentityLinksForProcessInstanceCmd struct {
	ProcessInstanceId string
	Ctx               context.Context
	Transactional     bool
}

func (g GetIdentityLinksForProcessInstanceCmd) IsTransactional() bool {
	return g.Transactional
}

func (g GetIdentityLinksForProcessInstanceCmd) Execute(commandContext engine.Context) (interface{}, error) {
	historicIdentityLinkEntityManager := entitymanager.GetHistoricIdentityLinkEntityManager()
	identityLinkEntities, err := historicIdentityLinkEntityManager.SelectByProcessInstanceId(g.ProcessInstanceId)
	return identityLinkEntities, err
}

func (g GetIdentityLinksForProcessInstanceCmd) Context() context.Context {
	return g.Ctx
}
