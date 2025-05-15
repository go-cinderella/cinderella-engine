package cmd

import (
	"context"
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/go-cinderella/cinderella-engine/engine/impl/converter"
	"github.com/samber/lo"
)

var _ engine.Command = (*ProcessMigrateCmd)(nil)

type DiffActivity struct {
	ActivityId    string
	ActivityName  string
	NewActivityId string
}

type ProcessMigrateCmd struct {
	Transactional bool
	Ctx           context.Context

	OldDeploymentId string
	NewDeploymentId string

	DiffActivities []DiffActivity
}

/**
 * 执行流程迁移命令，将流程从旧版本部署迁移到新版本部署
 *
 * 此方法执行以下迁移操作：
 * 1. 查找旧部署ID和新部署ID对应的流程定义实体
 * 2. 迁移流程实例的流程定义ID和起始活动ID
 * 3. 迁移执行实例的流程定义ID
 * 4. 迁移任务的流程定义ID
 * 5. 迁移历史活动实例的流程定义ID
 * 6. 针对每个差异活动，执行详细的迁移:
 *    - 迁移流程实例业务状态
 *    - 迁移执行实例活动ID
 *    - 迁移任务名称和任务定义键
 *    - 迁移历史活动实例
 *
 * @param commandContext 命令上下文
 * @return 返回nil和可能的错误
 */
func (g ProcessMigrateCmd) Execute(commandContext engine.Context) (interface{}, error) {
	// 获取流程定义实体管理器
	processDefinitionEntityManager := entitymanager.GetProcessDefinitionEntityManager()
	
	// 查找旧部署ID对应的流程定义
	oldProcessDefinitionEntity, err := processDefinitionEntityManager.FindByDeploymentId(g.OldDeploymentId)
	if err != nil {
		return nil, err
	}

	// 查找新部署ID对应的流程定义
	newProcessDefinitionEntity, err := processDefinitionEntityManager.FindByDeploymentId(g.NewDeploymentId)
	if err != nil {
		return nil, err
	}

	// 获取执行实体管理器
	executionEntityManager := entitymanager.GetExecutionEntityManager()
	
	// 迁移流程实例的流程定义ID和起始活动ID
	err = executionEntityManager.MigrateProcessInstanceProcDefIdAndStartActId(oldProcessDefinitionEntity, newProcessDefinitionEntity)
	if err != nil {
		return nil, err
	}

	// 迁移执行实例的流程定义ID
	err = executionEntityManager.MigrateExecutionProcDefID(oldProcessDefinitionEntity, newProcessDefinitionEntity)
	if err != nil {
		return nil, err
	}

	// 获取任务实体管理器
	taskEntityManager := entitymanager.GetTaskEntityManager()
	
	// 迁移任务的流程定义ID
	err = taskEntityManager.MigrateProcDefID(oldProcessDefinitionEntity, newProcessDefinitionEntity)
	if err != nil {
		return nil, err
	}

	// 获取历史活动实例实体管理器
	historicActivityInstanceEntityManager := entitymanager.GetHistoricActivityInstanceEntityManager()
	
	// 迁移历史活动实例的流程定义ID
	err = historicActivityInstanceEntityManager.MigrateProcDefID(oldProcessDefinitionEntity, newProcessDefinitionEntity)
	if err != nil {
		return nil, err
	}

	// 创建BPMN XML转换器，并获取新流程定义的流程模型
	bpmnXMLConverter := converter.BpmnXMLConverter{}
	bpmnModel := bpmnXMLConverter.ConvertToBpmnModel(newProcessDefinitionEntity.ResourceContent)
	process := bpmnModel.GetProcess()

	// 遍历差异活动并执行详细迁移
	lo.ForEachWhile(g.DiffActivities, func(item DiffActivity, index int) (goon bool) {
		// 迁移流程实例业务状态
		err = executionEntityManager.MigrateProcessInstanceBusinessStatus(newProcessDefinitionEntity, item.ActivityId, item.NewActivityId)
		if err != nil {
			return false
		}

		// 迁移执行实例活动ID
		err = executionEntityManager.MigrateExecutionActID(newProcessDefinitionEntity, item.ActivityId, item.NewActivityId)
		if err != nil {
			return false
		}

		// 获取新流程元素，并迁移任务名称和任务定义键
		newFlowElement := process.GetFlowElement(item.NewActivityId)
		err = taskEntityManager.MigrateNameAndTaskDefKey(newProcessDefinitionEntity, item.ActivityId, item.NewActivityId, newFlowElement.GetName())
		if err != nil {
			return false
		}

		// 迁移历史活动实例
		err = historicActivityInstanceEntityManager.MigrateAct(newProcessDefinitionEntity, item.ActivityId, item.NewActivityId, newFlowElement.GetName(), newFlowElement.GetHandlerType())
		if err != nil {
			return false
		}

		return true
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (g ProcessMigrateCmd) Context() context.Context {
	return g.Ctx
}

func (g ProcessMigrateCmd) IsTransactional() bool {
	return g.Transactional
}
