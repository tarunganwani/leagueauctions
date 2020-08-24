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
	// "github.com/leagueauctions/server/utils"
    "github.com/gorilla/websocket"
	pb "github.com/leagueauctions/server/auctioncmd" 
	"github.com/golang/protobuf/proto"
	"github.com/leagueauctions/server/database"
)



func initDBAndRouter(t *testing.T) (*router.MuxWrapper, *usermgmt.Router, *auctionctl.Router){

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
	err := r.Init(routerCfg)
	if (err != nil){
		t.Fatal(err)
	}

	//initialize user management module
	laMockStore, _ := database.GetLeagueAuctionMockstore()
	usrMgmtRouter := new(usermgmt.Router)
	err = usrMgmtRouter.Init(r, laMockStore.LAUserstore)
	if (err != nil){
		t.Fatal(err)
	}

	// initialize user connection pool
	userConnPool := new(auctionctl.UserConnectionPool)
	userConnPool.Init()

	//initialize auction ctl router
	auctionCtlRouter := new(auctionctl.Router)
	err = auctionCtlRouter.Init(r, laMockStore, userConnPool)
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
	
	// ---- Setup user data ----

	//Register new user
	var RegistrationJSONReqStr = []byte(`{"user_id":"mockuser@leagueauctions.com", "user_password": "pwd123"}`)
	req, _ := http.NewRequest("POST", "/user/register", bytes.NewBuffer(RegistrationJSONReqStr))
	response := executeRequest(r, req)

	//activate user
	var activationJSONRequest = []byte(`{"user_id":"mockuser@leagueauctions.com", "user_activation_code": "123456"}`)
	req, _ = http.NewRequest("POST", "/user/activation", bytes.NewBuffer(activationJSONRequest))
	response = executeRequest(r, req)

	//log in to fetch sesison token
	userid := "mockuser@leagueauctions.com"
	pwdstr := "pwd123"
	jsonRequest := "{\"user_id\": \"" + userid +"\", \"user_password\": \""+pwdstr+"\"}"
	log.Println("jsonRequest ", jsonRequest)
	var loginJSONRequest = []byte(jsonRequest)
	req, _ = http.NewRequest("POST", "/user/login", bytes.NewBuffer(loginJSONRequest))
	response = executeRequest(r, req)

	var loginResponse usermgmt.LogInResponse
	json.Unmarshal(response.Body.Bytes(), &loginResponse)
	if loginResponse.Token == "" || loginResponse.Expiry == ""{
		t.Fatal("CLIENT:: Expected valid login response. Actual: ", loginResponse)
	}


	// ---- set up server to serve ws request ----
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


	//Update player info
	updatePlayerInfoRequest := pb.UpdatePlayerInfoCommand{}
	updatePlayerInfoRequest.UserUuid = loginResponse.UserUUID
	updatePlayerInfoRequest.PlayerName = "TG"
	updatePlayerInfoRequest.PlayerType = 2
		
	updateAuctionCmdRequest := pb.AuctionCommand{}
	updateAuctionCmdRequest.CmdType = pb.AuctionCommand_UPDATE_PLAYER_INFO
	updateAuctionCmdRequest.Command = &pb.AuctionCommand_UpdatePlayerInfoCmd{ UpdatePlayerInfoCmd : &updatePlayerInfoRequest }


	//Get player info
	playerinfoRequest := pb.GetPlayerInfoCommand{}
	playerinfoRequest.UserUuid = loginResponse.UserUUID	

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
			log.Println("CLIENT :: RESPONSE BYTES RECD", message)
            if err != nil{
				errChan <- "CLIENT:: Response error" + err.Error()
			}
			if msgtype != websocket.BinaryMessage{
				errChan <- "CLIENT:: Invalid response type. Expected websocket.BinaryMessage(2) Got: " + string(msgtype)
			}
			respChan <- message
        }   
    }()

	//Write request on web socket channel

	updateAuctionCmdReqBytes, _ := proto.Marshal(&updateAuctionCmdRequest)
	err = conn.WriteMessage(websocket.BinaryMessage, updateAuctionCmdReqBytes)

	if err != nil{
		t.Fatal(err)
	}
	auctionCmdReqBytes, _ := proto.Marshal(&auctionCmdRequest)
	err = conn.WriteMessage(websocket.BinaryMessage, auctionCmdReqBytes)

	if err != nil{
		t.Fatal(err)
	}

	//Read from connection read channels

	for i := 0 ; i <2; i++{
		select{
		case errc := <-errChan:
			log.Println("CLIENT:: errc ", errc)
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
	}

	// log.Println("CLIENT:: WS connection Response ", resp)
}