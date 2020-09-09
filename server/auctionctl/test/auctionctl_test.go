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

func FetchPlayerInfoByUserUUIDRequest(playeruuidStr string) *pb.AuctionRequest{
	playerinfoRequest := new(pb.FetchPlayerInfoByUserUUIDRequest)
	playerinfoRequest.UserUuid = playeruuidStr

	auctionCmdRequest := new(pb.AuctionRequest)
	auctionCmdRequest.RequestType = pb.AuctionRequest_FETCH_PLAYER_INFO_BY_USER_UUID
	auctionCmdRequest.Request = &pb.AuctionRequest_FetchPlayerInfoByUserUuidRequest{ FetchPlayerInfoByUserUuidRequest : playerinfoRequest }
	return auctionCmdRequest
}

func FetchPlayerInfoByPlayerUUIDRequest(playeruuidStr string) *pb.AuctionRequest{
	playerinfoRequest := new(pb.FetchPlayerInfoByPlayerUUIDRequest)
	playerinfoRequest.PlayerUuid = playeruuidStr

	auctionCmdRequest := new(pb.AuctionRequest)
	auctionCmdRequest.RequestType = pb.AuctionRequest_FETCH_PLAYER_INFO_BY_PLAYER_UUID
	auctionCmdRequest.Request = &pb.AuctionRequest_FetchPlayerInfoByPlayerUuidRequest{ FetchPlayerInfoByPlayerUuidRequest : playerinfoRequest }
	return auctionCmdRequest
}

func GetUpdatePlayerInfoRequest(playeruuidStr string, playername string ,playertype int) *pb.AuctionRequest{
	updatePlayerInfoRequest := new(pb.UpdatePlayerInfoRequest)
	updatePlayerInfoRequest.UserUuid = playeruuidStr
	updatePlayerInfoRequest.PlayerName = playername
	updatePlayerInfoRequest.PlayerType = pb.PlayerType(playertype)
		
	updateAuctionCmdRequest := new(pb.AuctionRequest)
	updateAuctionCmdRequest.RequestType = pb.AuctionRequest_UPDATE_PLAYER_INFO
	updateAuctionCmdRequest.Request = &pb.AuctionRequest_UpdatePlayerInfoRequest{ UpdatePlayerInfoRequest : updatePlayerInfoRequest }
	return updateAuctionCmdRequest
}

func RegisterAndLoginUser(r *router.MuxWrapper, emailid string, pwd string, activationcode string) usermgmt.LogInResponse{
	
	userid := emailid
	pwdstr := pwd


	//Register new user
	var RegistrationJSONReqStr = []byte(`{"user_id":"`+emailid+`", "user_password":"`+pwdstr+`" }`)
	req, _ := http.NewRequest("POST", "/user/register", bytes.NewBuffer(RegistrationJSONReqStr))
	response := executeRequest(r, req)

	//activate user
	var activationJSONRequest = []byte(`{"user_id":"`+emailid+`", "user_activation_code": "`+activationcode+`"}`)
	req, _ = http.NewRequest("POST", "/user/activation", bytes.NewBuffer(activationJSONRequest))
	response = executeRequest(r, req)

	//log in to fetch sesison token
	
	jsonRequest := "{\"user_id\": \"" + userid +"\", \"user_password\": \""+pwdstr+"\"}"
	log.Println("jsonRequest ", jsonRequest)
	var loginJSONRequest = []byte(jsonRequest)
	req, _ = http.NewRequest("POST", "/user/login", bytes.NewBuffer(loginJSONRequest))
	response = executeRequest(r, req)

	var loginResponse usermgmt.LogInResponse
	json.Unmarshal(response.Body.Bytes(), &loginResponse)
	return loginResponse
}


type testCommand struct{
	cmd *pb.AuctionRequest
	cmdname string
}

func GetWebSocketConnection(t *testing.T, auctionRouter *auctionctl.Router, userid string, token string) *websocket.Conn{

	// ---- set up server to serve ws request ----
	h := http.HandlerFunc(auctionRouter.UpgradeToWsHandler)
	s := httptest.NewServer(h)
	defer s.Close()

	connectWsURL := "ws://" + s.Listener.Addr().String() + "/connect?user=" + userid + ";token=" + token
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
	return conn
}
func TestWsConnectionWithSampleCommand(t *testing.T){

	// ---- Setup user data ----

	r, _, auctionRouter := initDBAndRouter(t)
	userid := "mockuser@leagueauctions.com"
	loginResponse := RegisterAndLoginUser(r, userid, "pwd123", "123456")
	if loginResponse.Token == "" || loginResponse.Expiry == ""{
		t.Fatal("CLIENT:: Expected valid login response. Actual: ", loginResponse)
	}

	conn := GetWebSocketConnection(t, auctionRouter, userid, loginResponse.Token)
	cmdList := make([]testCommand, 0)

	//create player info
	updateAuctionCmdRequest := GetUpdatePlayerInfoRequest(loginResponse.UserUUID, "TG", 2)
	updateAuctionCmdRequest2 := GetUpdatePlayerInfoRequest(loginResponse.UserUUID, "TG-updated", 2)

	cmdList = append(cmdList, testCommand{updateAuctionCmdRequest, "Create test player TG"})
	cmdList = append(cmdList, testCommand{updateAuctionCmdRequest2, "Update test player TG name"})

	//Get player info
	auctionCmdRequest := FetchPlayerInfoByUserUUIDRequest(loginResponse.UserUUID)
	cmdList = append(cmdList, testCommand{auctionCmdRequest, "Get test player info by user uuid for TG"})

	
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
	err := error(nil)
	for _, msg := range cmdList{
		log.Println("[CLIENT] :: Command name ", msg.cmdname)
		log.Println("[CLIENT] :: Sending Command ", msg.cmd)
		protomsgCmdReqBytes, _ := proto.Marshal(msg.cmd)
		err = conn.WriteMessage(websocket.BinaryMessage, protomsgCmdReqBytes)

		if err != nil{
			t.Fatal(err)
		}
	}

	//Read from connection read channels

	for i := 0 ; i < len(cmdList); i++{
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