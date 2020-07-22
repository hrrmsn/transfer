package handlers

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/go-openapi/strfmt"

	"wheely/test/internal/predict/models"
	"wheely/test/pkg/transfer/utils"
)

func readBody(body io.ReadCloser) (*models.Position, error) {
	defer body.Close()

	resp, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	pos := &models.Position{}
	if err = pos.UnmarshalBinary(resp); err != nil {
		return nil, err
	}

	return pos, nil
}

func readPos(request *http.Request, formats strfmt.Registry) (*models.Position, error) {
	pos, err := readBody(request.Body)
	if err != nil {
		return nil, utils.WrapError("Error when reading body of POST request", err)
	}

	if err = pos.Validate(formats); err != nil {
		return nil, utils.WrapError("Predict position is invalid", err)
	}

	return pos, nil
}
