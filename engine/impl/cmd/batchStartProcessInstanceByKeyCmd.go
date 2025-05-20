package cmd

import (
	"context"
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/datamanager"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/go-cinderella/cinderella-engine/engine/impl/behavior"
	bpmn_model "github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/go-cinderella/cinderella-engine/engine/model"
	"github.com/go-cinderella/cinderella-engine/engine/utils"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/unionj-cloud/toolkit/stringutils"
	"time"
)

var _ engine.Command = (*BatchStartProcessInstanceByKeyCmd)(nil)

type BatchStartProcessInstanceByKeyCmd struct {
	ProcessDefinitionKey string
	ProcessInstanceId    string
	StartActivityId      string
	StartActivityName    string
	TenantId             string
	UserId               string
	Ctx                  context.Context
	Transactional        bool
}

func (receiver BatchStartProcessInstanceByKeyCmd) IsTransactional() bool {
	return receiver.Transactional
}

func (receiver BatchStartProcessInstanceByKeyCmd) Context() context.Context {
	return receiver.Ctx
}

func (receiver BatchStartProcessInstanceByKeyCmd) Execute(commandContext engine.Context) (interface{}, error) {
	variables := make(map[string]interface{})
	variables["initiator"] = receiver.UserId

	processDefinitionEntityManager := entitymanager.GetProcessDefinitionEntityManager()
	definitionEntity, err := processDefinitionEntityManager.FindLatestProcessDefinitionByKey(receiver.ProcessDefinitionKey)
	if err != nil {
		return entitymanager.ExecutionEntity{}, errors.WithStack(err)
	}

	processUtils := utils.ProcessDefinitionUtil{}
	process, err := processUtils.GetProcess(definitionEntity.GetId())
	if err != nil {
		return entitymanager.ExecutionEntity{}, errors.WithStack(err)
	}

	startActivityId := receiver.StartActivityId
	if stringutils.IsEmpty(startActivityId) {
		flowElementList := lo.Filter(process.FlowElementList, func(item delegate.FlowElement, index int) bool {
			switch item.(type) {
			case *bpmn_model.UserTask, *bpmn_model.ServiceTask:
				return true
			default:
				return false
			}
		})

		if len(flowElementList) == 0 {
			return entitymanager.ExecutionEntity{}, nil
		}

		if stringutils.IsNotEmpty(receiver.StartActivityName) {
			filtered := lo.Filter(flowElementList, func(item delegate.FlowElement, index int) bool {
				return item.GetName() == receiver.StartActivityName
			})

			if len(filtered) > 0 {
				startActivityId = filtered[0].GetId()
			}
		}

		if stringutils.IsEmpty(startActivityId) {
			startActivityId = flowElementList[0].GetId()
		}
	}

	// 默认从开始节点开始流程流转
	flowElement := process.InitialFlowElement
	element := flowElement.(*bpmn_model.StartEvent)

	processInstanceIds := stringutils.Split(receiver.ProcessInstanceId, ",")

	var processInstanceSaves []*model.ActRuExecution
	var historicProcessSaves []*model.ActHiProcinst
	var linkSaves []*model.ActRuIdentitylink
	var hisLinkSaves []*model.ActHiIdentitylink
	var hisActinstSaves []*model.ActHiActinst
	var taskSaves []*model.ActRuTask
	var histTaskSaves []*model.ActHiTaskinst

	now := time.Now().UTC()

	lo.ForEach(processInstanceIds, func(item string, index int) {
		processInstance := model.ActRuExecution{}
		processInstance.Rev_ = lo.ToPtr(int32(1))
		processInstance.ID_ = item
		processInstance.TenantID_ = &receiver.TenantId
		processInstance.StartTime_ = &now
		processInstance.ProcDefID_ = lo.ToPtr(definitionEntity.GetId())
		processInstance.IsActive_ = lo.ToPtr(true)
		processInstance.StartUserID_ = &receiver.UserId
		processInstance.StartActID_ = lo.ToPtr(element.GetId())
		processInstance.ProcInstID_ = &processInstance.ID_
		processInstance.RootProcInstID_ = &processInstance.ID_
		processInstance.IsConcurrent_ = lo.ToPtr(false)
		processInstance.IsScope_ = lo.ToPtr(true)
		processInstance.IsEventScope_ = lo.ToPtr(false)
		processInstance.IsMiRoot_ = lo.ToPtr(false)
		processInstance.SuspensionState_ = lo.ToPtr(int32(1))
		processInstance.IsCountEnabled_ = lo.ToPtr(true)
		processInstance.BusinessStatus_ = &startActivityId

		processInstanceSaves = append(processInstanceSaves, &processInstance)

		historicProcess := model.ActHiProcinst{}
		historicProcess.ID_ = processInstance.ID_
		historicProcess.ProcInstID_ = processInstance.ID_
		historicProcess.ProcDefID_ = cast.ToString(processInstance.ProcDefID_)
		historicProcess.TenantID_ = processInstance.TenantID_
		historicProcess.StartTime_ = cast.ToTime(processInstance.StartTime_)
		historicProcess.Name_ = processInstance.Name_
		historicProcess.BusinessKey_ = processInstance.BusinessKey_
		historicProcess.StartUserID_ = processInstance.StartUserID_
		historicProcess.StartActID_ = processInstance.StartActID_
		historicProcess.SuperProcessInstanceID_ = processInstance.SuperExec_
		historicProcess.BusinessStatus_ = &startActivityId

		historicProcessSaves = append(historicProcessSaves, &historicProcess)

		identityLink := model.ActRuIdentitylink{
			Rev_:        lo.ToPtr(int32(1)),
			Type_:       lo.ToPtr("starter"),
			UserID_:     processInstance.StartUserID_,
			ProcInstID_: processInstance.ProcInstID_,
		}
		linkSaves = append(linkSaves, &identityLink)

		var historicIdentityLink model.ActHiIdentitylink
		copier.Copy(&historicIdentityLink, &identityLink)
		historicIdentityLink.CreateTime_ = &now
		hisLinkSaves = append(hisLinkSaves, &historicIdentityLink)

		processInstanceEntity := entitymanager.ExecutionEntity{}
		processInstanceEntity.SetId(processInstance.ID_)
		processInstanceEntity.SetProcessDefinitionId(*processInstance.ProcDefID_)
		processInstanceEntity.SetProcessInstanceId(*processInstance.ProcInstID_)

		execution := entitymanager.CreateChildExecution(&processInstanceEntity)
		execution.SetCurrentFlowElement(element)

		actRuExecutionId, _ := contextutil.GetIDGenerator().NextID()
		execution.SetId(actRuExecutionId)

		// 最后还要删掉，所以这里就不存了
		//actRuExecution := model.ActRuExecution{
		//	ID_:                     actRuExecutionId,
		//	Rev_:                    lo.ToPtr(cast.ToInt32(1)),
		//	ProcInstID_:             &processInstanceEntity.ProcessInstanceId,
		//	ParentID_:               &processInstanceEntity.ParentId,
		//	ProcDefID_:              &processInstanceEntity.ProcessDefinitionId,
		//	RootProcInstID_:         &processInstanceEntity.ProcessInstanceId,
		//	ActID_:                  &processInstanceEntity.CurrentActivityId,
		//	IsActive_:               lo.ToPtr(true),
		//	IsConcurrent_:           lo.ToPtr(false),
		//	IsScope_:                lo.ToPtr(false),
		//	IsEventScope_:           lo.ToPtr(false),
		//	IsMiRoot_:               lo.ToPtr(false),
		//	SuspensionState_:        lo.ToPtr(cast.ToInt32(1)),
		//	StartTime_:              &processInstanceEntity.StartTime,
		//	IsCountEnabled_:         lo.ToPtr(true),
		//	EvtSubscrCount_:         lo.ToPtr(cast.ToInt32(0)),
		//	TaskCount_:              lo.ToPtr(cast.ToInt32(0)),
		//	JobCount_:               lo.ToPtr(cast.ToInt32(0)),
		//	TimerJobCount_:          lo.ToPtr(cast.ToInt32(0)),
		//	SuspJobCount_:           lo.ToPtr(cast.ToInt32(0)),
		//	DeadletterJobCount_:     lo.ToPtr(cast.ToInt32(0)),
		//	ExternalWorkerJobCount_: lo.ToPtr(cast.ToInt32(0)),
		//	VarCount_:               lo.ToPtr(cast.ToInt32(0)),
		//	IDLinkCount_:            lo.ToPtr(cast.ToInt32(0)),
		//}
		//processInstanceSaves = append(processInstanceSaves, &actRuExecution)

		historicActinst := model.ActHiActinst{}
		historicActinst.ExecutionID_ = execution.GetExecutionId()
		historicActinst.ProcDefID_ = execution.GetProcessDefinitionId()
		historicActinst.ProcInstID_ = execution.GetProcessInstanceId()
		historicActinst.ActID_ = execution.GetCurrentActivityId()
		historicActinst.StartTime_ = now
		if execution.GetCurrentFlowElement() != nil {
			historicActinst.ActName_ = lo.ToPtr(execution.GetCurrentFlowElement().GetName())
			historicActinst.ActType_ = execution.GetCurrentFlowElement().GetHandlerType()
		}
		historicActinst.EndTime_ = &now
		historicActinst.Duration_ = lo.ToPtr(int64(0))
		hisActinstSaves = append(hisActinstSaves, &historicActinst)

		flowElements := element.GetOutgoing()
		outgoingSequenceFlow := flowElements[0]

		executionEntity := entitymanager.CreateExecution(&execution)
		executionEntity.SetCurrentFlowElement(outgoingSequenceFlow)

		actRuExecutionId, _ = contextutil.GetIDGenerator().NextID()
		executionEntity.SetId(actRuExecutionId)

		// 最后还要删掉，所以这里就不存了
		//actRuExecution := model.ActRuExecution{
		//	ID_:                     actRuExecutionId,
		//	Rev_:                    lo.ToPtr(cast.ToInt32(1)),
		//	ProcInstID_:             &executionEntity.ProcessInstanceId,
		//	ParentID_:               &executionEntity.ParentId,
		//	ProcDefID_:              &executionEntity.ProcessDefinitionId,
		//	RootProcInstID_:         &executionEntity.ProcessInstanceId,
		//	ActID_:                  &executionEntity.CurrentActivityId,
		//	IsActive_:               lo.ToPtr(true),
		//	IsConcurrent_:           lo.ToPtr(false),
		//	IsScope_:                lo.ToPtr(false),
		//	IsEventScope_:           lo.ToPtr(false),
		//	IsMiRoot_:               lo.ToPtr(false),
		//	SuspensionState_:        lo.ToPtr(cast.ToInt32(1)),
		//	StartTime_:              &executionEntity.StartTime,
		//	IsCountEnabled_:         lo.ToPtr(true),
		//	EvtSubscrCount_:         lo.ToPtr(cast.ToInt32(0)),
		//	TaskCount_:              lo.ToPtr(cast.ToInt32(0)),
		//	JobCount_:               lo.ToPtr(cast.ToInt32(0)),
		//	TimerJobCount_:          lo.ToPtr(cast.ToInt32(0)),
		//	SuspJobCount_:           lo.ToPtr(cast.ToInt32(0)),
		//	DeadletterJobCount_:     lo.ToPtr(cast.ToInt32(0)),
		//	ExternalWorkerJobCount_: lo.ToPtr(cast.ToInt32(0)),
		//	VarCount_:               lo.ToPtr(cast.ToInt32(0)),
		//	IDLinkCount_:            lo.ToPtr(cast.ToInt32(0)),
		//}
		//processInstanceSaves = append(processInstanceSaves, &actRuExecution)

		seqHistActinst := model.ActHiActinst{}
		seqHistActinst.ProcDefID_ = executionEntity.GetProcessDefinitionId()
		seqHistActinst.ProcInstID_ = executionEntity.GetProcessInstanceId()
		seqHistActinst.ActID_ = executionEntity.GetCurrentActivityId()
		seqHistActinst.StartTime_ = now
		seqHistActinst.EndTime_ = &now
		seqHistActinst.Duration_ = lo.ToPtr(int64(0))
		if executionEntity.GetCurrentFlowElement() != nil {
			seqHistActinst.ActName_ = lo.ToPtr(executionEntity.GetCurrentFlowElement().GetName())
			seqHistActinst.ActType_ = executionEntity.GetCurrentFlowElement().GetHandlerType()
		}
		hisActinstSaves = append(hisActinstSaves, &seqHistActinst)

		sequenceFlow := executionEntity.GetCurrentFlowElement()
		executionEntity = entitymanager.CreateExecution(&executionEntity)
		executionEntity.SetCurrentFlowElement(sequenceFlow.GetTargetFlowElement())
		actRuExecutionId, _ = contextutil.GetIDGenerator().NextID()
		executionEntity.SetId(actRuExecutionId)

		actRuExecution := model.ActRuExecution{
			ID_:                     actRuExecutionId,
			Rev_:                    lo.ToPtr(cast.ToInt32(1)),
			ProcInstID_:             &executionEntity.ProcessInstanceId,
			ParentID_:               &executionEntity.ParentId,
			ProcDefID_:              &executionEntity.ProcessDefinitionId,
			RootProcInstID_:         &executionEntity.ProcessInstanceId,
			ActID_:                  &executionEntity.CurrentActivityId,
			IsActive_:               lo.ToPtr(true),
			IsConcurrent_:           lo.ToPtr(false),
			IsScope_:                lo.ToPtr(false),
			IsEventScope_:           lo.ToPtr(false),
			IsMiRoot_:               lo.ToPtr(false),
			SuspensionState_:        lo.ToPtr(cast.ToInt32(1)),
			StartTime_:              &executionEntity.StartTime,
			IsCountEnabled_:         lo.ToPtr(true),
			EvtSubscrCount_:         lo.ToPtr(cast.ToInt32(0)),
			TaskCount_:              lo.ToPtr(cast.ToInt32(0)),
			JobCount_:               lo.ToPtr(cast.ToInt32(0)),
			TimerJobCount_:          lo.ToPtr(cast.ToInt32(0)),
			SuspJobCount_:           lo.ToPtr(cast.ToInt32(0)),
			DeadletterJobCount_:     lo.ToPtr(cast.ToInt32(0)),
			ExternalWorkerJobCount_: lo.ToPtr(cast.ToInt32(0)),
			VarCount_:               lo.ToPtr(cast.ToInt32(0)),
			IDLinkCount_:            lo.ToPtr(cast.ToInt32(0)),
		}
		processInstanceSaves = append(processInstanceSaves, &actRuExecution)

		taskHistActinst := model.ActHiActinst{}
		taskHistActinst.ExecutionID_ = executionEntity.GetExecutionId()
		taskHistActinst.ProcDefID_ = executionEntity.GetProcessDefinitionId()
		taskHistActinst.ProcInstID_ = executionEntity.GetProcessInstanceId()
		taskHistActinst.ActID_ = executionEntity.GetCurrentActivityId()
		taskHistActinst.StartTime_ = now
		if executionEntity.GetCurrentFlowElement() != nil {
			taskHistActinst.ActName_ = lo.ToPtr(executionEntity.GetCurrentFlowElement().GetName())
			taskHistActinst.ActType_ = executionEntity.GetCurrentFlowElement().GetHandlerType()
		}
		hisActinstSaves = append(hisActinstSaves, &taskHistActinst)

		currentFlowElement := executionEntity.GetCurrentFlowElement()
		userTaskBehavior := currentFlowElement.GetBehavior().(*behavior.UserTaskActivityBehavior)

		task := entitymanager.NewTaskEntity(&executionEntity, userTaskBehavior.UserTask)
		if userTaskBehavior.UserTask.DueDate != nil && stringutils.IsNotEmpty(*userTaskBehavior.UserTask.DueDate) {
			if dueDate, err := time.Parse(time.DateTime, *userTaskBehavior.UserTask.DueDate); err == nil {
				task.DueDate = &dueDate
			}
		}

		ruTaskId, _ := contextutil.GetIDGenerator().NextID()
		task.SetId(ruTaskId)

		ruTaskModel := model.ActRuTask{
			ID_:              ruTaskId,
			Rev_:             lo.ToPtr(int32(1)),
			ExecutionID_:     &task.ExecutionId,
			ProcInstID_:      &task.ProcessInstanceId,
			ProcDefID_:       &task.ProcessDefinitionId,
			Name_:            &task.TaskDefineName,
			TaskDefKey_:      &task.TaskDefinitionKey,
			Assignee_:        task.Assignee,
			CreateTime_:      &task.StartTime,
			DueDate_:         task.DueDate,
			SuspensionState_: lo.ToPtr(int32(1)),
			TenantID_:        task.TenantId,
			FormKey_:         task.FormKey,
			ClaimTime_:       task.ClaimTime,
			IsCountEnabled_:  lo.ToPtr(true),
			VarCount_:        lo.ToPtr(int32(0)),
			IDLinkCount_:     lo.ToPtr(int32(0)),
			SubTaskCount_:    lo.ToPtr(int32(0)),
		}
		taskSaves = append(taskSaves, &ruTaskModel)

		taskHistActinst.Assignee_ = ruTaskModel.Assignee_
		taskHistActinst.TaskID_ = &ruTaskModel.ID_

		var historicTask model.ActHiTaskinst
		copier.Copy(&historicTask, &ruTaskModel)
		historicTask.StartTime_ = *ruTaskModel.CreateTime_
		histTaskSaves = append(histTaskSaves, &historicTask)

		user := userTaskBehavior.UserTask

		if user.CandidateUsers != nil {

			candidateUsers := *user.CandidateUsers
			var users []string

			if utils.IsExpr(candidateUsers) {
				users = utils.GetStringSliceFromExpression(variables, candidateUsers)
			} else {
				users = stringutils.Split(candidateUsers, ",")
				userSlices := lo.Map[string, []string](users, func(item string, index int) []string {
					if !utils.IsExpr(item) {
						return []string{item}
					}
					return utils.GetStringSliceFromExpression(variables, item)
				})
				users = lo.Reduce[[]string, []string](userSlices, func(agg []string, item []string, index int) []string {
					return append(agg, item...)
				}, []string{})
			}

			for _, userId := range users {
				if stringutils.IsEmpty(userId) {
					continue
				}

				candidateLink := model.ActRuIdentitylink{
					Rev_:        lo.ToPtr(int32(1)),
					Type_:       lo.ToPtr("candidate"),
					UserID_:     &userId,
					TaskID_:     &task.Id,
					ProcInstID_: &task.ProcessInstanceId,
				}
				linkSaves = append(linkSaves, &candidateLink)

				candidateHistLink := model.ActHiIdentitylink{}
				copier.Copy(&candidateHistLink, &candidateLink)
				candidateHistLink.CreateTime_ = &now
				hisLinkSaves = append(hisLinkSaves, &candidateHistLink)

				participantLink := model.ActRuIdentitylink{
					Rev_:        lo.ToPtr(int32(1)),
					Type_:       lo.ToPtr("participant"),
					UserID_:     &userId,
					ProcInstID_: &task.ProcessInstanceId,
				}
				linkSaves = append(linkSaves, &participantLink)

				participantHistLink := model.ActHiIdentitylink{}
				copier.Copy(&participantHistLink, &participantLink)
				participantHistLink.CreateTime_ = &now
				hisLinkSaves = append(hisLinkSaves, &participantHistLink)
			}
		}

		if user.CandidateGroups != nil {

			candidateGroups := *user.CandidateGroups
			var groups []string

			if utils.IsExpr(candidateGroups) {
				groups = utils.GetStringSliceFromExpression(variables, candidateGroups)
			} else {
				groups = stringutils.Split(candidateGroups, ",")
				groupSlices := lo.Map[string, []string](groups, func(item string, index int) []string {
					if !utils.IsExpr(item) {
						return []string{item}
					}
					return utils.GetStringSliceFromExpression(variables, item)
				})
				groups = lo.Reduce[[]string, []string](groupSlices, func(agg []string, item []string, index int) []string {
					return append(agg, item...)
				}, []string{})
			}

			for _, group := range groups {
				if stringutils.IsEmpty(group) {
					continue
				}

				link := model.ActRuIdentitylink{
					Rev_:        lo.ToPtr(int32(1)),
					Type_:       lo.ToPtr("candidate"),
					GroupID_:    &group,
					TaskID_:     &task.Id,
					ProcInstID_: &task.ProcessInstanceId,
				}
				linkSaves = append(linkSaves, &link)

				candidateHistLink := model.ActHiIdentitylink{}
				copier.Copy(&candidateHistLink, &link)
				candidateHistLink.CreateTime_ = &now
				hisLinkSaves = append(hisLinkSaves, &candidateHistLink)
			}

		}
	})

	executionDataManager := datamanager.GetExecutionDataManager()
	historicProcessManager := datamanager.GetHistoricProcessDataManager()
	linkDataManager := datamanager.GetIdentityLinkDataManager()
	hisLinkDataManager := datamanager.GetHistoricIdentityLinkDataManager()
	actinstDataManager := datamanager.GetHistoricActinstDataManager()
	taskDataManager := datamanager.GetTaskDataManager()
	historicTaskManager := datamanager.GetHistoricTaskDataManager()

	// 生成流程实例
	if len(processInstanceSaves) > 0 {
		if err = executionDataManager.Insert(processInstanceSaves); err != nil {
			return entitymanager.ExecutionEntity{}, errors.WithStack(err)
		}
	}

	if len(historicProcessSaves) > 0 {
		if err = historicProcessManager.Insert(historicProcessSaves); err != nil {
			return entitymanager.ExecutionEntity{}, errors.WithStack(err)
		}
	}

	if len(taskSaves) > 0 {
		if err = taskDataManager.Insert(taskSaves); err != nil {
			return entitymanager.ExecutionEntity{}, errors.WithStack(err)
		}
	}

	if len(histTaskSaves) > 0 {
		if err = historicTaskManager.Insert(histTaskSaves); err != nil {
			return entitymanager.ExecutionEntity{}, errors.WithStack(err)
		}
	}

	if len(linkSaves) > 0 {
		if err = linkDataManager.Insert(linkSaves); err != nil {
			return entitymanager.ExecutionEntity{}, errors.WithStack(err)
		}

		if err = hisLinkDataManager.Insert(hisLinkSaves); err != nil {
			return entitymanager.ExecutionEntity{}, errors.WithStack(err)
		}
	}

	if len(hisActinstSaves) > 0 {
		if err = actinstDataManager.Insert(hisActinstSaves); err != nil {
			return entitymanager.ExecutionEntity{}, errors.WithStack(err)
		}
	}

	return nil, nil
}
