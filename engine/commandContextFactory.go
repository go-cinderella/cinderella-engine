package engine

type ICommandContextFactory interface {
	CreateCommandContext() Context
}
