package auctionctl

import(
	"errors"
	pb "github.com/leagueauctions/server/auctioncmd" 
	"github.com/leagueauctions/server/database" 
	"github.com/leagueauctions/server/converterutil" 
	"github.com/google/uuid"
	"database/sql"
)


func processGetPlayerInfoRequest(getPlayerInfoCmd *pb.GetPlayerInfoCommand, 
				playerStore database.PlayerStore) (*pb.AuctionResponse, error) {

	userUUID, parseErr := uuid.Parse(getPlayerInfoCmd.UserUuid)
	if parseErr != nil{
		return nil, errors.New("ProcessAuctionCmd : get_player_info uuid parse error")
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
	auctionResponse.ResponseType = pb.AuctionResponse_GET_PLAYER_INFO
	auctionResponse.Response = &pb.AuctionResponse_GetPlayerInfoResponse{ 
									GetPlayerInfoResponse : playerInfoResponse,
								}
	return auctionResponse, nil
}


func processUpdatePlayerInfoRequest(updatePlayerInfoCmd *pb.UpdatePlayerInfoCommand, 
				playerStore database.PlayerStore) (*pb.AuctionResponse, error) {

	player := converterutil.GeneratePlayerDbObject(updatePlayerInfoCmd)
	userUUID, parseErr := uuid.Parse(updatePlayerInfoCmd.UserUuid)
	if parseErr != nil{
		return generateAuctionResponse("ProcessAuctionCmd : updatePlayeInfo uuid parse error"), nil
	}
	err := playerStore.UpdatePlayerInfoForUser(player, userUUID)
	if err != nil{
		return nil, errors.New("Update player info processor : player fetch error")
	}
	playerInfoResponse := pb.UpdatePlayerInfoResponse {UpdateSuccess : true}
	auctionResponse := new(pb.AuctionResponse)
	auctionResponse.ResponseType = pb.AuctionResponse_UPDATE_PLAYER_INFO
	auctionResponse.Response = &pb.AuctionResponse_UpdatePlayerInfoResponse{ 
		UpdatePlayerInfoResponse : &playerInfoResponse,
	}
	return auctionResponse, nil
}