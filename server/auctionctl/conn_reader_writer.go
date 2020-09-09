package auctionctl

import (
	"log"
	"errors"
    "github.com/gorilla/websocket"
	pb "github.com/leagueauctions/server/auctioncmd" 
	"github.com/golang/protobuf/proto"
)

type connectionReader struct{
	cmdProcessor *CommandProcessor
}


func (cr *connectionReader)Init() {
}

//TODO create custom logger - add connection id/uuid in every log for unique identifier
//TODO create logging levels, need not be too verbose
func (cr *connectionReader)SendErrorMsg(errmsg string, conn *websocket.Conn) {
	
	errResponse  := &pb.AuctionResponse{Errormsg : errmsg}
	log.Println("Marshalling error response", errmsg)
	responseBytes, protoErr := proto.Marshal(errResponse)
	if protoErr != nil{
		//log marshalling error
		log.Println("Listen marshalling error:", protoErr)
		return
	}
	log.Println("Sending error response over the wire:", string(responseBytes))
	if err := conn.WriteMessage(websocket.BinaryMessage, responseBytes); err != nil{
		log.Println("Listen write response error:", err)
		return
	}
	log.Println("Error sent")
}

//TODO  have a separate error handler struct for handling connection errors(logging plus sending)
//		some refactoring for a cleaner implementation
//		maybe exit loop after few repetitions of specific errors

func (cr *connectionReader)Listen(userid string, conn *websocket.Conn) (err error) {
	log.Println(userid + " listening... ")
	if cr.cmdProcessor == nil{
		log.Println("command processor can not be nil")
		err = errors.New("command processor can not be nil")
		return 
	}
	for {
		// read in a message
		messageType, msg,readErr := conn.ReadMessage()
		if readErr != nil {
			//log the error
			err = readErr
			log.Println("Listen Read error:", err)
			return
		}
		if (messageType == websocket.BinaryMessage){
			auctionReq  := &pb.AuctionRequest{}
			if err := proto.Unmarshal(msg, auctionReq); err == nil{
				log.Println("Listen ", userid, " Process auction command..")
				cmdResponse, cmdErr := cr.cmdProcessor.ProcessAuctionRequest(auctionReq)
				if cmdErr != nil{
					log.Println("Listen command processing error:", cmdErr)
					cr.SendErrorMsg(cmdErr.Error(), conn)
					continue
				}
				log.Println("Listen ", userid, "Marshalling repsonse..")
				responseBytes, protoErr := proto.Marshal(cmdResponse)
				if protoErr != nil{
					//log marshalling error
					log.Println("Listen marshalling error:", protoErr)
					continue
				}
				log.Println("Listen ", userid, " Writing repsonse")
				if err := conn.WriteMessage(websocket.BinaryMessage, responseBytes); err != nil{
					log.Println("Listen write response error:", err)
					continue
				}
				log.Println("Listen ", userid, "Response successfully sent")
			} else{
				log.Println("Listen unmarshal error:", err)
				cr.SendErrorMsg(err.Error(), conn)
				continue
			}
		} else{
			// log.Println("User ", userid, "Invalid message type")
			log.Println("Listen error: Invalid message type : user ", userid)
			cr.SendErrorMsg(err.Error(), conn)
			continue
		}

	}
}