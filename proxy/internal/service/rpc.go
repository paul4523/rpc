package service

import (
	"log"
	"net"
	"net/rpc"
	"proxy/internal/model"
)

type GeoProvider interface {
	AddressSearch(req model.RequestAddressSearch, res *model.ResponseAddress) error
	GeoCode(req model.RequestAddressGeocode, res *model.ResponseAddress) error
}
type RPCServicer struct {
	Service Service
}

func NewGeoProvider(service Service) GeoProvider {
	return &RPCServicer{
		Service: service,
	}
}

func (r *RPCServicer) AddressSearch(req model.RequestAddressSearch, res *model.ResponseAddress) error {
	addresses, err := r.Service.AddressSearch(req.Query)
	if err != nil {
		return err
	}
	res.Addresses = addresses
	return nil
}

func (r *RPCServicer) GeoCode(req model.RequestAddressGeocode, res *model.ResponseAddress) error {
	addresses, err := r.Service.GeoCode(req.Lat, req.Lon)
	if err != nil {
		return err
	}
	res.Addresses = addresses
	return nil
}

func RPCServe(service Service) {
	rpcService := NewGeoProvider(service)
	rpc.Register(rpcService)

	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("Ошибка при запуске сервера:", err)
	}

	log.Println("Сервер запущен на порту 1234")
	rpc.Accept(l)
}
func Clientgo() {
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
