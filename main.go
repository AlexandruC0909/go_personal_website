package main

import (
    "fmt"
    "go_api/database"
    server "go_api/handlers"
    "log"
)

func main() {
    // Run database migrations
    if err := migrateDatabase(); err != nil {
        log.Fatal(err)
    }

    store, err := database.NewPostgresDbConnection()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%+v\n", store)

    server := server.NewAPIServer(":3000", store)
    server.Run()
}

func migrateDatabase() error {
    store, err := database.NewPostgresDbConnection()
    if err != nil {
        return err
    }

    defer store.CloseDB()

    if err := store.migrate.Up(); err != nil && err != migrate.ErrNoChange {
        return err
    }

    fmt.Println("Database migrations completed.")
    return nil
}
