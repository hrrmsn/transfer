package transfer

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"wheely/test/pkg/transfer/utils"
)

// mock cars service
type CarsServerMock struct {
	Data string
}

func (csm *CarsServerMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}

	switch r.URL.Path {
	case "/fake-eta/cars":
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(csm.Data))
	}
}

// mock predict service
type PredictServerMock struct {
	Data string
}

func (psm *PredictServerMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	switch r.URL.Path {
	case "/fake-eta/predict":
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(psm.Data))
	}
}

// TestCase struct is a helper
type TestCase struct {
	Name string
	*TestRequest
	*TestResponse
}

type TestRequest struct {
	TransferData string
	CarsData     string
	PredictData  string
}

type TestResponse struct {
	ExpectedStatus int
	ExpectedResult string // equals or contains
}

func (tc *TestCase) serverMocks() (*httptest.Server, *httptest.Server) {
	carsServerMock := httptest.NewServer(&CarsServerMock{Data: tc.CarsData})
	predictServerMock := httptest.NewServer(&PredictServerMock{Data: tc.PredictData})
	return carsServerMock, predictServerMock
}

func (tc *TestCase) config(carsSrvUrl, predictSrvUrl string) *utils.Config {
	cfgTest, _ := utils.NewConfig()

	cfgTest.CarsConfig.Host = strings.TrimPrefix(carsSrvUrl, "http://")
	cfgTest.CarsConfig.Schemes = []string{"http"}

	cfgTest.PredictConfig.Host = strings.TrimPrefix(predictSrvUrl, "http://")
	cfgTest.PredictConfig.Schemes = []string{"http"}

	return cfgTest
}

func (tc *TestCase) check(
	resp []byte,
	rec *httptest.ResponseRecorder,
	testName string,
	t *testing.T,
) bool {

	passed := true
	if status := rec.Code; status != tc.TestResponse.ExpectedStatus {
		t.Errorf(
			"Test '%s': wrong status request: got %d want %d\n",
			testName,
			status,
			tc.TestResponse.ExpectedStatus,
		)
		passed = false
	}

	if result := string(resp); result != tc.TestResponse.ExpectedResult &&
		!strings.Contains(result, tc.TestResponse.ExpectedResult) {

		t.Errorf(
			"Test '%s': wrong response: got %s want %s\n",
			testName,
			result,
			tc.TestResponse.ExpectedResult,
		)
		passed = false
	}

	return passed
}

func (tc *TestCase) run(testName string, t *testing.T) bool {
	carsServerMock, predictServerMock := tc.serverMocks()

	defer func() {
		carsServerMock.Close()
		predictServerMock.Close()
	}()

	cfgTest := tc.config(carsServerMock.URL, predictServerMock.URL)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/transfer", strings.NewReader(tc.TransferData))

	transferServer := NewServer(cfgTest)
	transferServer.Handler.ServeHTTP(rec, req)

	defer rec.Result().Body.Close()
	resp, _ := ioutil.ReadAll(rec.Result().Body)

	if passed := tc.check(resp, rec, testName, t); passed {
		return true
	}
	return false
}

// init test cases
var testCases []*TestCase

func initTestCases() {
	testCases = make([]*TestCase, 0)

	testCases = append(testCases, &TestCase{
		Name: "test route handler",
		TestRequest: &TestRequest{
			TransferData: `{"lat": 17.986511, "lng": 63.441092}`,
			CarsData: "[" +
				`{"id":16,"lat":55.7575429,"lng":37.6135117},` +
				`{"id":229,"lat":55.74837156167371,"lng":37.61180107665421},` +
				`{"id":8,"lat":55.7532706,"lng":37.6076902}` +
				"]",
			PredictData: `[7,2,3]`,
		},
		TestResponse: &TestResponse{
			ExpectedStatus: http.StatusOK,
			ExpectedResult: `{"response": 2}`,
		},
	})

	testCases = append(testCases, &TestCase{
		Name: "test invalid input data",
		TestRequest: &TestRequest{
			TransferData: `{"lat": -200, "lng": 63.441092}`,
			CarsData:     `[{"id":16,"lat":55.7575429,"lng":37.6135117}]`,
			PredictData:  `[7,2,3]`,
		},
		TestResponse: &TestResponse{
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedResult: `input data error: predict position is invalid`,
		},
	})

	testCases = append(testCases, &TestCase{
		Name: "test cars data invalid",
		TestRequest: &TestRequest{
			TransferData: `{"lat": 17.986511, "lng": 63.441092}`,
			CarsData:     `[{"id":16,"lat":-200,"lng":37.6135117}]`,
			PredictData:  `[7,2,3]`,
		},
		TestResponse: &TestResponse{
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedResult: `cars data is invalid`,
		},
	})

	testCases = append(testCases, &TestCase{
		Name: "test receiving cars data error",
		TestRequest: &TestRequest{
			TransferData: `{"lat": 17.986511, "lng": 63.441092}`,
			CarsData:     `[{"id":16,,"lat":55.7575429,"lng":37.6135117}]`,
			PredictData:  `[7,2,3]`,
		},
		TestResponse: &TestResponse{
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedResult: `error when receiving data from cars service`,
		},
	})

	testCases = append(testCases, &TestCase{
		Name: "test receiving predict data error",
		TestRequest: &TestRequest{
			TransferData: `{"lat": 17.986511, "lng": 63.441092}`,
			CarsData:     `[{"id":16,"lat":55.7575429,"lng":37.6135117}]`,
			PredictData:  `[7,2,,3]`,
		},
		TestResponse: &TestResponse{
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedResult: `error when receiving data from predict service`,
		},
	})
}

// run tests
func TestRoute(t *testing.T) {
	initTestCases()

	for _, testCase := range testCases {
		if passed := testCase.run(testCase.Name, t); passed {
			t.Logf("Test '%s' is passed\n", testCase.Name)
		} else {
			t.Logf("Test '%s' isn't passed\n", testCase.Name)
		}
	}
}
