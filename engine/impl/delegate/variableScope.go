package delegate

type VariableScope interface {
	GetVariables() (map[string]interface{}, error)

	GetVariablesLocal() (map[string]interface{}, error)

	GetVariableLocal(variableName string) (value interface{}, ok bool, err error)

	SetVariableLocal(variableName string, value interface{}) error

	SetProcessVariables(variables map[string]interface{}) error

	RemoveVariablesLocal(variableNames []string) error

	SetVariablesLocal(variables map[string]interface{}) error
}
