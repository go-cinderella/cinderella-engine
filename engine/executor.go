package engine

type Executor interface {
	Exe(conf Command) (interface{}, error)
}
