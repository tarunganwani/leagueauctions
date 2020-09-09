package auctionctl

import(
	"errors"
	pb "github.com/leagueauctions/server/auctioncmd" 
	"github.com/leagueauctions/server/database" 
	"github.com/leagueauctions/server/converterutil" 
	"github.com/google/uuid"
	"database/sql"
)


func processFetchPlayerInfoByUserUUIDRequest(getPlayerInfoReq *pb.FetchPlayerInfoByUserUUIDRequest, 
				playerStore database.PlayerStore) (*pb.AuctionResponse, error) {

	userUUID, parseErr := uuid.Parse(getPlayerInfoReq.UserUuid)
	if parseErr != nil{
		return nil, errors.New("Player info processor : get_player_info uuid parse error")
	}
	player, storeErr := playerStore.GetPlayerByUserUUID(userUUID)
	if storeErr != nil{
		if storeErr == sql.ErrNoRows{
			return generateAuctionResponse("GetPlayerInfo processor : Player not found"), nil
		}
		return nil, errors.New("Player info processor : player fetch error")
	}
	if player == nil {
		return nil, errors.New("Player info processor : database player fetch error")
	}
	playerInfoResponse := converterutil.GenerateGetPlayerInfoResponse(player)
	auctionResponse := new(pb.AuctionResponse)
	auctionResponse.ResponseType = pb.AuctionResponse_FETCH_PLAYER_INFO
	auctionResponse.Response = &pb.AuctionResponse_FetchPlayerInfoResponse{ 
									FetchPlayerInfoResponse : playerInfoResponse,
								}
	return auctionResponse, nil
}

func processFetchPlayerInfoByPlayerUUIDRequest(getPlayerInfoReq *pb.FetchPlayerInfoByPlayerUUIDRequest, 
				playerStore database.PlayerStore) (*pb.AuctionResponse, error) {

	playerUUID, parseErr := uuid.Parse(getPlayerInfoReq.PlayerUuid)
	if parseErr != nil{
		return nil, errors.New("Player info processor : get_player_info player uuid parse error")
	}
	player, storeErr := playerStore.GetPlayerByPlayerUUID(playerUUID)
	if storeErr != nil{
		if storeErr == sql.ErrNoRows{
			return generateAuctionResponse("GetPlayerInfo processor : Player not found"), nil
		}
		return nil, errors.New("Player info processor : player fetch error")
	}
	if player == nil {
		return nil, errors.New("Player info processor : database player fetch error")
	}
	playerInfoResponse := converterutil.GenerateGetPlayerInfoResponse(player)
	auctionResponse := new(pb.AuctionResponse)
	auctionResponse.ResponseType = pb.AuctionResponse_FETCH_PLAYER_INFO
	auctionResponse.Response = &pb.AuctionResponse_FetchPlayerInfoResponse{ 
									FetchPlayerInfoResponse : playerInfoResponse,
								}
	return auctionResponse, nil
}

func processUpdatePlayerInfoRequest(updatePlayerInfoReq *pb.UpdatePlayerInfoRequest, 
				playerStore database.PlayerStore) (*pb.AuctionResponse, error) {

	player := converterutil.GeneratePlayerDbObject(updatePlayerInfoReq)
	userUUID, parseErr := uuid.Parse(updatePlayerInfoReq.UserUuid)
	if parseErr != nil{
		return generateAuctionResponse("Update player info processor : updatePlayeInfo uuid parse error"), nil
	}
	err := playerStore.UpdatePlayerInfoForUser(player, userUUID)
	if err != nil{
		return nil, errors.New("Update player info processor : player fetch error")
	}
	playerInfoResponse := pb.UpdatePlayerInfoResponse {UpdateSuccess : true, 
												PlayerUuid : player.PlayerID.String()}
	auctionResponse := new(pb.AuctionResponse)
	auctionResponse.ResponseType = pb.AuctionResponse_UPDATE_PLAYER_INFO
	auctionResponse.Response = &pb.AuctionResponse_UpdatePlayerInfoResponse{ 
		UpdatePlayerInfoResponse : &playerInfoResponse,
	}
	return auctionResponse, nil
}