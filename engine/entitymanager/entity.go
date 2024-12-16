package entitymanager

type Entity interface {
	GetId() string
	SetId(id string)
}
