package cmd

import (
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/constant"
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/go-cinderella/cinderella-engine/engine/eventmanager"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/cmd/taskcmd"
	"github.com/go-cinderella/cinderella-engine/engine/impl/handler"
	"github.com/go-cinderella/cinderella-engine/engine/internal/errs"
	"github.com/go-cinderella/cinderella-engine/engine/utils"
	log "github.com/sirupsen/logrus"
)

var _ taskcmd.ITaskCmd = (*CompleteCmd)(nil)

type CompleteCmd struct {
	taskcmd.NeedsActiveTaskCmd
	ProcessVariables map[string]interface{}
	UserId           *string
}

func (completeCmd CompleteCmd) TaskExecute(commandContext engine.Context, entity entitymanager.TaskEntity) (interface{}, error) {
	err := completeCmd.executeTaskComplete(entity, commandContext)
	return entity, err
}

func (completeCmd CompleteCmd) executeTaskComplete(task entitymanager.TaskEntity, commandContext engine.Context) (err error) {
	if completeCmd.UserId != nil {
		assignTaskCmd := NewAssignTaskCmd(completeCmd.Ctx, completeCmd.TaskId, completeCmd.UserId)
		if _, err = assignTaskCmd.Execute(commandContext); err != nil {
			return err
		}
	}

	executionEntity, err := entitymanager.GetExecutionEntityManager().FindById(task.GetExecutionId())
	if err != nil {
		return err
	}

	processUtils := utils.ProcessDefinitionUtil{}
	process, err := processUtils.GetProcess(executionEntity.GetProcessDefinitionId())
	if err != nil {
		return err
	}

	currentTask := process.GetFlowElement(executionEntity.GetCurrentActivityId())

	executionEntity.SetCurrentFlowElement(currentTask)

	if err = executionEntity.SetVariable(&executionEntity, completeCmd.ProcessVariables); err != nil {
		return err
	}

	userTask, ok := currentTask.(*model.UserTask)
	if ok {
		taskListeners := userTask.ExtensionElements.TaskListener
		for _, listener := range taskListeners {
			if listener.EventType == constant.TASK_TYPE_COMPLETED {
				err = handler.PerformTaskListener(&task, userTask.Name, task.GetCurrentActivityId())
				if err != nil {
					return err
				}
			}
		}
	} else {
		log.Error("not task")
		return errs.CinderellaError{Code: "not task"}
	}

	// All properties set, now firing 'create' events
	entityEvent := eventmanager.CreateEntityEvent(eventmanager.TASK_COMPLETED, task)
	if err = eventmanager.GetEventDispatcher().DispatchEvent(entityEvent); err != nil {
		return
	}

	if err = entitymanager.GetTaskEntityManager().DeleteTask(task, nil); err != nil {
		return
	}

	contextutil.GetAgendaFromContext(commandContext).PlanTriggerExecutionOperation(&executionEntity)
	return nil
}

func NewCompleteCmd(taskId string, formData map[string]any, userId *string, options ...taskcmd.Options) CompleteCmd {
	completeCmd := CompleteCmd{
		ProcessVariables: formData,
		UserId:           userId,
	}
	completeCmd.NeedsActiveTaskCmd = taskcmd.NeedsActiveTaskCmd{
		ITaskCmd: &completeCmd,
		TaskId:   taskId,
	}

	for _, option := range options {
		option(&completeCmd.NeedsActiveTaskCmd)
	}

	return completeCmd
}
