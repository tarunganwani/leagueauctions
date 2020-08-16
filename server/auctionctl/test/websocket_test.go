package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"bytes"
	"strings"
	// "fmt"
	_ "github.com/lib/pq"
	"github.com/leagueauctions/server/usermgmt"
	"github.com/leagueauctions/server/libs/router"
	"github.com/leagueauctions/server/utils"
)



func initDBAndRouter(t *testing.T) *router.MuxWrapper{
	db, err := utils.OpenPostgreDatabase("postgres", "postgres", "leagueauction")
	if err != nil{
		t.Fatal(err)
	}
	//err = clearUserTable(t, db)
	if (err != nil){
		t.Fatal(err)
	}

	var r *router.MuxWrapper = new(router.MuxWrapper)
	routerCfg := router.Config{
		HostAddress: "localhost", 
		PortNo : 8080,
		Secure : false,
	}
	// routerCfg := router.Config{
	// 	HostAddress: "localhost", 
	// 	PortNo : 8081, 
	// 	Secure: true,
	// 	CertFilePath : "../../certs/cert.pem",
	// 	KeyPath : "../../certs/key.pem",
	// }
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

func executeRequest(r *router.MuxWrapper, req *http.Request) *httptest.ResponseRecorder {
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    return rr
}

func leagueauctionwsconnect(w http.ResponseWriter, r *http.Request) {

    user, ok := r.URL.Query()["user"]
    if !ok || len(user[0]) < 1 {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("invalid user"))
        return
    }

    log.Println("opening websocket for user " + user[0])
    c, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("websocket upgrade error"))
        return
    }
    // w.WriteHeader(http.StatusOK) // response.WriteHeader on hijacked connection
    defer c.Close()
}

func DialWebSocket(handlerFunc func(w http.ResponseWriter, r *http.Request), uuid string) (*websocket.Conn, *http.Response, error){

        // Create test server with the echo handler.
    s := httptest.NewServer(http.HandlerFunc(handlerFunc))
    defer s.Close()

    // Convert http://127.0.0.1 to ws://127.0.0.1
    u := "ws://localhost:8080/connect"

    // Connect to the server
    return websocket.DefaultDialer.Dial(u, nil)
}

func TestLeagueAuctionWsConnect(t *testing.T) {

    ws, _, err := DialWebSocket(leagueauctionwsconnect, "uuid1")
    if err != nil {
        t.Fatalf("server :: %v", err)
    }
    defer ws.Close()

}
func TestRegisterActivationLoginExistingUser(t *testing.T){

	r := initDBAndRouter(t)
	
	// var RegistrationJSONReqStr = []byte(`{"user_id":"x@x.com", "user_password": "pwd123"}`)
	// req, _ := http.NewRequest("POST", "/user/register", bytes.NewBuffer(RegistrationJSONReqStr))
	// response := executeRequest(r, req)
	// if response.Code != http.StatusOK{
	// 	t.Fatal("Actual code ", response.Code)
	// }
	// var m map[string]string
	// json.Unmarshal(response.Body.Bytes(), &m)
	// if m["status"] != "awaiting verification"{
	// 	t.Fatal("Actual status ", m["status"])
	// }

	// var activationJSONRequest = []byte(`{"user_id":"x@x.com", "user_activation_code": "123456"}`)
	// req, _ = http.NewRequest("POST", "/user/activation", bytes.NewBuffer(activationJSONRequest))
	// response = executeRequest(r, req)
	// if response.Code != http.StatusOK{
	// 	var m map[string]string
	// 	json.Unmarshal(response.Body.Bytes(), &m)
	// 	t.Error("Actual code ", response.Code)
	// 	t.Fatal("Actual error ", m["error"])
	// }
	// var loginResponse usermgmt.LogInResponse
	// json.Unmarshal(response.Body.Bytes(), &loginResponse)
	// if loginResponse.Token == "" || loginResponse.Expiry == ""{
	// 	t.Fatal("Expected valid login response. Actual: ", loginResponse)
	// }

	//assume x@x.com/pwd123 is already a valid account
	var loginJSONRequest = []byte(`{"user_id":"x@x.com", "user_password": "pwd123"}`)
	req, _ := http.NewRequest("POST", "/user/login", bytes.NewBuffer(loginJSONRequest))
	response := executeRequest(r, req)
	if response.Code != http.StatusOK {
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

}