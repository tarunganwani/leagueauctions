package converterutil

import(
	pb "github.com/leagueauctions/server/auctioncmd" 
	"github.com/leagueauctions/server/database" 
)

//GenerateGetPlayerInfoResponse - convert player db object into proto response
func GenerateGetPlayerInfoResponse(playerDbObj *database.Player) (*pb.GetPlayerInfoResponse) {
	if playerDbObj == nil{
		return nil
	}
	playerInfoResponse := new(pb.GetPlayerInfoResponse)
	playerInfoResponse.IsPlayerActive = playerDbObj.IsActive
	playerInfoResponse.PlayerBio = playerDbObj.PlayerBio
	playerInfoResponse.PlayerName = playerDbObj.PlayerName
	playerInfoResponse.PlayerPicture = playerDbObj.PlayerPicture
	playerInfoResponse.PlayerProfileLink = playerDbObj.PlayerProfileLink
	playerInfoResponse.PlayerType = pb.PlayerType(playerDbObj.PlayerType)
	return playerInfoResponse
}


//GeneratePlayerDbObject - convert proto response into player db object
func GeneratePlayerDbObject (updatePlayerInfoCmd *pb.UpdatePlayerInfoCommand) (*database.Player){
	if updatePlayerInfoCmd == nil{
		return nil
	}
	playerDbObj := new(database.Player)
	playerDbObj.IsActive = updatePlayerInfoCmd.IsPlayerActive
	playerDbObj.PlayerBio = updatePlayerInfoCmd.PlayerBio
	playerDbObj.PlayerName = updatePlayerInfoCmd.PlayerName
	playerDbObj.PlayerPicture = updatePlayerInfoCmd.PlayerPicture
	playerDbObj.PlayerProfileLink = updatePlayerInfoCmd.PlayerProfileLink
	playerDbObj.PlayerType = int(updatePlayerInfoCmd.PlayerType)
	
	return playerDbObj
}