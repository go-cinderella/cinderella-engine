package entitymanager

type AbstractEntity struct {
	Id string `json:"id"`
}

func (entity *AbstractEntity) GetId() string {
	return entity.Id
}

func (entity *AbstractEntity) SetId(id string) {
	entity.Id = id
}
