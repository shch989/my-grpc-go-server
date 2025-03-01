package port

import (
	"time"

	"github.com/google/uuid"
	dbank "github.com/shch989/my-grpc-go-server/internal/application/domain/bank"
)

type HelloServicePort interface {
	GenerateHello(name string) string
}

type BankServicePort interface {
	FindCurrentBalance(acct string) float64
	CreateExchangeRate(r dbank.ExchangeRate) (uuid.UUID, error)
	FindExchangeRate(fromCur string, toCur string, ts time.Time) float64
	CreateTransaction(acct string, t dbank.Transaction) (uuid.UUID, error)
	CalculateTransactionSummary(tcur *dbank.TransactionSummary, trans dbank.Transaction) error
	Transfer(tt dbank.TransferTransaction) (uuid.UUID, bool, error)
}
