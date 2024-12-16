package behavior

import (
	. "github.com/go-cinderella/cinderella-engine/engine/constant"
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	. "github.com/go-cinderella/cinderella-engine/engine/impl/handler"
	"reflect"
	"time"
)

type UserAutoTaskActivityBehavior struct {
	UserTask   model.UserTask
	ProcessKey string
}

// 自动通过用户节点处理
func (user UserAutoTaskActivityBehavior) Execute(execution delegate.DelegateExecution) (err error) {
	task := entitymanager.TaskEntity{}
	task.ProcessInstanceId = execution.GetProcessInstanceId()
	task.SetAssignee(user.UserTask.Assignee)
	task.SetStartTime(time.Now().UTC())
	task.SetTaskDefineKey(user.UserTask.Id)
	task.SetTaskDefineName(user.UserTask.Name)
	dataManager := entitymanager.GetTaskEntityManager()
	err = dataManager.InsertTask(&task)
	if err != nil {
		return err
	}
	activitiConstructor, err := GetConstructorByName(user.ProcessKey)
	if err != nil {
		dataManager.DeleteTask(task, nil)
		contextutil.GetAgenda().PlanTriggerExecutionOperation(execution)
		return nil
	}
	constructor := activitiConstructor(execution)
	reflectConstructor := reflect.ValueOf(constructor)
	taskParams := []reflect.Value{reflectConstructor}

	method, b := reflectConstructor.Type().MethodByName(user.UserTask.Name)
	if !b {
		dataManager.DeleteTask(task, nil)
		contextutil.GetAgenda().PlanTriggerExecutionOperation(execution)
		return err
	}

	callResponse := method.Func.Call(taskParams)

	code := callResponse[0].Interface()
	errRes := callResponse[1].Interface()
	code = code.(string)
	if code != ACTIVITI_HANDLER_CODE {
		err := errRes.(error)
		return err
	}
	dataManager.DeleteTask(task, nil)
	contextutil.GetAgenda().PlanTriggerExecutionOperation(execution)
	return err
}

// 普通用户节点处理
func (user UserAutoTaskActivityBehavior) Trigger(execution delegate.DelegateExecution) error {
	user.Leave(execution)
	return nil
}

func (user UserAutoTaskActivityBehavior) Leave(execution delegate.DelegateExecution) {
	element := execution.GetCurrentFlowElement()
	execution.SetCurrentFlowElement(element)
	contextutil.GetAgenda().PlanTakeOutgoingSequenceFlowsOperation(execution, true)
}
