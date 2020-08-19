package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"bytes"
	"time"
	// "strings"
	// "fmt"
	"log"
	_ "github.com/lib/pq"
	"github.com/leagueauctions/server/usermgmt"
	"github.com/leagueauctions/server/auctionctl"
	"github.com/leagueauctions/server/libs/router"
	"github.com/leagueauctions/server/utils"
    "github.com/gorilla/websocket"
	pb "github.com/leagueauctions/server/auctioncmd" 
	"github.com/golang/protobuf/proto"
)



func initDBAndRouter(t *testing.T) (*router.MuxWrapper, *usermgmt.Router, *auctionctl.Router){
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

	// initialize user connection pool
	userConnPool := new(auctionctl.UserConnectionPool)
	userConnPool.Init()

	auctionCtlRouter := new(auctionctl.Router)
	err = auctionCtlRouter.Init(r, db, userConnPool)
	if (err != nil){
		t.Fatal(err)
	}
	return r, usrMgmtRouter, auctionCtlRouter
}

func executeRequest(r *router.MuxWrapper, req *http.Request) *httptest.ResponseRecorder {
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    return rr
}



func TestWsConnectionWithSampleCommand(t *testing.T){

	r, _, auctionRouter := initDBAndRouter(t)
	
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
	userid := "x@x.com"
	pwdstr := "pwd123"
	jsonRequest := "{\"user_id\": \"" + userid +"\", \"user_password\": \""+pwdstr+"\"}"
	log.Println("jsonRequest ", jsonRequest)
	var loginJSONRequest = []byte(jsonRequest)
	req, _ := http.NewRequest("POST", "/user/login", bytes.NewBuffer(loginJSONRequest))
	response := executeRequest(r, req)
	if response.Code != http.StatusOK {
		var m map[string]string
		json.Unmarshal(response.Body.Bytes(), &m)
		t.Error("CLIENT:: Actual code ", response.Code)
		t.Fatal("CLIENT:: Actual error ", m["error"])
	}
	var loginResponse usermgmt.LogInResponse
	json.Unmarshal(response.Body.Bytes(), &loginResponse)
	if loginResponse.Token == "" || loginResponse.Expiry == ""{
		t.Fatal("CLIENT:: Expected valid login response. Actual: ", loginResponse)
	}


	//set up server to serve ws request
	h := http.HandlerFunc(auctionRouter.UpgradeToWsHandler)
	s := httptest.NewServer(h)
	defer s.Close()

	connectWsURL := "ws://" + s.Listener.Addr().String() + "/connect?user=" + userid + ";token=" + loginResponse.Token
	// connectWsURL := "ws://localhost:8080/connect?user=" + userid + ";token=" + loginResponse.Token
	conn, resp, err := websocket.DefaultDialer.Dial(connectWsURL, nil)
	if err != nil{
		t.Fatal("CLIENT:: ", err.Error())
	}
	if conn == nil{
		t.Fatal("CLIENT:: Expected non nil connection object")
	}
	if resp == nil{
		t.Fatal("CLIENT:: Expected non nil response object")
	}
	if got, want := resp.StatusCode, http.StatusSwitchingProtocols; got != want {
		t.Errorf("CLIENT:: resp.StatusCode = %q, want %q", got, want)
	}

	playerinfoRequest := pb.GetPlayerInfoCommand{}
	playerinfoRequest.UserUuid = "x@x.com"	//TODO replace by user uuid

	auctionCmdRequest := pb.AuctionCommand{}
	auctionCmdRequest.CmdType = pb.AuctionCommand_GET_PLAYER_INFO
	auctionCmdRequest.Command = &pb.AuctionCommand_GetPlayerInfoCmd{ GetPlayerInfoCmd : &playerinfoRequest }

	respChan := make(chan []byte)
	errChan := make(chan string)

	defer func(){
		close(respChan)
		close(errChan)
	}()

    go func(){
        for{
            msgtype, message, err :=  conn.ReadMessage()
            if err != nil{
				errChan <- "CLIENT:: Response error" + err.Error()
                return
			}
			if msgtype != websocket.BinaryMessage{
				errChan <- "CLIENT:: Invalid response type. Expected websocket.BinaryMessage(2) Got: " + string(msgtype)
				return
			}
			respChan <- message
			errChan <- ""
			return
        }   
    }()


	auctionCmdReqBytes, _ := proto.Marshal(&auctionCmdRequest)
	err = conn.WriteMessage(websocket.BinaryMessage, auctionCmdReqBytes)

	if err != nil{
		t.Fatal(err)
	}

	select{
	case errc := <-errChan:
		t.Error(errc)
	case auctionResponseBytes := <-respChan:

		var auctionResp pb.AuctionResponse
		err = proto.Unmarshal(auctionResponseBytes, &auctionResp)
		if err != nil{
			t.Error("CLIENT:: Error" + err.Error())
		}
		log.Println("CLIENT:: Auction response ", auctionResp.String())
		if auctionResp.GetErrormsg() != ""{
			t.Error("CLIENT:: Auction response error msg ", auctionResp.GetErrormsg())
		}
	case <-time.After(2 * time.Second):
		t.Error("CLIENT:: reader Timed-out!")
	}

	// log.Println("CLIENT:: WS connection Response ", resp)
}