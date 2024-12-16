package uuidgenerator

import (
	"github.com/go-cinderella/cinderella-engine/engine/idgenerator"
	"github.com/google/uuid"
)

var _ idgenerator.IDGenerator = (*UUIDGenerator)(nil)

type UUIDGenerator struct {
}

func (u UUIDGenerator) NextID() (string, error) {
	id, err := uuid.NewUUID()
	return id.String(), err
}
