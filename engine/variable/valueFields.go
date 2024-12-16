package variable

type ValueFields interface {
	GetName() string

	GetProcessInstanceId() string

	GetTaskId() string

	GetNumberValue() int

	SetNumberValue(value int)

	GetTextValue() string

	SetTextValue(value string)
}
