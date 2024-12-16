package agenda

import (
	"github.com/go-cinderella/cinderella-engine/engine/dto/execution"
	"github.com/go-cinderella/cinderella-engine/engine/dto/request"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/utils"
	"github.com/samber/lo"
	"math"
)

type EvaluateConditionalEventsOperation struct {
	AbstractOperation
}

func (cont *EvaluateConditionalEventsOperation) Run() (err error) {
	processUtils := utils.ProcessDefinitionUtil{}
	var process model.Process
	process, err = processUtils.GetProcess(cont.Execution.GetProcessDefinitionId())
	if err != nil {
		return err
	}

	processInstanceId := cont.Execution.GetProcessInstanceId()

	executionEntityManager := entitymanager.GetExecutionEntityManager()
	childExecutions, err := executionEntityManager.List(execution.ListRequest{
		ListCommonRequest: request.ListCommonRequest{
			Size: math.MaxInt32,
		},
		ProcessInstanceId: processInstanceId,
		ChildOnly:         lo.ToPtr(true),
	})
	if err != nil {
		return err
	}

	lo.ForEachWhile(childExecutions, func(item entitymanager.ExecutionEntity, index int) (goon bool) {
		activityId := item.GetCurrentActivityId()
		currentFlowElement := process.GetFlowElement(activityId)
		intermediateCatchEvent, ok := currentFlowElement.(*model.IntermediateCatchEvent)
		if ok && intermediateCatchEvent.ConditionalEventDefinition != nil {
			item.SetCurrentFlowElement(currentFlowElement)
			cont.GetAgenda().PlanTriggerExecutionOperation(&item)
		}
		return true
	})

	return err
}
