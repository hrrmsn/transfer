package handlers

import (
	"reflect"

	cars_ops "wheely/test/internal/cars/client/operations"
	predict_ops "wheely/test/internal/predict/client/operations"
	predict_mods "wheely/test/internal/predict/models"
)

var (
	predictOKTest = &predict_ops.PredictOK{
		Payload: []int64{1, 2},
	}
)

type PredictClientMock struct{}

func (pcm *PredictClientMock) GetPredict(
	pos *predict_mods.Position,
	carsData *cars_ops.GetCarsOK,
) (*predict_ops.PredictOK, error) {

	if reflect.DeepEqual(pos, posTest) && reflect.DeepEqual(carsData, getCarsOKTest) {
		return predictOKTest, nil
	}

	return &predict_ops.PredictOK{}, nil
}

func (pcm *PredictClientMock) Healthy() bool {
	return true
}
