package main

import (
	"fmt"
	"go_api/postgres"
	"go_api/server"
	"log"
)

func main() {
	store, err := postgres.NewPostgressStore()
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