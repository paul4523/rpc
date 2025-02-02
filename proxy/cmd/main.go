package main

import (
	"log"
	"net"
	"net/rpc"
	_ "proxy/docs"
	"proxy/internal/controller"
	"proxy/internal/model"

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

	go func() {

		server.Serve()
	}()
	go func() {
		rpcService := &servicer.RPCServicer{Service: service}
		rpc.Register(rpcService)

		l, err := net.Listen("tcp", ":1234")
		if err != nil {
			log.Fatal("Ошибка при запуске сервера:", err)
		}

		log.Println("Сервер запущен на порту 1234")
		rpc.Accept(l)
	}()

	//клиент для теста
	go clientgo()
	go clientgo()

	server.Shutdown()
}
func clientgo() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Ошибка при подключении к RPC-серверу:", err)
	}
	defer client.Close()

	// Пример вызова метода AddressSearch
	req := model.RequestAddressSearch{Query: "Москва"}
	var res model.ResponseAddress
	err = client.Call("RPCServicer.AddressSearch", req, &res)
	if err != nil {
		log.Fatal("Ошибка при вызове метода AddressSearch:", err)
	}
	for _, addr := range res.Addresses {
		log.Printf("Адрес: %+v\n", addr)
	}

	// Пример вызова метода GeoCode
	geoReq := model.RequestAddressGeocode{Lat: "55.7558", Lon: "37.6173"}
	var geoRes model.ResponseAddress
	err = client.Call("RPCServicer.GeoCode", geoReq, &geoRes)
	if err != nil {
		log.Fatal("Ошибка при вызове метода GeoCode:", err)
	}
	for _, addr := range geoRes.Addresses {
		log.Printf("Адрес: %+v\n", addr)
	}
}
