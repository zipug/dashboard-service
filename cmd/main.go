package main

import (
	"dashboard/internal/application"
	"dashboard/internal/presentation/http/chi"
)

func main() {
	app := application.DashboardAppService

	server := chi.NewHttpServer(app)
	server.Start()
}
