package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"proxy/internal/model"

	"github.com/ekomobile/dadata/v2/api/suggest"
	"github.com/ekomobile/dadata/v2/client"
)

type GeoServicer interface {
	AddressSearch(input string) ([]*model.Address, error)
	GeoCode(lat, lng string) ([]*model.Address, error)
}

type Service struct {
	api       *suggest.Api
	apiKey    string
	secretKey string
}

func (s *Service) AddressSearch(input string) ([]*model.Address, error) {

	var res []*model.Address
	rawRes, err := s.api.Address(context.Background(), &suggest.RequestParams{Query: input})
	if err != nil {
		return nil, err
	}
	for _, r := range rawRes {
		if r.Data.City == "" || r.Data.Street == "" {
			continue
		}
		res = append(res, &model.Address{City: r.Data.City, Street: r.Data.Street, House: r.Data.House, Lat: r.Data.GeoLat, Lon: r.Data.GeoLon})
	}
	return res, nil
}

func (s *Service) GeoCode(lat, lng string) ([]*model.Address, error) {

	httpClient := &http.Client{}
	var data = strings.NewReader(fmt.Sprintf(`{"lat": %s, "lon": %s}`, lat, lng))
	req, err := http.NewRequest("POST", "https://suggestions.dadata.ru/suggestions/api/4_1/rs/geolocate/address", data)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", s.apiKey))
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var geoCode model.GeoCode
	err = json.NewDecoder(resp.Body).Decode(&geoCode)
	if err != nil {
		return nil, err
	}
	var res []*model.Address
	for _, r := range geoCode.Suggestions {
		var address model.Address
		address.City = string(r.Data.City)
		address.Street = string(r.Data.Street)
		address.House = r.Data.House
		address.Lat = r.Data.GeoLat
		address.Lon = r.Data.GeoLon
		res = append(res, &address)
	}
	return res, nil
}

func New() Service {
	endpointUrl, _ := url.Parse("https://suggestions.dadata.ru/suggestions/api/4_1/rs/")
	creds := client.Credentials{
		ApiKeyValue:    "a24f884fd9737d9d7604e288667bc4719b6bd9b6",
		SecretKeyValue: "a455650abcfca33d33239b6a286784586a6a0b71",
	}
	api := suggest.Api{
		Client: client.NewClient(endpointUrl, client.WithCredentialProvider(&creds)),
	}

	service := Service{
		api:       &api,
		apiKey:    "a24f884fd9737d9d7604e288667bc4719b6bd9b6",
		secretKey: "a455650abcfca33d33239b6a286784586a6a0b71",
	}
	return service
}
