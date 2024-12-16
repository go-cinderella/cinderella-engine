package eventmanager

type CinderellaEventListener interface {
	OnEvent(event CinderellaEvent) error
}
