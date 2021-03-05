package main

import (
	"log"

	"github.com/riskiramdan/evos/databases"
	"github.com/riskiramdan/evos/seeder"
)

func main() {
	databases.MigrateUp()
	err := seeder.SeedUp()
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
}
