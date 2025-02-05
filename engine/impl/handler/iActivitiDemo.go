package handler

import (
	"github.com/go-cinderella/cinderella-engine/engine/constant"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

func init() {
	RegisterConstructor("userAuto", NewTestIActiviti)
}

func NewTestIActiviti(entity delegate.DelegateExecution) IActiviti {
	return &TestIActiviti{
		Entity: entity,
	}
}

type TestIActiviti struct {
	Entity delegate.DelegateExecution
	InPut  string
	OutPut string
}

func (test *TestIActiviti) GetInPut() interface{} {
	return test.InPut
}

func (test *TestIActiviti) GetOutPut() interface{} {
	return test.OutPut
}

func (test *TestIActiviti) User001() (code interface{}, err error) {
	//variable := test.Entity.GetVariables()
	//fmt.Println(variable)
	return constant.ACTIVITI_HANDLER_CODE, nil
}

func (test *TestIActiviti) User002() (code interface{}, err error) {
	return constant.ACTIVITI_HANDLER_CODE, nil
}
