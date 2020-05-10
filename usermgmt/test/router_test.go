package test


import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"bytes"
	// "fmt"
	_ "github.com/lib/pq"
	"github.com/leagueauctions/usermgmt"
	"github.com/leagueauctions/router"
)


func initDBAndRouter(t *testing.T) router.Wrapper{
	db, err := openPostgreDatabase("postgres", "postgres", "leagueauction")
	if err != nil{
		t.Fatal(err)
	}
	err = clearUserTable(t, db)
	if (err != nil){
		t.Fatal(err)
	}

	var r router.Wrapper = new(router.MuxWrapper)
	routerCfg := router.Config{HostAddress: "localhost", PortNo : 8081}
	err = r.Init(routerCfg)
	if (err != nil){
		t.Fatal(err)
	}

	usrMgmtRouter := new(usermgmt.Router)
	err = usrMgmtRouter.Init(r, db)
	if (err != nil){
		t.Fatal(err)
	}
	return r
}

func executeRequest(r router.Wrapper, req *http.Request) *httptest.ResponseRecorder {
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


func TestRegisterNewUser(t *testing.T){

	r := initDBAndRouter(t)
	
	var RegistrationJSONReqStr = []byte(`{"user_id":"x@x.com", "password": "pwd123"}`)
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
	if userInfoResponse.UserSerialID < 0 || userInfoResponse.IsActive == false{
		t.Fatal("Expected valid user info. Actual: ", userInfoResponse)
	}
	// fmt.Println("userInfoResponse", userInfoResponse)
}

//TODO: write more tests for negative cases

//TODO:refactor error response - make life easy for client so that it does not has to rememeber
//		to use different response(map[string]string) for error handling


//TODO: create testing stub for activation code
//probably a functor for creating and verifying user activation code
//also another functor for creating and using JWT could prove handy

//TODO: check if mocking the db is an easy option for router test, right now router testing is too heavy
