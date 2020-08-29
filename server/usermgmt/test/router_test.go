package test


import (
	"github.com/leagueauctions/server/database"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"bytes"
	// "fmt"
	_ "github.com/lib/pq"
	"github.com/leagueauctions/server/usermgmt"
	"github.com/leagueauctions/server/libs/router"
)


func initDBAndRouter(t *testing.T) *router.MuxWrapper{

	//use mock database
	mockUserstore := database.GetMockUserStore()

	var r *router.MuxWrapper = new(router.MuxWrapper)
	routerCfg := router.Config{
		HostAddress: "localhost", 
		PortNo : 8081, 
		Secure: true,
		CertFilePath : "../../certs/cert.pem",
		KeyPath : "../../certs/key.pem",
	}
	err := r.Init(routerCfg)
	if (err != nil){
		t.Fatal(err)
	}

	usrMgmtRouter := new(usermgmt.Router)
	err = usrMgmtRouter.Init(r, mockUserstore)
	if (err != nil){
		t.Fatal(err)
	}
	return r
}

func executeRequest(r *router.MuxWrapper, req *http.Request) *httptest.ResponseRecorder {
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    return rr
}

func TestEmptyTokenRoute(t *testing.T) {

	r := initDBAndRouter(t)
	
	var jsonReqStr = []byte(`{"user_id":"x@x.com", "token_string": ""}`)
	req, _ := http.NewRequest("GET", "/user/info", bytes.NewBuffer(jsonReqStr))
	response := executeRequest(r, req)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "invalid token"{
		t.Fatal("Token must be invalidated. Actual error: ", m["error"])
	}
}


func TestRegisterActivationLoginExistingUser(t *testing.T){

	r := initDBAndRouter(t)
	
	var RegistrationJSONReqStr = []byte(`{"user_id":"x@x.com", "user_password": "pwd123"}`)
	req, _ := http.NewRequest("POST", "/user/register", bytes.NewBuffer(RegistrationJSONReqStr))
	response := executeRequest(r, req)
	if response.Code != http.StatusOK{
		t.Fatal("Actual code ", response.Code)
	}
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["status"] != "awaiting verification"{
		t.Fatal("Actual status ", m["status"])
	}

	var activationJSONRequest = []byte(`{"user_id":"x@x.com", "user_activation_code": "123456"}`)
	req, _ = http.NewRequest("POST", "/user/activation", bytes.NewBuffer(activationJSONRequest))
	response = executeRequest(r, req)
	if response.Code != http.StatusOK{
		var m map[string]string
		json.Unmarshal(response.Body.Bytes(), &m)
		t.Error("Actual code ", response.Code)
		t.Fatal("Actual error ", m["error"])
	}
	var loginResponse usermgmt.LogInResponse
	json.Unmarshal(response.Body.Bytes(), &loginResponse)
	if loginResponse.Token == "" || loginResponse.Expiry == ""{
		t.Fatal("Expected valid login response. Actual: ", loginResponse)
	}

	var loginJSONRequest = []byte(`{"user_id":"x@x.com", "user_password": "pwd123"}`)
	req, _ = http.NewRequest("POST", "/user/login", bytes.NewBuffer(loginJSONRequest))
	response = executeRequest(r, req)
	if response.Code != http.StatusOK {
		var m map[string]string
		json.Unmarshal(response.Body.Bytes(), &m)
		t.Error("Actual code ", response.Code)
		t.Fatal("Actual error ", m["error"])
	}
	json.Unmarshal(response.Body.Bytes(), &loginResponse)
	if loginResponse.Token == "" || loginResponse.Expiry == ""{
		t.Fatal("Expected valid login response. Actual: ", loginResponse)
	}

	getUserInfoRequestStr := "{\"user_id\":\"x@x.com\", \"token_string\": \""+loginResponse.Token+"\"}"
	var getUserInfoRequest = []byte(getUserInfoRequestStr)
	req, _ = http.NewRequest("GET", "/user/info", bytes.NewBuffer(getUserInfoRequest))
	response = executeRequest(r, req)
	if response.Code != http.StatusOK{
		var m map[string]string
		json.Unmarshal(response.Body.Bytes(), &m)
		t.Error("Actual code ", response.Code)
		t.Fatal("Actual error ", m["error"])
	}

	var userInfoResponse usermgmt.UserInfoResponse
	json.Unmarshal(response.Body.Bytes(), &userInfoResponse)
	if userInfoResponse.UserSerialID == "" || userInfoResponse.IsActive == false{
		t.Fatal("Expected valid user info. Actual: ", userInfoResponse)
	}
	// fmt.Println("userInfoResponse", userInfoResponse)
}


func TestInvalidActivationCode(t *testing.T){

	r := initDBAndRouter(t)
	
	var RegistrationJSONReqStr = []byte(`{"user_id":"x@x.com", "user_password": "pwd123"}`)
	req, _ := http.NewRequest("POST", "/user/register", bytes.NewBuffer(RegistrationJSONReqStr))
	response := executeRequest(r, req)

	var activationJSONRequest = []byte(`{"user_id":"x@x.com", "user_activation_code": "wrong_code"}`)
	req, _ = http.NewRequest("POST", "/user/activation", bytes.NewBuffer(activationJSONRequest))
	response = executeRequest(r, req)
	if response.Code != http.StatusUnauthorized {
		var m map[string]string
		json.Unmarshal(response.Body.Bytes(), &m)
		t.Error("Actual error ", m["error"])
		t.Fatal("Actual code ", response.Code, "Expected" , http.StatusUnauthorized)
	}
}

