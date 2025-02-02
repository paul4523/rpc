package service

import "proxy/internal/model"

type RPCServicer struct {
	Service Service
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
