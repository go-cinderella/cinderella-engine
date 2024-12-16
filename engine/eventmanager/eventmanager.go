package eventmanager

var eventDispatcher = &CinderellaEventDispatcher{}

func GetEventDispatcher() *CinderellaEventDispatcher {
	return eventDispatcher
}
