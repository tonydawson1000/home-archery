package main

import (
	"context"
	"log"
	"net/http"
	"os"

	httpapi "github.com/tonydawson1000/home-archery/api/internal/adapters/http"
	"github.com/tonydawson1000/home-archery/api/internal/adapters/postgres"
	"github.com/tonydawson1000/home-archery/api/internal/application"
	"github.com/tonydawson1000/home-archery/api/internal/application/memory"
)

func main() {
	ctx := context.Background()
	svc, backend, cleanup := mustService(ctx)
	defer cleanup()

	addr := ":8080"
	if v := os.Getenv("HTTP_ADDR"); v != "" {
		addr = v
	}
	log.Printf("home-archery API (%s) listening on %s", backend, addr)
	if err := http.ListenAndServe(addr, httpapi.NewHandler(svc)); err != nil {
		log.Fatal(err)
	}
}

func mustService(ctx context.Context) (*application.Service, string, func()) {
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		store, err := postgres.Open(ctx, dsn)
		if err != nil {
			log.Fatal(err)
		}
		return &application.Service{
			Archers:  store.Archers(),
			Sessions: store.Sessions(),
			Clock:    application.SystemClock{},
			IDs:      application.RandomIDs{},
		}, "postgres", store.Close
	}

	store := memory.NewStore()
	return &application.Service{
		Archers:  store.Archers(),
		Sessions: store.Sessions(),
		Clock:    application.SystemClock{},
		IDs:      application.RandomIDs{},
	}, "in-memory", func() {}
}
