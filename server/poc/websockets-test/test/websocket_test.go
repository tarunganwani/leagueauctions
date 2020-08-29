package test

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "log"
    "time"
    "github.com/gorilla/websocket"
	// "github.com/leagueauctions/server/libs/router"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:1024,
    WriteBufferSize:1024,
}

func echo(w http.ResponseWriter, r *http.Request) {

    c, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    defer c.Close()
    for {
        mt, message, err := c.ReadMessage()
        if err != nil {
            break
        }
        err = c.WriteMessage(mt, message)
        if err != nil {
            break
        }
    }
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

    // Convert http://127.0.0.1 to ws://127.0.0.
    u := "ws" + strings.TrimPrefix(s.URL, "http") + "?user=" + uuid

    // Connect to the server
    return websocket.DefaultDialer.Dial(u, nil)
}

func TestLeagueAuctionWsConnect(t *testing.T) {

    ws, _, err := DialWebSocket(leagueauctionwsconnect, "uuid1")
    if err != nil {
        t.Fatalf("server :: %v", err)
    }
    // if resp.StatusCode != http.StatusOK{
    //     t.Fatal("resp.StatusCode ", resp.StatusCode)
    // }
    defer ws.Close()

    // Send message to server, read response and check to see if it's what we expect.
	// if err := ws.WriteMessage(websocket.TextMessage, []byte("hello")); err != nil {
	// 	t.Fatalf("%v", err)
	// }
	// _, p, err := ws.ReadMessage()
	// if err != nil {
	// 	t.Fatalf("%v", err)
	// }
	// if string(p) != "hello" {
	// 	t.Fatalf("bad message")
	// }
}

const (
    pingtime = 1 * time.Second
    pongwait = 2 * time.Second 
	writewait = 2 * time.Second
)

func serverPinger(clientid string, c *websocket.Conn){
    log.Println("server :: [",clientid,"] Init ping generator")
    ticker := time.NewTicker(pingtime)
	defer func() {
        log.Println("server :: [",clientid,"] Closing connection..")
        ticker.Stop()   
		c.Close()
	}()
	for {
		select {
		case <-ticker.C:
			log.Println("server :: [",clientid,"] Sending ping message..")
			c.SetWriteDeadline(time.Now().Add(writewait))
			if err := c.WriteMessage(websocket.PingMessage, nil); err != nil {
                if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                    log.Printf("error: %v", err)
                }
                log.Printf("pinger normal closure")
                return
			} else{
                log.Println("server :: [",clientid,"] PINGED - waiting for client PONG..")
            }
		}
	}
}

func pongHandler(w http.ResponseWriter, r *http.Request) {

    user, ok := r.URL.Query()["user"]
    if !ok || len(user[0]) < 1 {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("invalid user"))
        return
    }

    log.Println("server :: opening websocket for user " + user[0])
    c, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("server :: websocket upgrade error"))
        return
    }

    go serverPinger(user[0], c)
    c.SetPongHandler(func(s string) error { 
		log.Println("server :: [",user[0],"] Received pong ", s)
        // c.SetReadDeadline(time.Now().Add(pongwait)); 
        return nil 
	})

	// c.SetReadDeadline(time.Now().Add(pongwait))
    for {
		log.Println("server :: [",user[0],"] Waiting for message..")
		_, message, err := c.ReadMessage()
		log.Println("server :: [",user[0],"] Received message", string(message))
		if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
            log.Printf("pong handler normal closure")
			break
		}
		// message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
	}
    

    // w.WriteHeader(http.StatusOK) // response.WriteHeader on hijacked connection
    c.Close()
}

//Server pings and client MUST pong
func TestLeagueAuctionServerPongHandler(t *testing.T) {

    clientWsConn, _, err := DialWebSocket(pongHandler, "uuid1")
    if err != nil {
        t.Fatalf("%v", err)
    }

    log.Println("client uuid1 connected - accepting ping... ")
    clientWsConn.SetPingHandler(func(s string) error { 
        log.Println("client uuid1 Received ping ", s)
        if err := clientWsConn.WriteMessage(websocket.PongMessage, nil); err != nil {
            t.Fatal("client uuid1 ponging error")
        }
        log.Println("client uuid1 sent PONG ", s)
        return nil 
	})

    log.Println("client uuid1 waiting..")
	
    //read message from websocket and put it on channel
    c := make (chan string)
    go func(){
        for{
            _, message, err :=  clientWsConn.ReadMessage()
            if err != nil{
                    c <- "Client :: Error"
                    return
            }
            c <- string(message)
        }   
    }()

    //listen to messages on channel, close connection after 6 seconds 
    go func(){
        ticker := time.NewTicker(6 * time.Second)
        defer func() {
            log.Println("client :: uuid Closing connection..")
            ticker.Stop()   
            clientWsConn.Close()
        }()
        for {
            select {
            case message := <-c:
                log.Println("[uuid1] Got message", string(message))
            case <-ticker.C:
                return
            }
        }
    }()
    
    time.Sleep(10 * time.Second)
    
}
