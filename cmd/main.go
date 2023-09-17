package main

import (
	"fmt"
	database "go_api/database"
	server "go_api/handlers"
	"log"
)

func main() {
	store, err := database.NewDbConnection()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", store)
	server := server.NewAPIServer(":3000", store)
	server.Run()
}
