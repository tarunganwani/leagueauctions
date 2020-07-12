package main


import (
	"net/http"
	"log"
	"fmt"
	"sync"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/gorilla/mux"
)



//---------------------------------------------------------------------
//--------------------- USER GROUP DATA AND FUNCS ---------------------
//---------------------------------------------------------------------

//UserConnectionMap has id to conn map
type UserConnectionMap map [string](*websocket.Conn)

//UserGroups - maps client id to a set of groups it belongs to 
type UserGroups map[string](map[string]interface{})

//GroupMembers - maps group id to a set/list of group members
type GroupMembers map[string](map[string]interface{})


var pool UserConnectionMap
var poolMutex *sync.Mutex

// var chatReaderMutex *sync.Mutex
// var chatWriteMutex *sync.Mutex

var groupMembers GroupMembers
var userGroups UserGroups

func addUserConnection(user string, conn *websocket.Conn){
	defer poolMutex.Unlock()
	poolMutex.Lock()
	if _, ok := pool[user]; ok {
		// already opened connection. todo - handle this
	}
	pool[user] = conn
}

func initsyncvars(){
	poolMutex = new(sync.Mutex)
	// chatReaderMutex = new(sync.Mutex)
	// chatWriteMutex = new(sync.Mutex)

}

func initmockdata(){

	pool = make(map [string](*websocket.Conn))
	
	userGroups = make(map[string]map[string]interface{})
	groupMembers = make(map[string]map[string]interface{})

	//user 1 groups
	user1Groups := map[string]interface{} {
		"group1":nil,
		"group2":nil,
		"group3":nil,
	}
	userGroups["user1"] = user1Groups

	//user 2 groups
	user2Groups := map[string]interface{} {
		"group2":nil,
		"group3":nil,
	}
	userGroups["user2"] = user2Groups

	//user 3 groups
	user3Groups := map[string]interface{} {
		"group3":nil,
	}
	userGroups["user3"] = user3Groups

	//user 4 groups
	user4Groups := map[string]interface{} {
		"group1":nil,
	}
	userGroups["user4"] = user4Groups


	// add groupMembers

	group1users:= map[string]interface{} {
		"user1":nil,
		"user4":nil,
	}
	groupMembers["group1"] = group1users

	group2users:= map[string]interface{} {
		"user1":nil,
		"user2":nil,
	}
	groupMembers["group2"] = group2users

	group3users:= map[string]interface{} {
		"user1":nil,
		"user2":nil,
		"user3":nil,
	}
	groupMembers["group3"] = group3users

}




//---------------------------------------------------------------------
//--------------------- ROUTER ENDPOINTS AND WEB SOCKETS --------------
//---------------------------------------------------------------------

//Message - auction message
type usermessage struct{
	User 	string `json:"user"`
	Group 	string `json:"group"`
	Message string `json:"message"`
}

func broadcast2group(msg usermessage) error{
	// defer chatMutex.Unlock()
	// chatMutex.Lock()
	//fetch group members
	if membersMap, ok := groupMembers[msg.Group]; ok{
		log.Println("broadcasting message to group ",msg.Group)
		for member := range membersMap{
			//check if they have active connection
			if conn, online := pool[member]; online{
				if member != msg.User{
					log.Println("sending message to ",member)
					// chatWriteMutex.Lock()
					// if err := conn.WriteJSON(msg); err != nil{
					sendmsg := "[" + msg.Group + "] " + msg.Message
					if err := conn.WriteMessage(1, []byte(sendmsg)); err != nil{
						log.Println("send msg error", err)
					}
					// chatWriteMutex.Unlock()
				}
			}
		}
		log.Println("broadcasting done")
	} else {
		log.Println("Error : group", msg.Group, "not found")
	}

	return nil
}

func reader(userid string, conn *websocket.Conn) {
	for {
		// read in a message
		log.Println("User ", userid, " listening for messages..")
		// chatReaderMutex.Lock()
		messageType, p, err := conn.ReadMessage()
		// chatReaderMutex.Unlock()
		log.Println("User ", userid, " message received = ", string(p))
		if err != nil {
			log.Println("reader ",  err)
			return
		}
		if (messageType == 1){
			var msg usermessage
			err = json.Unmarshal(p, &msg)
			if err == nil {
				log.Println("User ", userid, "p " , string(p))
				log.Println("User ", userid, "msg ", msg)
				broadcast2group(msg)
				// print out that message for clarity
				// log.Println(messageType, " <- type :: message ->", string(p))
		
				// if err := conn.WriteMessage(messageType, p); err != nil {
				// 	log.Println(err)
				// 	return
				// }
			} else{
				log.Println("User ", userid, "JSON decode error: ", err)
			}
		} else{
			log.Println("User ", userid, "Invalid message type")
		}

	}
}



// define a readerwriter which will listen and send new messages 
func readerwriter(userid string, conn *websocket.Conn) {
	//for now server will only read and broadcast messages to registered groups
	reader(userid, conn)
	
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return conn, nil
}


func serveWs(userid string, w http.ResponseWriter, r *http.Request) {
	// log.Println("WebSocket Endpoint Hit")
	conn, err := upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}
	log.Println(userid, "connected..")
	addUserConnection(userid, conn)
	readerwriter(userid, conn)
}

var muxrouter *mux.Router 

func setupRoutes() {
	
	muxrouter = mux.NewRouter()
	muxrouter.HandleFunc("/user/{userid}/connect", func(w http.ResponseWriter, r *http.Request) {
		serveWs(mux.Vars(r)["userid"], w, r)
	})
}



func main() {

	log.Println("Distributed Chat App v0.01")
	log.Println("Init mock user and group data")
	initmockdata()
	initsyncvars()
	log.Println("Setting up routes")
	setupRoutes()
	log.Println("Listening at 8080")
	http.ListenAndServe(":8080", muxrouter)
}

