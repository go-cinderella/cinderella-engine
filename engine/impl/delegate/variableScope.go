package delegate

type VariableScope interface {
	GetProcessVariables() (map[string]interface{}, error)

	GetVariablesLocal() (map[string]interface{}, error)

	SetVariableLocal(variableName string, value interface{}) error
}
