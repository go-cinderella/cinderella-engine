package entitymanager

type VariableEntity struct {
	AbstractEntity
	Type   string
	Delete bool
}

func (variableEntity VariableEntity) GetType() string {
	return variableEntity.Type
}

func (variableEntity VariableEntity) SetDeleted(b bool) {
	variableEntity.Delete = b
}
