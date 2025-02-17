package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	dbmigration "github.com/shch989/my-grpc-go-server/db"

	mydb "github.com/shch989/my-grpc-go-server/internal/adapter/database"
	mygrpc "github.com/shch989/my-grpc-go-server/internal/adapter/grpc"
	app "github.com/shch989/my-grpc-go-server/internal/application"
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
	runDummyOrm(databaseAdapter)

	hs := &app.HelloService{}
	bs := &app.BankService{}

	grpcAdapter := mygrpc.NewGrpcAdapter(hs, bs, 9090)

	grpcAdapter.Run()

}

func runDummyOrm(da *mydb.DatabaseAdapter) {
	now := time.Now()

	uuid, _ := da.Save(
		&mydb.DummyOrm{
			UserId:    uuid.New(),
			UserName:  "Time " + time.Now().Format("15:04:03"),
			CreatedAt: now,
			UpdatedAt: now,
		},
	)

	res, _ := da.GetByUuid(&uuid)

	log.Println("res :", res)
}
