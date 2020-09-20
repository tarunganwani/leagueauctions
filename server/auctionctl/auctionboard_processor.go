package auctionctl

import(
	pb "github.com/leagueauctions/server/auctioncmd" 
	"github.com/leagueauctions/server/database" 
	"github.com/leagueauctions/server/converterutil" 
	"github.com/leagueauctions/server/utils" 
	"errors" 
)


func processCreateAuctionBoardRequest(createAuctionBoardReq *pb.CreateAuctionBoardRequest, 
				auctionStore database.AuctionStore) (*pb.AuctionResponse, error) {
	auctionBoardDbObject, err := converterutil.CreateRequestToAuctionBoardDBObject(createAuctionBoardReq)
	if err != nil {
		return nil, err
	}
	err = auctionStore.CreateAuctionBoard(auctionBoardDbObject)
	if err != nil {
		return nil, err
	}
	creteAuctionBoardResponsePb := new(pb.CreateAuctionBoardResponse)
	creteAuctionBoardResponsePb.AuctionBoardUuid = auctionBoardDbObject.AuctionBoardUUID.String()
	creteAuctionBoardResponsePb.AuctionCode = auctionBoardDbObject.AuctionCode

	auctionResponse := new(pb.AuctionResponse)
	auctionResponse.ResponseType = pb.AuctionResponse_CREATE_AUCTION_BOARD
	auctionResponse.Response = &pb.AuctionResponse_CreateAuctionBoardResponse{ 
									CreateAuctionBoardResponse : creteAuctionBoardResponsePb,
								}
	return auctionResponse, nil
}


func processUpdateAuctionBoardRequest(updateAuctionBoardReq *pb.UpdateAuctionBoardRequest, 
	auctionStore database.AuctionStore) (*pb.AuctionResponse, error) {

	auctionBoardDbObject, err := converterutil.UpdateRequestToAuctionBoardDBObject(updateAuctionBoardReq)
	if err != nil {
		return nil, err
	}
	err = auctionStore.UpdateAuctionBoardInfo(auctionBoardDbObject)
	if err != nil {
		return nil, err
	}
	updateAuctionBoardResponsePb := new(pb.UpdateAuctionBoardResponse)
	updateAuctionBoardResponsePb.Success = true

	auctionResponse := new(pb.AuctionResponse)
	auctionResponse.ResponseType = pb.AuctionResponse_UPDATE_AUCTION_BOARD
	auctionResponse.Response = &pb.AuctionResponse_UpdateAuctionBoardResponse{ 
							UpdateAuctionBoardResponse : updateAuctionBoardResponsePb,
						}
	return auctionResponse, nil
}

func processDeleteAuctionBoardRequest(deleteAuctionBoardReq *pb.DeleteAuctionBoardRequest, 
	auctionStore database.AuctionStore) (*pb.AuctionResponse, error) {

	auctionBoardUUID, parseErr := utils.GetUUIDFromString(deleteAuctionBoardReq.AuctionBoardUuid)
	if parseErr != nil{
		return nil, errors.New("delete auction board info processor : delete_auction_board_info player uuid parse error")
	}
	err := auctionStore.DeleteAuctionBoardInfo(auctionBoardUUID)
	if err != nil {
		return nil, err
	}
	deleteAuctionBoardResponsePb := new(pb.DeleteAuctionBoardResponse)
	deleteAuctionBoardResponsePb.Success = true

	auctionResponse := new(pb.AuctionResponse)
	auctionResponse.ResponseType = pb.AuctionResponse_DELETE_AUCTION_BOARD
	auctionResponse.Response = &pb.AuctionResponse_DeleteAuctionBoardResponse{ 
							DeleteAuctionBoardResponse : deleteAuctionBoardResponsePb,
						}
	return auctionResponse, nil
}

func processFetchAuctionBoardRequest(fetchAuctionBoardReq *pb.FetchAuctionBoardRequest, 
	auctionStore database.AuctionStore) (*pb.AuctionResponse, error) {

	auctionBoardUUID, parseErr := utils.GetUUIDFromString(fetchAuctionBoardReq.AuctionBoardUuid)
	if parseErr != nil{
		return nil, errors.New("fetch auction board info processor : fetch_auction_board_info player uuid parse error")
	}
	auctionBoardInfoDbObj, err := auctionStore.GetAuctionBoardInfo(auctionBoardUUID)
	if err != nil {
		return nil, err
	}
	
	fetchAuctionBoardResponsePb, err := converterutil.GenerateFetchAuctionBoardResponse(auctionBoardInfoDbObj)
	if err != nil {
		return nil, err
	}

	auctionResponse := new(pb.AuctionResponse)
	auctionResponse.ResponseType = pb.AuctionResponse_FETCH_AUCTION_BOARD
	auctionResponse.Response = &pb.AuctionResponse_FetchAuctionBoardResponse{ 
							FetchAuctionBoardResponse : fetchAuctionBoardResponsePb,
						}
	return auctionResponse, nil
}