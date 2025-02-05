package delegate

type VariableScope interface {
	GetVariables() map[string]interface{}

	GetProcessVariables() map[string]interface{}

	GetVariablesLocal() map[string]interface{}
}
