package errs

import (
	"errors"
	"github.com/unionj-cloud/toolkit/stringutils"
)

type CinderellaError struct {
	Code string
	Msg  string
}

func (error CinderellaError) Error() string {
	if stringutils.IsNotEmpty(error.Code) {
		return error.Code + "---" + error.Msg
	}
	return error.Msg
}

var (
	ErrProcessInstanceNotFound         = errors.New("process instance not found")
	ErrHistoricProcessInstanceNotFound = errors.New("history process instance not found")
	ErrProcessDefinitionNotFound       = errors.New("process definition not found")
	ErrTaskNotFound                    = errors.New("task not found")
	ErrDeploymentNotFound              = errors.New("deployment not found")
	ErrExecutionNotFound               = errors.New("execution not found")
	ErrInternalError                   = errors.New("internal error")
)
