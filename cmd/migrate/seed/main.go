package main

/* do this after the seed is completed
import (
	"log"

	"github.com/Raghav847/social-go/internal/db"
	"github.com/Raghav847/social-go/internal/env"
	"github.com/Raghav847/social-go/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/socialnetwork?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := store.NewStorage(conn)

	db.Seed(store, conn)
}
*/
