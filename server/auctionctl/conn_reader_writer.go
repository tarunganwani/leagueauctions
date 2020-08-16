package auctionctl

import (
	"log"
    "github.com/gorilla/websocket"
	pb "github.com/leagueauctions/server/auctioncmd" 
	"github.com/golang/protobuf/proto"
)

type connectionReader struct{
	cmdProcessor *CommandProcessor
}



func (cr *connectionReader)Init() {
	cr.cmdProcessor = new(CommandProcessor)
}

//TODO: add error handling logs  and tests
func (cr *connectionReader)Listen(userid string, conn *websocket.Conn) {
	log.Println(userid + " listening... ")
	for {
		// read in a message
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			//log the error
			return
		}
		if (messageType == websocket.BinaryMessage){
			auctioncmd  := &pb.AuctionCommand{}
			if err := proto.Unmarshal(msg, auctioncmd); err != nil{
				cmdResponse, cmdErr := cr.cmdProcessor.ProcessAuctionCmd(auctioncmd)
				if cmdErr != nil{
					//log error - may be an invalid command as well
					continue
				}
				responseBytes, protoErr := proto.Marshal(cmdResponse)
				if protoErr != nil{
					//log marshalling error
					continue
				}
				if err := conn.WriteMessage(websocket.BinaryMessage, responseBytes); err != nil{
					//log.Println("send msg error", err)
					continue
				}

			} else{
				//log parsing error and continue???
			}
		} else{
			// log.Println("User ", userid, "Invalid message type")
		}

	}
}