package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	emphttp "github.com/arseniizyk/internal/employees/delivery/http"
	emprepo "github.com/arseniizyk/internal/employees/repository/postgres"
	empusecase "github.com/arseniizyk/internal/employees/usecase"
	"github.com/gin-gonic/gin"

	"github.com/arseniizyk/internal/employees"
	pg "github.com/arseniizyk/internal/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	empUC employees.EmployeeUsecase
	pool  *pgxpool.Pool
}

func NewApp() *App {
	pool := initDB()

	empRepo := emprepo.New(pool)

	return &App{
		empUC: empusecase.New(empRepo),
		pool:  pool,
	}
}

func (a *App) Run() error {
	r := gin.Default()

	api := r.Group("/api/v1")
	emphttp.RegisterEmployeesEndpoints(api, a.empUC)

	httpServer := http.Server{
		Handler:      r,
		Addr:         ":" + os.Getenv("SERVER_PORT"),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	a.pool.Close()

	return httpServer.Shutdown(ctx)
}

func initDB() *pgxpool.Pool {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_SSL"),
	)

	pool, err := pg.OpenPool(dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}

	return pool
}
