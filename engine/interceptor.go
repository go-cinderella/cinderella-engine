package engine

type Interceptor interface {
	Execute(command Command) (interface{}, error)

	SetNext(next Interceptor)
}
