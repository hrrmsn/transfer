package transfer

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"wheely/test/pkg/transfer/utils"
)

type CarsServerMock struct {
}

func (csm *CarsServerMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("CarsServiceMock call")

	if r.Method != "GET" {
		// DEBUG
		log.Println("only GET request to cars service allowed")
		return
	}

	switch r.RequestURI {
	case "/fake-eta/cars":
		lat, lng, limit := r.FormValue("lat"), r.FormValue("lng"), r.FormValue("limit")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(``))

		// DEBUG
		log.Printf("lat=%s, lng=%s, limit=%s\n", lat, lng, limit)
	case "/fake-eta/_health":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(``))
	}
}

type PredictServerMock struct {
}

func (psm *PredictServerMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("PredictServerMock call")

	if r.Method != "POST" {
		// DEBUG
		log.Println("only POST requests to predict server allowed")
		return
	}

	switch r.RequestURI {
	case "fake-eta/predict":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(``))

		// DEBUG
		defer r.Body.Close()

		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(string(payload))
	case "/fake-eta/_health":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(``))
	}
}

func TestRoute(t *testing.T) {
	log.Println("TestRoute started")

	carsServiceMock := httptest.NewServer(&CarsServerMock{})
	predictServiceMock := httptest.NewServer(&PredictServerMock{})

	defer func() {
		carsServiceMock.Close()
		predictServiceMock.Close()
	}()

	cfgTest, _ := utils.NewConfig()

	cfgTest.CarsConfig.Host = strings.TrimPrefix(carsServiceMock.URL, "http://")
	cfgTest.CarsConfig.Schemes = []string{"http"}

	cfgTest.PredictConfig.Host = strings.TrimPrefix(predictServiceMock.URL, "http://")
	cfgTest.PredictConfig.Schemes = []string{"http"}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/transfer", strings.NewReader(`{"lat": 17.986511, "lng": 63.441092}`))

	transferServer := NewServer(cfgTest)
	transferServer.Handler.ServeHTTP(rec, req)

	resp, _ := ioutil.ReadAll(rec.Result().Body)
	fmt.Println("response body -> ", string(resp))

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("Something wrong: status is %d\n", status)
	}
}
