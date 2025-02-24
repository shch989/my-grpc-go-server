package port

import (
	"time"

	"github.com/google/uuid"
	db "github.com/shch989/my-grpc-go-server/internal/adapter/database"
)

type DummyDatabasePort interface {
	Save(data *db.DummyOrm) (uuid.UUID, error)
	GetByUuid(uuid *uuid.UUID) (db.DummyOrm, error)
}

type BankDatabasePort interface {
	GetBankAccountByAccountNumber(acct string) (db.BankAccountOrm, error)
	CreateExchangeRate(r db.BankExchangeRateOrm) (uuid.UUID, error)
	GetExchangeRateAtTimestamp(fromCur string, toCur string, ts time.Time) (db.BankExchangeRateOrm, error)
	CreateTransaction(acct db.BankAccountOrm, t db.BankTransactionOrm) (uuid.UUID, error)
}
