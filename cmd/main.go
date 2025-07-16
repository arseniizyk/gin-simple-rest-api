package main

import (
	"github.com/arseniizyk/internal/handlers"
	"github.com/arseniizyk/internal/routes"
	"github.com/arseniizyk/internal/storage"
)


func main() {
	s := storage.NewMemoryStorage()
	h := handlers.New(s)

	r := routes.Setup(h)
	r.Run(":8080")
}
