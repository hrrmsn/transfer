package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"wheely/test/pkg/transfer/utils"
)

var (
	cfgTest, _ = utils.NewConfig()
)

func TestRouteEmptyBody(t *testing.T) {
	req, err := http.NewRequest("POST", "/transfer", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	routeHandler := NewRouteHandler(cfgTest)

	routeHandler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusInternalServerError {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	handlerResponse := rec.Body.String()
	containsExpected := "input data error"
	if !strings.Contains(handlerResponse, containsExpected) {
		t.Errorf("Incorrect error message: got %v want %v", handlerResponse, containsExpected)
	}
}

func TestRouteWithBody(t *testing.T) {
	req, err := http.NewRequest(
		"POST",
		"/transfer",
		strings.NewReader(`{"lat": 12.34, "lng": 56.78}`),
	)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	routeHandler := NewRouteHandler(cfgTest)

	routeHandler.CarsClient = &CarsClientMock{}
	routeHandler.PredictClient = &PredictClientMock{}

	routeHandler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}

	handlerResponse, expected := rec.Body.String(), `{"response": "1"}`
	if handlerResponse != expected {
		t.Errorf("Incorrect error message: got %v want %v", handlerResponse, expected)
	}
}