func TestInvalidUserPassword(t *testing.T){

	r := initDBAndRouter(t)
	
	var RegistrationJSONReqStr = []byte(`{"user_id":"x@x.com", "user_password": "pwd123"}`)
	req, _ := http.NewRequest("POST", "/user/register", bytes.NewBuffer(RegistrationJSONReqStr))
	response := executeRequest(r, req)

	var activationJSONRequest = []byte(`{"user_id":"x@x.com", "user_activation_code": "123456"}`)
	req, _ = http.NewRequest("POST", "/user/activation", bytes.NewBuffer(activationJSONRequest))
	response = executeRequest(r, req)
	
	var loginJSONRequest = []byte(`{"user_id":"x@x.com", "user_password": "pwd1234"}`)
	req, _ = http.NewRequest("POST", "/user/login", bytes.NewBuffer(loginJSONRequest))
	response = executeRequest(r, req)
	if response.Code != http.StatusUnauthorized {
		var m map[string]string
		json.Unmarshal(response.Body.Bytes(), &m)
		t.Error("Actual error ", m["error"])
		t.Fatal("Actual code ", response.Code, "Expected" , http.StatusUnauthorized)
	}
}


func TestLoginWithoutActivation(t *testing.T){

	r := initDBAndRouter(t)
	
	var RegistrationJSONReqStr = []byte(`{"user_id":"x@x.com", "user_password": "pwd123"}`)
	req, _ := http.NewRequest("POST", "/user/register", bytes.NewBuffer(RegistrationJSONReqStr))
	response := executeRequest(r, req)

	var activationJSONRequest = []byte(`{"user_id":"x@x.com", "user_activation_code": "123456"}`)
	req, _ = http.NewRequest("POST", "/user/activation", bytes.NewBuffer(activationJSONRequest))
	response = executeRequest(r, req)
	
	req, _ = http.NewRequest("POST", "/user/register", bytes.NewBuffer(RegistrationJSONReqStr))
	response = executeRequest(r, req)
	if response.Code != http.StatusCreated {
		t.Fatal("should get created response once user activation is done")
	}
}


func TestAlreadyRegisteredUser(t *testing.T){

	r := initDBAndRouter(t)
	
	var RegistrationJSONReqStr = []byte(`{"user_id":"x@x.com", "user_password": "pwd123"}`)
	req, _ := http.NewRequest("POST", "/user/register", bytes.NewBuffer(RegistrationJSONReqStr))
	response := executeRequest(r, req)

	var loginJSONRequest = []byte(`{"user_id":"x@x.com", "user_password": "pwd123"}`)
	req, _ = http.NewRequest("POST", "/user/login", bytes.NewBuffer(loginJSONRequest))
	response = executeRequest(r, req)
	if response.Code != http.StatusForbidden {
		t.Fatal("Expected activation error")
	}
}

var ErrorRequestTests = []struct {
	reqquestURI		string
	reqquestMethod	string
	payloadJSON		[]byte
	expectedStatus  int
}{
	{"/user/register", "POST", []byte(`{"id":"x@x.com", "user_password": "pwd123"}`), http.StatusBadRequest},
	{"/user/register", "POST", []byte(`{"user_id":"x@x.com", "password": "pwd123"}`), http.StatusBadRequest},
	{"/user/activation", "POST", []byte(`{"user_id":"x@x.com", "code": "123456"}`), http.StatusBadRequest},
	{"/user/activation", "POST", []byte(`{"id":"x@x.com", "user_activation_code": "123456"}`), http.StatusBadRequest},
	{"/user/login", "POST", []byte(`{"id":"x@x.com", "user_password": "pwd123"}`), http.StatusBadRequest},
	{"/user/login", "POST", []byte(`{"user_id":"x@x.com", "password": "pwd123"}`), http.StatusBadRequest},

	{"/user/login", "POST", []byte(`{"user_id":"x@x.com", "user_password": "pwd123"}`), http.StatusNotFound},
	{"/user/activation", "POST", []byte(`{"user_id":"x@x.com", "user_activation_code": "123456"}`), http.StatusNotFound},

}
func TestErrorRequests(t *testing.T){

	r := initDBAndRouter(t)
	
	for _, test := range ErrorRequestTests{
		req, _ := http.NewRequest(test.reqquestMethod, test.reqquestURI, bytes.NewBuffer(test.payloadJSON))
		response := executeRequest(r, req)
		if response.Code != test.expectedStatus{
			t.Fatal("Expected code ", test.expectedStatus, "Actual code", response.Code)
		}
	}

}

//TODO: write more tests for negative cases

//TODO:refactor error response - make life easy for client so that it does not has to rememeber
//		to use different response(map[string]string) for error handling


//TODO: create testing stub for activation code
//probably a functor for creating and verifying user activation code
//also another functor for creating and using JWT could prove handy
