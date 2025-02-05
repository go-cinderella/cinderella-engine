package parse

import "github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"

type activityGetter interface {
	GetActivity() model.Activity
}
