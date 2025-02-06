package behavior

import "github.com/go-cinderella/cinderella-engine/engine/impl/delegate"

var _ delegate.ActivityBehavior = (*abstractBpmnActivityBehavior)(nil)
var _ delegate.TriggerableActivityBehavior = (*abstractBpmnActivityBehavior)(nil)

type abstractBpmnActivityBehavior struct {
	flowNodeActivityBehavior
	multiInstanceActivityBehavior *AbstractMultiInstanceActivityBehavior
}

func (f *abstractBpmnActivityBehavior) SetMultiInstanceActivityBehavior(multiInstanceActivityBehavior *AbstractMultiInstanceActivityBehavior) {
	f.multiInstanceActivityBehavior = multiInstanceActivityBehavior
}

func (f abstractBpmnActivityBehavior) MultiInstanceActivityBehavior() *AbstractMultiInstanceActivityBehavior {
	return f.multiInstanceActivityBehavior
}

func (f abstractBpmnActivityBehavior) leave(execution delegate.DelegateExecution) error {
	var err error
	if !f.hasLoopCharacteristics() {
		err = f.flowNodeActivityBehavior.leave(execution)
	} else if f.hasMultiInstanceCharacteristics() {
		err = f.multiInstanceActivityBehavior.leave(execution)
	}
	return err
}

func (f abstractBpmnActivityBehavior) hasLoopCharacteristics() bool {
	return f.hasMultiInstanceCharacteristics()
}

func (f abstractBpmnActivityBehavior) hasMultiInstanceCharacteristics() bool {
	return f.multiInstanceActivityBehavior != nil
}
