package delegate

// 通用字段
type BaseElement interface {
	GetId() string
	GetName() string
	GetHandlerType() string
}

type BaseHandlerType interface {
	GetType() string
}
