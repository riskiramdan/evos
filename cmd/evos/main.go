package main

import (
	"log"

	"github.com/riskiramdan/evos/config"
	"github.com/riskiramdan/evos/databases"
	"github.com/riskiramdan/evos/internal/data"
	"github.com/riskiramdan/evos/internal/hosts"
	internalhttp "github.com/riskiramdan/evos/internal/http"
	"github.com/riskiramdan/evos/internal/user"
	userPg "github.com/riskiramdan/evos/internal/user/postgres"
	"github.com/riskiramdan/evos/seeder"
	"github.com/riskiramdan/evos/util"

	"github.com/jmoiron/sqlx"
)

// InternalServices represents all the internal domain services
type InternalServices struct {
	userService user.ServiceInterface
}

func buildInternalServices(db *sqlx.DB) *InternalServices {
	userPostgresStorage := userPg.NewPostgresStorage(
		data.NewPostgresStorage(db, "users", user.Users{}),
	)
	userService := user.NewService(userPostgresStorage)
	return &InternalServices{
		userService: userService,
	}
}

func main() {

	config, err := config.GetConfiguration()
	if err != nil {
		log.Fatalln("failed to get configuration: ", err)
	}
	db, err := sqlx.Open("postgres", config.DBConnectionString)
	if err != nil {
		log.Fatalln("failed to open database x: ", err)
	}

	util := &util.Utility{}
	httpManager := &hosts.HTTPManager{}
	defer db.Close()
	dataManager := data.NewManager(db)
	internalServices := buildInternalServices(db)
	// Migrate the db
	databases.MigrateUp()
	// Seeder
	err = seeder.SeedUp()
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	s := internalhttp.NewServer(
		internalServices.userService,
		dataManager,
		config,
		util,
		httpManager,
	)
	s.Serve()
}
