package main

import (
	"fmt"
	"go_api/database"
	server "go_api/handlers"
	"log"
)

func main() {

	Store, err := database.NewPostgresDbConnection()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", Store)

	server := server.NewAPIServer(":3000", Store)
	server.Run()
}
