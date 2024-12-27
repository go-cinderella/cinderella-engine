package cmd

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/dto/historicactinst"
	"github.com/go-cinderella/cinderella-engine/engine/dto/request"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
)

var json = sonic.ConfigDefault

var _ engine.Command = (*GetPipelineActivityInfoCmd)(nil)

type GetPipelineActivityInfoCmd struct {
	Ctx               context.Context
	ProcessInstanceId string
	ActivityId        string
	Transactional     bool
}

func (getPipelineActivityInfoCmd GetPipelineActivityInfoCmd) IsTransactional() bool {
	return getPipelineActivityInfoCmd.Transactional
}

func (getPipelineActivityInfoCmd GetPipelineActivityInfoCmd) Execute(commandContext engine.Context) (interface{}, error) {
	historicActivityInstanceEntityManager := entitymanager.GetHistoricActivityInstanceEntityManager()

	histActInsts, err := historicActivityInstanceEntityManager.List(historicactinst.ListRequest{
		ListCommonRequest: request.ListCommonRequest{
			Size:  1,
			Sort:  "start",
			Order: "desc",
		},
		ProcessInstanceId: getPipelineActivityInfoCmd.ProcessInstanceId,
		ActivityId:        getPipelineActivityInfoCmd.ActivityId,
	})
	if err != nil {
		return nil, err
	}

	if len(histActInsts) == 0 {
		return nil, nil
	}
	histActInst := histActInsts[0]

	var result historicactinst.PipelineActivityDTO

	if histActInst.BusinessParameter != nil {
		if err = json.UnmarshalFromString(*histActInst.BusinessParameter, &result); err != nil {
			return nil, err
		}
	}

	if histActInst.BusinessResult != nil {
		result.BusinessResult = *histActInst.BusinessResult
	}

	return result, nil
}

func (getPipelineActivityInfoCmd GetPipelineActivityInfoCmd) Context() context.Context {
	return getPipelineActivityInfoCmd.Ctx
}
