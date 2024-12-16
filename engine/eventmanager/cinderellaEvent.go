package eventmanager

type CinderellaEventType string

const (
	TASK_CREATED CinderellaEventType = "TASK_CREATED"

	TASK_ASSIGNED CinderellaEventType = "TASK_ASSIGNED"

	TASK_COMPLETED CinderellaEventType = "TASK_COMPLETED"
)

type CinderellaEvent interface {
	GetType() CinderellaEventType
}
