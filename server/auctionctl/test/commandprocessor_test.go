package test

import(
	"testing"
	"github.com/leagueauctions/server/auctionctl"
	"github.com/leagueauctions/server/database"
	pb "github.com/leagueauctions/server/auctioncmd"
	"github.com/golang/protobuf/ptypes"
)

func TestCreateAuctionCommand(t *testing.T){

	laMockStore, _ := database.GetLeagueAuctionMockstore()
	emailid := "dummy@leagueauctions.com"
	user := database.User{EmailID : emailid}
	_ = laMockStore.LAUserstore.CreateUser(&user)

	player := database.Player{PlayerName : "TG"}
	_ = laMockStore.LAPlayerStore.UpdatePlayerInfoForUser(&player, user.UserID)

	cmdProcesor, _ := auctionctl.GetCommandProcessor(laMockStore)
	auctionboardRequest := new(pb.CreateAuctionBoardRequest)
	auctionboardRequest.AuctioneerPlayerUuid = player.PlayerID.String()
	auctionboardRequest.PurseMoney = 10000000
	auctionboardRequest.ScheduleTime = ptypes.TimestampNow()
	catA := pb.PlayerCategory{CategoryName : "CategoryA", PlayerBasePrice : 1000000}
	catB := pb.PlayerCategory{CategoryName : "CategoryB", PlayerBasePrice : 500000}
	catC := pb.PlayerCategory{CategoryName : "CategoryC", PlayerBasePrice : 300000}
	auctionboardRequest.PlayerCategory = append(auctionboardRequest.PlayerCategory, &catA, &catB, &catC)


	auctionRequest := new(pb.AuctionRequest)
	auctionRequest.RequestType = pb.AuctionRequest_CREATE_AUCTION_BOARD
	auctionRequest.Request = &pb.AuctionRequest_CreateAuctionBoardRequest{ CreateAuctionBoardRequest : auctionboardRequest }
	
	/*auctionResponse*/ _, err := cmdProcesor.ProcessAuctionRequest(auctionRequest)
	if err != nil{
		t.Fatal("ProcessAuctionRequest err :: ", err)
	}
}