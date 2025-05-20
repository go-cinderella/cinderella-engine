package handler

import (
	. "github.com/go-cinderella/cinderella-engine/engine/constant"
	"github.com/go-cinderella/cinderella-engine/engine/errs"
	delegate2 "github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"reflect"
	"sync"
)

var gConstructorMap map[string]ActivitiConstructor
var lock sync.Mutex

func init() {
	gConstructorMap = make(map[string]ActivitiConstructor, 0)
}

type IActiviti interface {
	GetInPut() interface{}
	GetOutPut() interface{}
}

type ActivitiConstructor func(entity delegate2.DelegateExecution) IActiviti

func RegisterConstructor(name string, constructor ActivitiConstructor) error {
	lock.Lock()
	defer lock.Unlock()
	_, ok := gConstructorMap[name]
	if !ok {
		gConstructorMap[name] = constructor
	} else {
		return errs.CinderellaError{Code: "1005", Msg: "name has register"}
	}
	return nil
}

func GetConstructorByName(name string) (ActivitiConstructor, error) {
	lock.Lock()
	defer lock.Unlock()
	constructor, ok := gConstructorMap[name]
	if !ok {
		log.Error("not find Constructor handle name:", name)
		return nil, errs.CinderellaError{Code: "1006", Msg: "name not find"}
	}
	return constructor, nil
}

func PerformTaskListener(entity delegate2.DelegateExecution, taskName, processKey string) error {
	activitiConstructor, err := GetConstructorByName(processKey)
	if err != nil {
		return err
	}
	constructor := activitiConstructor(entity)
	reflectConstructor := reflect.ValueOf(constructor)
	taskParams := []reflect.Value{reflectConstructor}

	method, b := reflectConstructor.Type().MethodByName(taskName)
	if !b {
		return nil
	}

	callResponse := method.Func.Call(taskParams)
	code := callResponse[0].Interface()
	errRes := callResponse[1].Interface()
	code = cast.ToString(code)
	if code != ACTIVITI_HANDLER_CODE {
		err := errRes.(error)
		return err
	}
	return nil
}
