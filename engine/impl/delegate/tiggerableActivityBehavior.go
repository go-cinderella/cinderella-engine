package delegate

type TriggerableActivityBehavior interface {
	Trigger(execution DelegateExecution) error
}
