package main

import (
	"database/sql"
	"log"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
	dbmigration "github.com/shch989/my-grpc-go-server/db"

	mydb "github.com/shch989/my-grpc-go-server/internal/adapter/database"
	mygrpc "github.com/shch989/my-grpc-go-server/internal/adapter/grpc"
	app "github.com/shch989/my-grpc-go-server/internal/application"
	"github.com/shch989/my-grpc-go-server/internal/application/domain/bank"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(log.Writer())

	sqlDB, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5555/grpc?sslmode=disable")

	if err != nil {
		log.Fatalln("Can't connect database :", err)
	}

	dbmigration.Migrate(sqlDB)

	databaseAdapter, err := mydb.NewDatabaseAdapter(sqlDB)

	if err != nil {
		log.Fatalln("Can't create database adapter :", err)
	}

	// runDummyOrm(databaseAdapter)

	hs := &app.HelloService{}
	// bs := &app.BankService{}
	bs := app.NewBankService(databaseAdapter)

	go generateExchangeRates(bs, "USD", "IDR", 5*time.Second)

	grpcAdapter := mygrpc.NewGrpcAdapter(hs, bs, 9090)

	grpcAdapter.Run()

}

// func runDummyOrm(da *mydb.DatabaseAdapter) {
// 	now := time.Now()

// 	uuid, _ := da.Save(
// 		&mydb.DummyOrm{
// 			UserId:    uuid.New(),
// 			UserName:  "Time " + time.Now().Format("15:04:03"),
// 			CreatedAt: now,
// 			UpdatedAt: now,
// 		},
// 	)

// 	res, _ := da.GetByUuid(&uuid)

// 	log.Println("res :", res)
// }

func generateExchangeRates(bs *app.BankService, fromCurrency, toCurrency string, duration time.Duration) {
	ticker := time.NewTicker(duration)

	for range ticker.C {
		now := time.Now()
		validFrom := now.Truncate(time.Second).Add(3 * time.Second)
		validTo := validFrom.Add(duration).Add(-1 * time.Millisecond)

		dummyRate := bank.ExchangeRate{
			FromCurrency:       fromCurrency,
			ToCurrency:         toCurrency,
			ValidFromTimestamp: validFrom,
			ValidToTimestamp:   validTo,
			Rate:               2000 + float64(rand.Intn(300)),
		}

		bs.CreateExchangeRate(dummyRate)
	}
}
