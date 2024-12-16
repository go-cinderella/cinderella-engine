package cmd

import (
	"context"
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/dto/historicprocess"
	"github.com/go-cinderella/cinderella-engine/engine/dto/request"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
)

var _ engine.Command = (*GetProcessStatusCmd)(nil)

type GetProcessStatusCmd struct {
	Ctx               context.Context
	ProcessInstanceId string
	Transactional     bool
}

func (g GetProcessStatusCmd) IsTransactional() bool {
	return g.Transactional
}

func (g GetProcessStatusCmd) Execute(commandContext engine.Context) (interface{}, error) {
	historicProcessInstanceEntityManager := entitymanager.GetHistoricProcessInstanceEntityManager()
	procInsts, err := historicProcessInstanceEntityManager.List(historicprocess.ListRequest{
		ListCommonRequest: request.ListCommonRequest{
			Size: 1,
		},
		ProcessInstanceId: g.ProcessInstanceId,
	})
	if err != nil {
		return "", err
	}

	if len(procInsts) == 0 {
		return "", err
	}
	processInstance := procInsts[0]
	return processInstance.BusinessStatus, nil
}

func (g GetProcessStatusCmd) Context() context.Context {
	return g.Ctx
}
