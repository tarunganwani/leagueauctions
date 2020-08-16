package auctionctl

import(
	"errors"
	pb "github.com/leagueauctions/server/auctioncmd" 
)

//CommandProcessor -  auction command processor object
type CommandProcessor struct{

}

//ProcessAuctionCmd - function to process auction command
func (mp *CommandProcessor)ProcessAuctionCmd(auctionCmd *pb.AuctionCommand) (*pb.AuctionResponse, error){
	switch auctionCmd.GetCmdType(){
	case pb.AuctionCommand_GET_PLAYER_INFO:
		//do stuff
		return nil, errors.New("ProcessAuctionCmd : get_player_info unimplemented")
	case pb.AuctionCommand_UPDATE_PLAYER_INFO:
		//do stuff
		return nil, errors.New("ProcessAuctionCmd : update_player_info unimplemented")
	default:
		return nil, errors.New("ProcessAuctionCmd : unsupported command")
	}
}