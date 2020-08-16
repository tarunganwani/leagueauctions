package auctionctl

import(
	"log"
	"net/http"
	"errors"
	"database/sql"
	"github.com/leagueauctions/server/libs/router"
	"github.com/leagueauctions/server/utils"
    "github.com/gorilla/websocket"
)

//Router - user management router object
type Router struct {
	router 		*router.MuxWrapper
	modelDB 	*sql.DB
	upgrader 	*websocket.Upgrader
	userConnPool *UserConnectionPool
}


//Init - Init auctions router
func (ar *Router)Init(r *router.MuxWrapper, db *sql.DB, conPool *UserConnectionPool) error{

	if r == nil{
		return errors.New("router wrapper object can not be nil")
	}
	if db == nil{
		return errors.New("database object can not be nil")
	}
	if conPool == nil{
		return errors.New("connection pool object can not be nil")
	}

	// initialize user connection pool
	ar.userConnPool = conPool

	//instantiate and  initilaize upgrader
	ar.upgrader = new(websocket.Upgrader)
	//0 sets to default number of bytes(4096 at the time of writing this code)
	ar.upgrader.ReadBufferSize = 0	
	ar.upgrader.WriteBufferSize = 0	

	ar.modelDB = db
	ar.router = r
	err := ar.router.HandleRoute("/connect", "", ar.UpgradeToWsHandler)
	if err != nil{
		return err
	}
	return nil
}

//UpgradeToWsHandler - upgrade http connection to web socket
func (ar *Router)UpgradeToWsHandler(w http.ResponseWriter, r* http.Request){

	//example query - ws://localhost:8080/connect?user=x@x.com;token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InhAeC5jb20iLCJleHAiOjE1OTc1OTczOTF9.gm7epuoT7dVK6hZFEo-TcEmkpTIZx3OMtijB_Wwy8ME
    user, ok := r.URL.Query()["user"]
    if !ok || len(user[0]) < 1 {
        // w.WriteHeader(http.StatusBadRequest)
		log.Println("connect request error : invalid user" )
        w.Write([]byte("connect invalid user"))
        return
    }
    tokenstring, ok := r.URL.Query()["token"]
    if !ok || len(tokenstring[0]) < 1 {
        // w.WriteHeader(http.StatusBadRequest)
		log.Println("connect request error : empty access token" )
        w.Write([]byte("connect empty access token"))
        return
	}
	
	log.Println("Authenticating user " + user[0])
	//Validate JSON web token
	tokenValidationStatus := utils.ValidateJWTToken(tokenstring[0])
	if tokenValidationStatus != http.StatusOK{
		log.Println("connect request error : invalid access token" )
		utils.RespondWithErrorWS(w, "invalid access token")
		return
	}

	// #checkpoint 
	log.Println("Establishing ebsocket connection for user " + user[0])
	//TODO: can use response header to neogtiate sub protocol
	ar.upgrader.CheckOrigin = func(r *http.Request) bool { return true }
    con, err := ar.upgrader.Upgrade(w, r, nil/*response header*/)
    if err != nil {
		log.Println("connect request error :" + err.Error())
		utils.RespondWithErrorWS(w, "connect request - " + err.Error())
        return
	}
	
	log.Println("Websocket connection established for user " + user[0])
	
	//add connection to user-connection pool 
	ar.userConnPool.AddToUserConnectionPool(user[0], con)
	log.Println("user " + user[0] + " added to conenction pool")

	//spawn reader and writer
	cr := connectionReader{}
	log.Println("user " + user[0] + " listening.. ")
	cr.Listen(user[0], con)	//TODO consider making user id and conenction class attribute 
}