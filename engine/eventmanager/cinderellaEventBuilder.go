package eventmanager

type CinderellaEventBuilder struct {
}

func CreateEntityEvent(eventType CinderellaEventType, entity interface{}) CinderellaEntityEvent {
	entityEventImpl := CinderellaEntityEventImpl{}
	entityEventImpl.CinderellaEventImpl = CinderellaEventImpl{}
	entityEventImpl.EventType = eventType
	entityEventImpl.Entity = entity
	return entityEventImpl
}
