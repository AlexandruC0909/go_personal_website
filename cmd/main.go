package main

import (
	"fmt"
	database "go_api/database"
	server "go_api/handlers"
	"log"
	"os"
)

func main() {
	store, err := database.NewPostgresDbConnection()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := store.Init(); err != nil {
		log.Fatal(err)
	}
	os.Setenv("JWT_SECRET", "9V7$2kP&6a#R@5bT1yZ!8wG*4qS%F3eU")
	fmt.Printf("%+v\n", store)
	server := server.NewAPIServer(":3000", store)
	server.Run()

}
