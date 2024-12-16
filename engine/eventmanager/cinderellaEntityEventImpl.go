package eventmanager

type CinderellaEntityEventImpl struct {
	CinderellaEntityEvent
	CinderellaEventImpl
	Entity interface{}
}

func (CinderellaEntityEventImpl) GetType() CinderellaEventType {
	return TASK_CREATED
}
