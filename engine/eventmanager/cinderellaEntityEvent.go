package eventmanager

type CinderellaEntityEvent interface {
	CinderellaEvent
	GetEntity() interface{}
}
