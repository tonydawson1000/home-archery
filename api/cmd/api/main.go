package main

import (
	"log"
	"net/http"
	"os"

	httpapi "github.com/tonydawson1000/home-archery/api/internal/adapters/http"
	"github.com/tonydawson1000/home-archery/api/internal/application"
	"github.com/tonydawson1000/home-archery/api/internal/application/memory"
)

func main() {
	store := memory.NewStore()
	svc := &application.Service{
		Archers:  store.Archers(),
		Sessions: store.Sessions(),
		Clock:    application.SystemClock{},
		IDs:      application.RandomIDs{},
	}

	addr := ":8080"
	if v := os.Getenv("HTTP_ADDR"); v != "" {
		addr = v
	}
	log.Printf("home-archery API (in-memory) listening on %s", addr)
	if err := http.ListenAndServe(addr, httpapi.NewHandler(svc)); err != nil {
		log.Fatal(err)
	}
}
