package auctionctl

import(
	"errors"
	pb "github.com/leagueauctions/server/auctioncmd" 
	"github.com/leagueauctions/server/database" 
)

//CommandProcessor -  auction command processor object
type CommandProcessor struct{
	auctionDatastore 	*database.LeagueAuctionDatastore
}

func (mp *CommandProcessor)init(laDatastore *database.LeagueAuctionDatastore) error{
	if laDatastore == nil {
		return errors.New("laDatastore can not be nil")
	}
	mp.auctionDatastore = laDatastore
	return nil
}

//ProcessAuctionCmd - function to process auction command
func (mp *CommandProcessor)ProcessAuctionCmd(auctionCmd *pb.AuctionCommand) (*pb.AuctionResponse, error){
	if mp.auctionDatastore == nil{
		return nil, errors.New("Auction data store can not be nil")
	}
	switch auctionCmd.GetCmdType(){
	case pb.AuctionCommand_GET_PLAYER_INFO:
		if getPlayerInfoCmd := auctionCmd.GetGetPlayerInfoCmd(); getPlayerInfoCmd != nil{
			return processGetPlayerInfoRequest(getPlayerInfoCmd, mp.auctionDatastore.LAPlayerStore)
		}
		//TODO: handle error - send error in auction response? or move this to processor
		return nil, errors.New("Bad GetPlayerInfo Command")
	case pb.AuctionCommand_UPDATE_PLAYER_INFO:
		if updatePlayerInfoCmd := auctionCmd.GetUpdatePlayerInfoCmd(); updatePlayerInfoCmd != nil{
			return processUpdatePlayerInfoRequest(updatePlayerInfoCmd, mp.auctionDatastore.LAPlayerStore)
		}
		//TODO: handle error - send error in auction response? or move this to processors
		return nil, errors.New("Bad UpdatePlayerInfo Command")
	default:
		return nil, errors.New("ProcessAuctionCmd : unsupported command")
	}
}