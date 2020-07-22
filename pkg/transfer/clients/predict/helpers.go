package predict

import (
	"github.com/go-openapi/strfmt"

	cars_ops "wheely/test/internal/cars/client/operations"
	"wheely/test/internal/predict/models"
	"wheely/test/pkg/transfer/utils"
)

func carsToPositions(carsData *cars_ops.GetCarsOK, formats strfmt.Registry) ([]models.Position, error) {
	predictPositions := make([]models.Position, len(carsData.Payload))
	var err error

	for i := 0; i < len(predictPositions); i++ {
		var car = carsData.Payload[i]
		predictPositions[i] = models.Position{
			Lat: car.Lat,
			Lng: car.Lng,
		}
		if err = predictPositions[i].Validate(formats); err != nil {
			return nil, utils.WrapError("Predict position is invalid", err)
		}
	}

	return predictPositions, nil
}
