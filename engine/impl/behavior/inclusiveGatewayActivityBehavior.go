package behavior

import (
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/datamanager"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/go-cinderella/cinderella-engine/engine/utils"
	"github.com/spf13/cast"
)

type InclusiveGatewayActivityBehavior struct {
}

// 包容网关
func (exclusive InclusiveGatewayActivityBehavior) Execute(execution delegate.DelegateExecution) error {
	return exclusive.Leave(execution)
}

// 执行逻辑：获取当前所有执行的节点，判断是否可达当前网关可以停止执行，等待完成
func (exclusive InclusiveGatewayActivityBehavior) Leave(execution delegate.DelegateExecution) error {
	processInstanceId := execution.GetProcessInstanceId()
	taskManager := datamanager.GetTaskDataManager()
	//查询当前执行节点
	tasks, errS := taskManager.FindByProcessInstanceId(processInstanceId)
	var oneExecutionCanReachGateway = false
	var err error
	if errS != nil {
		for _, task := range tasks {
			if cast.ToString(task.TaskDefKey_) != execution.GetCurrentActivityId() {
				//判断是否可以继续执行
				oneExecutionCanReachGateway, err = utils.IsReachable(execution.GetProcessDefinitionId(), cast.ToString(task.TaskDefKey_), execution.GetCurrentActivityId())
				if err != nil {
					return err
				}
			} else {
				oneExecutionCanReachGateway = true
			}
		}
	}
	if !oneExecutionCanReachGateway {
		//执行出口逻辑，设置条件判断
		contextutil.GetAgenda().PlanTakeOutgoingSequenceFlowsOperation(execution, true)
	}
	return nil
}
