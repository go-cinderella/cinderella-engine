package cmd

import (
	"context"
	"fmt"
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/constant"
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/go-cinderella/cinderella-engine/engine/impl/handler"
	"github.com/go-cinderella/cinderella-engine/engine/utils"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

var _ engine.Command = (*MoveActivityCmd)(nil)

type MoveActivityCmd struct {
	Ctx               context.Context
	TargetActivityId  string
	Variables         map[string]interface{}
	UserId            string
	ProcessInstanceId string
	Transactional     bool
}

func (moveActivityCmd MoveActivityCmd) IsTransactional() bool {
	return moveActivityCmd.Transactional
}

var deleteReason = "Change activity to %s"

func (moveActivityCmd MoveActivityCmd) Execute(commandContext engine.Context) (result interface{}, err error) {
	executionEntityManager := entitymanager.GetExecutionEntityManager()
	taskEntityManager := entitymanager.GetTaskEntityManager()
	var taskEntities []entitymanager.TaskEntity
	taskEntities, err = taskEntityManager.FindByProcessInstanceId(moveActivityCmd.ProcessInstanceId)
	if err != nil {
		return nil, err
	}

	task := taskEntities[0]
	var executionEntity entitymanager.ExecutionEntity
	executionEntity, err = executionEntityManager.FindById(task.GetExecutionId())
	if err != nil {
		return nil, err
	}

	processUtils := utils.ProcessDefinitionUtil{}
	var process model.Process
	process, err = processUtils.GetProcess(executionEntity.GetProcessDefinitionId())
	if err != nil {
		return nil, err
	}

	currentTask := process.GetFlowElement(task.GetTaskDefineKey())
	targetFlowElement := process.GetFlowElement(moveActivityCmd.TargetActivityId)
	sequenceFlow := newSequenceFlow(targetFlowElement, currentTask.GetId())
	currentTask.SetOutgoing([]delegate.FlowElement{sequenceFlow})
	executionEntity.SetCurrentFlowElement(currentTask)

	if moveActivityCmd.Variables != nil {
		if err = executionEntity.SetVariable(&executionEntity, moveActivityCmd.Variables); err != nil {
			return nil, err
		}
	}

	userTask, ok := currentTask.(*model.UserTask)
	if ok {
		taskListeners := userTask.ExtensionElements.TaskListener
		for _, listener := range taskListeners {
			if listener.EventType == constant.TASK_TYPE_MOVED {
				err = handler.PerformTaskListener(&task, userTask.Name, task.GetCurrentActivityId())
				if err != nil {
					return nil, err
				}
			}
		}
	}

	historicActivityInstanceEntityManager := entitymanager.GetHistoricActivityInstanceEntityManager()
	reason := fmt.Sprintf(deleteReason, moveActivityCmd.TargetActivityId)
	lo.ForEachWhile(taskEntities, func(item entitymanager.TaskEntity, index int) (goon bool) {
		assignTaskCmd := NewAssignTaskCmd(moveActivityCmd.Ctx, item.GetId(), &moveActivityCmd.UserId)
		if _, err = assignTaskCmd.Execute(commandContext); err != nil {
			return false
		}

		if err = taskEntityManager.DeleteTask(item, &reason); err != nil {
			return false
		}

		if err = executionEntityManager.DeleteExecution(item.GetExecutionId()); err != nil {
			return false
		}

		if err = historicActivityInstanceEntityManager.RecordActEndByExecutionIdAndActId(item.GetExecutionId(), item.GetTaskDefineKey(), &reason); err != nil {
			return false
		}

		return true
	})

	contextutil.GetAgendaFromContext(commandContext).PlanTriggerExecutionOperation(&executionEntity)
	return nil, nil
}

func newSequenceFlow(target delegate.FlowElement, sourceRef string) delegate.FlowElement {
	sequenceFlow := model.SequenceFlow{}
	id, _ := uuid.NewUUID()
	sequenceFlow.Id = id.String()
	sequenceFlow.SourceRef = sourceRef
	sequenceFlow.TargetRef = target.GetId()
	sequenceFlow.FlowNode = model.FlowNode{
		BaseHandlerType: sequenceFlow,
		DefaultBaseElement: model.DefaultBaseElement{
			Id: sequenceFlow.Id,
		},
	}

	sequenceFlow.SetTargetFlowElement(target)
	return &sequenceFlow
}

func (moveActivityCmd MoveActivityCmd) Context() context.Context {
	return moveActivityCmd.Ctx
}
