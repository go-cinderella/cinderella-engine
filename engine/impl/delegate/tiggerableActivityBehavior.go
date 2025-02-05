package delegate

type TriggerableActivityBehavior interface {
	ActivityBehavior
	Trigger(execution DelegateExecution) error
}
