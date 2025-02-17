package port

import (
	"github.com/google/uuid"
	db "github.com/shch989/my-grpc-go-server/internal/adapter/database"
)

type DummyDatabasePort interface {
	Save(data *db.DummyOrm) (uuid.UUID, error)
	GetByUuid(uuid *uuid.UUID) (db.DummyOrm, error)
}
