package handlers

import (
	"reflect"

	cars_ops "wheely/test/internal/cars/client/operations"
	cars_mods "wheely/test/internal/cars/models"
	predict_mods "wheely/test/internal/predict/models"
)

var (
	posTest = &predict_mods.Position{
		Lat: 12.34,
		Lng: 56.78,
	}

	getCarsOKTest = &cars_ops.GetCarsOK{
		Payload: []cars_mods.Car{
			cars_mods.Car{
				ID:  1,
				Lat: 42.17,
				Lng: 59.63,
			},
			cars_mods.Car{
				ID:  2,
				Lat: 40.71,
				Lng: 53.04,
			},
		},
	}
)

type CarsClientMock struct{}

func (ccm *CarsClientMock) GetCars(pos *predict_mods.Position) (*cars_ops.GetCarsOK, error) {
	if reflect.DeepEqual(pos, posTest) {
		return getCarsOKTest, nil
	}
	return &cars_ops.GetCarsOK{}, nil
}

func (ccm *CarsClientMock) Validate(carsData *cars_ops.GetCarsOK) error {
	return nil
}

func (ccm *CarsClientMock) Healthy() bool {
	return true
}
