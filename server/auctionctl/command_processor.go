package auctionctl

import(
	"errors"
	pb "github.com/leagueauctions/server/auctioncmd" 
	"github.com/leagueauctions/server/database" 
)

//CommandProcessor -  command processor object
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

//GetCommandProcessor - initialize and return auction command processor object
func GetCommandProcessor(laStore *database.LeagueAuctionDatastore) (*CommandProcessor, error){
	cmdProcessor := new(CommandProcessor)
	return cmdProcessor, cmdProcessor.init(laStore)
}

//ProcessAuctionRequest - function to process auction command
func (mp *CommandProcessor)ProcessAuctionRequest(auctionReq *pb.AuctionRequest) (*pb.AuctionResponse, error){
	if mp.auctionDatastore == nil{
		return nil, errors.New("Auction data store can not be nil")
	}
	switch auctionReq.GetRequestType(){

	case pb.AuctionRequest_FETCH_PLAYER_INFO_BY_USER_UUID:
		if fetchPlayerInfoByUserUUIDReq := auctionReq.GetFetchPlayerInfoByUserUuidRequest(); fetchPlayerInfoByUserUUIDReq != nil{
			return processFetchPlayerInfoByUserUUIDRequest(fetchPlayerInfoByUserUUIDReq, mp.auctionDatastore.LAPlayerStore)
		}
		//TODO: handle error - send error in auction response? or move this to processor
		return nil, errors.New("Bad FetchPlayerInfoByUserUUID Request")

	case pb.AuctionRequest_FETCH_PLAYER_INFO_BY_PLAYER_UUID:
		if fetchPlayerInfoByPlayerUUIDReq := auctionReq.GetFetchPlayerInfoByPlayerUuidRequest(); fetchPlayerInfoByPlayerUUIDReq != nil{
			return processFetchPlayerInfoByPlayerUUIDRequest(fetchPlayerInfoByPlayerUUIDReq, mp.auctionDatastore.LAPlayerStore)
		}
		//TODO: handle error - send error in auction response? or move this to processor
		return nil, errors.New("Bad FetchPlayerInfoByPlayerUUID Request")

	case pb.AuctionRequest_UPDATE_PLAYER_INFO:
		if updatePlayerInfoReq := auctionReq.GetUpdatePlayerInfoRequest(); updatePlayerInfoReq != nil{
			return processUpdatePlayerInfoRequest(updatePlayerInfoReq, mp.auctionDatastore.LAPlayerStore)
		}
		//TODO: handle error - send error in auction response? or move this to processors
		return nil, errors.New("Bad UpdatePlayerInfo Request")

	case pb.AuctionRequest_CREATE_AUCTION_BOARD:
		if createAuctionReq := auctionReq.GetCreateAuctionBoardRequest(); createAuctionReq != nil{
			return processCreateAuctionBoardRequest(createAuctionReq, mp.auctionDatastore.LAAuctionStore)
		}
		return nil, errors.New("Bad CreateAuctionBoard Request")

	case pb.AuctionRequest_UPDATE_AUCTION_BOARD:
		if updateAuctionReq := auctionReq.GetUpdateAuctionBoardRequest(); updateAuctionReq != nil{
			return processUpdateAuctionBoardRequest(updateAuctionReq, mp.auctionDatastore.LAAuctionStore)
		}
		return nil, errors.New("Bad UpdateAuctionBoard Request")

	case pb.AuctionRequest_DELETE_AUCTION_BOARD:
		if deleteAuctionReq := auctionReq.GetDeleteAuctionBoardRequest(); deleteAuctionReq != nil{
			return processDeleteAuctionBoardRequest(deleteAuctionReq, mp.auctionDatastore.LAAuctionStore)
		}
		return nil, errors.New("Bad DeleteAuctionBoard Request")

	case pb.AuctionRequest_FETCH_AUCTION_BOARD_INFO:
		if fetchAuctionReq := auctionReq.GetFetchAuctionBoardRequest(); fetchAuctionReq != nil{
			return processFetchAuctionBoardRequest(fetchAuctionReq, mp.auctionDatastore.LAAuctionStore)
		}
		return nil, errors.New("Bad FetchAuctionBoard Request")

	default:
		return nil, errors.New("ProcessAuctionRequest : unsupported request")
	}
}