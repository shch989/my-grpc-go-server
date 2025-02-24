package application

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	db "github.com/shch989/my-grpc-go-server/internal/adapter/database"
	dbank "github.com/shch989/my-grpc-go-server/internal/application/domain/bank"
	"github.com/shch989/my-grpc-go-server/internal/port"
)

type BankService struct {
	db port.BankDatabasePort
}

func NewBankService(dbPort port.BankDatabasePort) *BankService {
	return &BankService{
		db: dbPort,
	}
}

func (s *BankService) FindCurrentBalance(acct string) float64 {
	bankAccount, err := s.db.GetBankAccountByAccountNumber(acct)

	if err != nil {
		log.Println("Error on FindCurrentBalance :", err)
	}

	return bankAccount.CurrentBalance
}

func (s *BankService) CreateExchangeRate(r dbank.ExchangeRate) (uuid.UUID, error) {
	newUuid := uuid.New()
	now := time.Now()

	exchangeRateOrm := db.BankExchangeRateOrm{
		ExchangeRateUuid:   newUuid,
		FromCurrency:       r.FromCurrency,
		ToCurrency:         r.ToCurrency,
		Rate:               r.Rate,
		ValidFromTimestamp: r.ValidFromTimestamp,
		ValidToTimestamp:   r.ValidToTimestamp,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	return s.db.CreateExchangeRate(exchangeRateOrm)
}

func (s *BankService) FindExchangeRate(fromCur string, toCur string, ts time.Time) float64 {
	exchangeRate, err := s.db.GetExchangeRateAtTimestamp(fromCur, toCur, ts)

	if err != nil {
		return 0
	}

	return float64(exchangeRate.Rate)
}

func (s *BankService) CreateTransaction(acct string, t dbank.Transaction) (uuid.UUID, error) {
	newUuid := uuid.New()
	now := time.Now()

	bankAccountOrm, err := s.db.GetBankAccountByAccountNumber(acct)

	if err != nil {
		log.Printf("Can't create transaction for %v : %v\n", acct, err)
		return uuid.Nil, err
	}

	transactionOrm := db.BankTransactionOrm{
		TransactionUuid:      newUuid,
		AccountUuid:          bankAccountOrm.AccountUuid,
		TransactionTimestamp: now,
		Amount:               t.Amount,
		TransactionType:      t.TransactionType,
		Notes:                t.Notes,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	savedUuid, err := s.db.CreateTransaction(bankAccountOrm, transactionOrm)

	return savedUuid, err
}

func (s *BankService) CalculateTransactionSummary(tcur *dbank.TransactionSummary, trans dbank.Transaction) error {
	switch trans.TransactionType {
	case dbank.TransactionTypeIn:
		tcur.SumIn += trans.Amount
	case dbank.TransactionTypeOut:
		tcur.SumOut += trans.Amount
	default:
		return fmt.Errorf("Unknown transaction type %v", trans.TransactionType)
	}

	tcur.SumTotal = tcur.SumIn - tcur.SumOut

	return nil
}

func (s *BankService) Transfer(tt dbank.TransferTransaction) (uuid.UUID, bool, error) {
	now := time.Now()

	fromAccountOrm, err := s.db.GetBankAccountByAccountNumber(tt.FromAccountNumber)

	if err != nil {
		log.Printf("Can't find transfer from account %v : %v\n", tt.FromAccountNumber, err)
		return uuid.Nil, false, err
	}

	toAccountOrm, err := s.db.GetBankAccountByAccountNumber(tt.ToAccountNumber)

	if err != nil {
		log.Printf("Can't find transfer to account %v : %v\n", tt.ToAccountNumber, err)
		return uuid.Nil, false, err
	}

	fromTransactionOrm := db.BankTransactionOrm{
		TransactionUuid:      uuid.New(),
		TransactionTimestamp: now,
		TransactionType:      dbank.TransactionTypeOut,
		AccountUuid:          fromAccountOrm.AccountUuid,
		Amount:               tt.Amount,
		Notes:                "Transfer out to " + tt.ToAccountNumber,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	toTransactionOrm := db.BankTransactionOrm{
		TransactionUuid:      uuid.New(),
		TransactionTimestamp: now,
		TransactionType:      dbank.TransactionTypeIn,
		AccountUuid:          toAccountOrm.AccountUuid,
		Amount:               tt.Amount,
		Notes:                "Transfer out from " + tt.FromAccountNumber,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	// create transfer request
	newTransferUuid := uuid.New()

	transferOrm := db.BankTransferOrm{
		TransferUuid:      newTransferUuid,
		FromAccountUuid:   fromAccountOrm.AccountUuid,
		ToAccountUuid:     toAccountOrm.AccountUuid,
		Currency:          tt.Currency,
		Amount:            tt.Amount,
		TransferTimestamp: now,
		TransferSuccess:   false,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if _, err := s.db.CreateTransfer(transferOrm); err != nil {
		log.Printf("Can't create teransfer from %v to %v\n", tt.FromAccountNumber, tt.ToAccountNumber, err)
		return uuid.Nil, false, err
	}

	if transferPairSuccess, err := s.db.CreateTransferTransactionPair(fromAccountOrm, toAccountOrm, fromTransactionOrm, toTransactionOrm); transferPairSuccess {
		s.db.UpdateTransferStatus(transferOrm, true)
		return newTransferUuid, true, nil
	} else {
		return newTransferUuid, false, err
	}
}
