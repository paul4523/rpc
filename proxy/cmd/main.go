package main

import (
	_ "proxy/docs"
	"proxy/internal/controller"

	servicer "proxy/internal/service"
)

// @title Swagger GeoService
// @version 1.0
// @description GeoService from KataAcademy
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	service := servicer.New()
	handler := controller.Handler{Responder: &service}
	r := controller.NewRouter(&handler)
	server := controller.NewServer(":8080", r)

	go server.Serve()
	go servicer.RPCServe(service)

	//клиент для теста
	go servicer.Clientgo()
	go servicer.Clientgo()

	server.Shutdown()
}
