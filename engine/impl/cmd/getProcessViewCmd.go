package cmd

import (
	"context"
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/constant"
	"github.com/go-cinderella/cinderella-engine/engine/dto/historicactinst"
	"github.com/go-cinderella/cinderella-engine/engine/dto/historicprocess"
	"github.com/go-cinderella/cinderella-engine/engine/dto/request"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/samber/lo"
	"math"
)

var _ engine.Command = (*GetProcessViewCmd)(nil)

type GetProcessViewCmd struct {
	Ctx               context.Context
	ProcessInstanceId string
	Transactional     bool
}

func (g GetProcessViewCmd) IsTransactional() bool {
	return g.Transactional
}

type ProcessViewDTO struct {
	HighLightedFlows     []string
	ActiveActivityIds    []string
	HisActiveActivityIds []string
	ModelXml             string
	ModelName            string
}

func (g GetProcessViewCmd) Execute(commandContext engine.Context) (interface{}, error) {
	var activeActivityIds []string
	var hisActiveActivityIds []string
	var highLightedFlows []string

	historicActivityInstanceEntityManager := entitymanager.GetHistoricActivityInstanceEntityManager()
	histActInsts, err := historicActivityInstanceEntityManager.List(historicactinst.ListRequest{
		ListCommonRequest: request.ListCommonRequest{
			Size: math.MaxInt32,
		},
		ProcessInstanceId: g.ProcessInstanceId,
	})
	if err != nil {
		return nil, err
	}

	if len(histActInsts) > 0 {
		historicSquenceFlows := lo.Filter(histActInsts, func(item entitymanager.HistoricActivityInstanceEntity, index int) bool {
			return item.ActivityType == constant.ELEMENT_SEQUENCE_FLOW
		})
		for _, historicActivityInstance := range historicSquenceFlows {
			highLightedFlows = append(highLightedFlows, historicActivityInstance.ActivityId)
		}
		histActInsts = lo.Filter(histActInsts, func(item entitymanager.HistoricActivityInstanceEntity, index int) bool {
			return item.ActivityType != constant.ELEMENT_SEQUENCE_FLOW
		})
		for _, historicActivityInstance := range histActInsts {
			hisActiveActivityIds = append(hisActiveActivityIds, historicActivityInstance.ActivityId)
		}

		activeActinsts := lo.Filter(histActInsts, func(item entitymanager.HistoricActivityInstanceEntity, index int) bool {
			return item.EndTime == nil || item.ActivityType == constant.ELEMENT_EVENT_END
		})

		for _, historicActivityInstance := range activeActinsts {
			activeActivityIds = append(activeActivityIds, historicActivityInstance.ActivityId)
		}
	}

	historicProcessInstanceEntityManager := entitymanager.GetHistoricProcessInstanceEntityManager()
	histProcInsts, err := historicProcessInstanceEntityManager.List(historicprocess.ListRequest{
		ListCommonRequest: request.ListCommonRequest{
			Size: math.MaxInt32,
		},
		ProcessInstanceId: g.ProcessInstanceId,
	})
	if err != nil {
		return nil, err
	}
	if len(histProcInsts) == 0 {
		return nil, nil
	}

	processInstance := histProcInsts[0]
	processDefinitionId := processInstance.ProcessDefinitionId
	modelName := processInstance.ProcessDefinitionName

	processDefinitionEntityManager := entitymanager.GetProcessDefinitionEntityManager()
	resourceEntity, err := processDefinitionEntityManager.FindResourceEntityByProcessDefinitionById(processDefinitionId)
	if err != nil {
		return nil, err
	}

	bpmnContent := string(resourceEntity.GetBytes())
	return &ProcessViewDTO{
		HighLightedFlows:     highLightedFlows,
		ActiveActivityIds:    activeActivityIds,
		ModelXml:             bpmnContent,
		ModelName:            modelName,
		HisActiveActivityIds: hisActiveActivityIds,
	}, nil
}

func (g GetProcessViewCmd) Context() context.Context {
	return g.Ctx
}
