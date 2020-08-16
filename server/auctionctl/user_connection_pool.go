package auctionctl

import (
	"sync"
	"github.com/gorilla/websocket"
)

//UserConnectionPool - maintain a connection pool
type UserConnectionPool struct {
	useruuidToWSConnection map[string]*websocket.Conn
	maplock *sync.RWMutex
}

//Init - initialze syncronized user conenction pool
func(p *UserConnectionPool)Init(){
	p.useruuidToWSConnection = make(map[string]*websocket.Conn)
	p.maplock = new(sync.RWMutex)
}

//AddToUserConnectionPool - add user uuid and corr web socket connection to user pool
func(p *UserConnectionPool)AddToUserConnectionPool(userUUID string, conn *websocket.Conn){
	p.maplock.Lock()
	p.useruuidToWSConnection[userUUID] = conn
	p.maplock.Unlock()
}

//LookupUserConnectionPool - fetch websocket connection object by acquiring mutex read lock
func(p *UserConnectionPool)LookupUserConnectionPool(userUUID string) *websocket.Conn{
	p.maplock.RLock()
	conn, found := p.useruuidToWSConnection[userUUID]
	p.maplock.RUnlock()
	if found == true{
		return conn
	}
	return nil
}