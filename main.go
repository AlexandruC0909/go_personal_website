package main

import (
	"fmt"
	"go_api/database"
	server "go_api/handlers"
	"log"
)

func main() {

	store, err := database.NewPostgresDbConnection()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", store)

	server := server.NewAPIServer(":3000", store)
	server.Run()
}
