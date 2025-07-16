package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arseniizyk/internal/handlers"
	"github.com/arseniizyk/internal/routes"
	"github.com/arseniizyk/internal/storage"
	"github.com/arseniizyk/internal/storage/pg"
)

func main() {
	dsn := fmt.Sprintf("postgres://postgres:%s@postgres:5432/postgres?sslmode=disable", os.Getenv("POSTGRES_PASSWORD"))
	pool, err := pg.OpenPool(dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	s := storage.New(pool)
	h := handlers.New(s)

	r := routes.Setup(h)

	srv := http.Server{
		Handler:      r,
		Addr:         ":8000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Shutting down...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
