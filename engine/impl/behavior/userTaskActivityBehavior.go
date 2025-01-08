package behavior

import (
	. "github.com/go-cinderella/cinderella-engine/engine/constant"
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/go-cinderella/cinderella-engine/engine/eventmanager"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	. "github.com/go-cinderella/cinderella-engine/engine/impl/handler"
	model2 "github.com/go-cinderella/cinderella-engine/engine/internal/model"
	"github.com/go-cinderella/cinderella-engine/engine/utils"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/unionj-cloud/toolkit/stringutils"
	"strings"
	"time"
)

var _ delegate.ActivityBehavior = (*UserTaskActivityBehavior)(nil)
var _ delegate.TriggerableActivityBehavior = (*UserTaskActivityBehavior)(nil)

type UserTaskActivityBehavior struct {
	abstractBpmnActivityBehavior
	UserTask   model.UserTask
	ProcessKey string
}

// 普通用户节点处理
func (user UserTaskActivityBehavior) Execute(execution delegate.DelegateExecution) error {
	task := entitymanager.NewTaskEntity(execution, user.UserTask)
	if user.UserTask.DueDate != nil && stringutils.IsNotEmpty(*user.UserTask.DueDate) {
		if dueDate, err := time.Parse(time.DateTime, *user.UserTask.DueDate); err == nil {
			task.DueDate = &dueDate
		}
	}

	assignee := user.UserTask.Assignee

	if assignee != nil && stringutils.IsNotEmpty(*assignee) {
		assigneeStr := *assignee

		if strings.HasPrefix(assigneeStr, "${") && strings.HasSuffix(assigneeStr, "}") {
			expressionManager := contextutil.GetExpressionManager()
			context := expressionManager.EvaluationContext()

			variable := execution.GetProcessVariable()
			if len(variable) > 0 {
				context.SetVariables(variable)
			}

			expression := expressionManager.CreateExpression(assigneeStr)
			value := expression.GetValueContext(&context)
			b := cast.ToString(value)
			if stringutils.IsNotEmpty(b) {
				task.Assignee = &b
			}
		} else {
			task.Assignee = assignee
		}

		if task.Assignee != nil {
			task.ClaimTime = &task.StartTime
		}
	}

	taskEntityManager := entitymanager.GetTaskEntityManager()
	if err := taskEntityManager.InsertTask(&task); err != nil {
		return err
	}

	if err := handleAssignments(user.UserTask, task, execution); err != nil {
		return err
	}

	// All properties set, now firing 'create' events
	activitiEntityEvent := eventmanager.CreateEntityEvent(eventmanager.TASK_CREATED, task)
	if err := eventmanager.GetEventDispatcher().DispatchEvent(activitiEntityEvent); err != nil {
		return err
	}

	extensionElements := user.UserTask.ExtensionElements
	if extensionElements.TaskListener != nil && len(extensionElements.TaskListener) > 0 {
		taskListeners := extensionElements.TaskListener
		for _, listener := range taskListeners {
			if listener.EventType == TASK_TYPE_CREATE {
				if err := PerformTaskListener(execution, user.UserTask.Name, user.ProcessKey); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func handleAssignments(user model.UserTask, task entitymanager.TaskEntity, execution delegate.DelegateExecution) (err error) {
	identityLinkManager := entitymanager.GetIdentityLinkManager()

	if task.Assignee != nil {
		link := model2.ActRuIdentitylink{
			Rev_:        lo.ToPtr(int32(1)),
			Type_:       lo.ToPtr("assignee"),
			UserID_:     task.Assignee,
			TaskID_:     &task.Id,
			ProcInstID_: &task.ProcessInstanceId,
		}
		err = identityLinkManager.CreateIdentityLink(link)
		if err != nil {
			return err
		}

		link = model2.ActRuIdentitylink{
			Rev_:        lo.ToPtr(int32(1)),
			Type_:       lo.ToPtr("participant"),
			UserID_:     task.Assignee,
			ProcInstID_: &task.ProcessInstanceId,
		}
		err = identityLinkManager.CreateIdentityLink(link)
		if err != nil {
			return err
		}
	}

	variables := execution.GetProcessVariable()

	if user.CandidateUsers != nil {

		candidateUsers := *user.CandidateUsers
		var users []string

		if utils.IsExpr(candidateUsers) {
			users = utils.GetStringSliceFromExpression(variables, candidateUsers)
		} else {
			users = strings.Split(candidateUsers, ",")
			userSlices := lo.Map[string, []string](users, func(item string, index int) []string {
				if !utils.IsExpr(item) {
					return []string{item}
				}
				return utils.GetStringSliceFromExpression(variables, item)
			})
			users = lo.Reduce[[]string, []string](userSlices, func(agg []string, item []string, index int) []string {
				if len(item) > 0 {
					return append(agg, item...)
				}
				return agg
			}, []string{})
		}

		for _, userId := range users {
			link := model2.ActRuIdentitylink{
				Rev_:        lo.ToPtr(int32(1)),
				Type_:       lo.ToPtr("candidate"),
				UserID_:     &userId,
				TaskID_:     &task.Id,
				ProcInstID_: &task.ProcessInstanceId,
			}
			err = identityLinkManager.CreateIdentityLink(link)
			if err != nil {
				return err
			}

			link = model2.ActRuIdentitylink{
				Rev_:        lo.ToPtr(int32(1)),
				Type_:       lo.ToPtr("participant"),
				UserID_:     &userId,
				ProcInstID_: &task.ProcessInstanceId,
			}
			err = identityLinkManager.CreateIdentityLink(link)
			if err != nil {
				return err
			}
		}
	}

	if user.CandidateGroups != nil {

		candidateGroups := *user.CandidateGroups
		var groups []string

		if utils.IsExpr(candidateGroups) {
			groups = utils.GetStringSliceFromExpression(variables, candidateGroups)
		} else {
			groups = strings.Split(candidateGroups, ",")
			groupSlices := lo.Map[string, []string](groups, func(item string, index int) []string {
				if !utils.IsExpr(item) {
					return []string{item}
				}
				return utils.GetStringSliceFromExpression(variables, item)
			})
			groups = lo.Reduce[[]string, []string](groupSlices, func(agg []string, item []string, index int) []string {
				if len(item) > 0 {
					return append(agg, item...)
				}
				return agg
			}, []string{})
		}

		for _, group := range groups {
			link := model2.ActRuIdentitylink{
				Rev_:        lo.ToPtr(int32(1)),
				Type_:       lo.ToPtr("candidate"),
				GroupID_:    &group,
				TaskID_:     &task.Id,
				ProcInstID_: &task.ProcessInstanceId,
			}
			err = identityLinkManager.CreateIdentityLink(link)
			if err != nil {
				return err
			}
		}

	}

	return err
}

// Trigger 普通用户节点处理
func (user UserTaskActivityBehavior) Trigger(execution delegate.DelegateExecution) error {
	user.leave(execution)
	return nil
}
