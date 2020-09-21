package test

import(
	"testing"
	"github.com/leagueauctions/server/auctionctl"
	"github.com/leagueauctions/server/database"
	"github.com/leagueauctions/server/utils"
	pb "github.com/leagueauctions/server/auctioncmd"
	"github.com/golang/protobuf/ptypes"
)


func TestCRUDAuctionCommands(t *testing.T){

	laMockStore, _ := database.GetLeagueAuctionMockstore()
	emailid := "dummy@leagueauctions.com"
	user := database.User{EmailID : emailid}
	_ = laMockStore.LAUserstore.CreateUser(&user)

	player := database.Player{PlayerName : "TG"}
	_ = laMockStore.LAPlayerStore.UpdatePlayerInfoForUser(&player, user.UserID)

	cmdProcesor, _ := auctionctl.GetCommandProcessor(laMockStore)

	// -------- CREATE ---------

	purse := uint64(10000000)
	ccy := "COINS"
	auctionboardName := "FCL-AUCTIONS"

	newBoardName := "NEW-NAME"
	newPurse := uint64(20000000)
	newCcy := "RUPEES"

	auctionboardRequest := new(pb.CreateAuctionBoardRequest)
	auctionboardRequest.AuctioneerPlayerUuid = player.PlayerID.String()
	auctionboardRequest.PurseMoney = purse
	auctionboardRequest.PurseCcy = ccy
	auctionboardRequest.ScheduleTime = ptypes.TimestampNow()
	auctionboardRequest.AuctionBoardName = auctionboardName
	catA := pb.PlayerCategory{CategoryName : "CategoryA", PlayerBasePrice : 1000000}
	catB := pb.PlayerCategory{CategoryName : "CategoryB", PlayerBasePrice : 500000}
	catC := pb.PlayerCategory{CategoryName : "CategoryC", PlayerBasePrice : 300000}
	auctionboardRequest.PlayerCategoryList = append(auctionboardRequest.PlayerCategoryList, &catA, &catB, &catC)


	auctionRequest := new(pb.AuctionRequest)
	auctionRequest.RequestType = pb.AuctionRequest_CREATE_AUCTION_BOARD
	auctionRequest.Request = &pb.AuctionRequest_CreateAuctionBoardRequest{ CreateAuctionBoardRequest : auctionboardRequest }
	
	auctionResponse, err := cmdProcesor.ProcessAuctionRequest(auctionRequest)
	if err != nil{
		t.Error("ProcessAuctionRequest err :: ", err)
	}

	if auctionResponse.ResponseType != pb.AuctionResponse_CREATE_AUCTION_BOARD ||
	auctionResponse.GetCreateAuctionBoardResponse() == nil {
		t.Error("CreateAucitonBoard : Bad response")
	}

	// -------- UPDATE ---------

	updateAuctionBoardRequest := new(pb.UpdateAuctionBoardRequest)
	updateAuctionBoardRequest.AuctionBoardUuid = auctionResponse.GetCreateAuctionBoardResponse().AuctionBoardUuid
	updateAuctionBoardRequest.ScheduleTime = auctionboardRequest.ScheduleTime
	//fields to update
	updateAuctionBoardRequest.AuctionBoardName = newBoardName 
	updateAuctionBoardRequest.PurseMoney = newPurse
	updateAuctionBoardRequest.PurseCcy = newCcy

	updateAuctionRequest := new(pb.AuctionRequest)
	updateAuctionRequest.RequestType = pb.AuctionRequest_UPDATE_AUCTION_BOARD
	updateAuctionRequest.Request = &pb.AuctionRequest_UpdateAuctionBoardRequest{UpdateAuctionBoardRequest : updateAuctionBoardRequest}

	updateAuctionResponse, err := cmdProcesor.ProcessAuctionRequest(updateAuctionRequest)
	if err != nil{
		t.Fatal("UpdateAuctionRequest err :: ", err)
	}

	if updateAuctionResponse.ResponseType != pb.AuctionResponse_UPDATE_AUCTION_BOARD ||
		updateAuctionResponse.GetUpdateAuctionBoardResponse() == nil  || 
		updateAuctionResponse.GetUpdateAuctionBoardResponse().Success == false{
		t.Fatal("UpdateAucitonBoard : Bad response")
	}
	

	// -------- DELETE ---------

	deleteAuctionBoardRequest := new(pb.DeleteAuctionBoardRequest)
	deleteAuctionBoardRequest.AuctionBoardUuid = auctionResponse.GetCreateAuctionBoardResponse().AuctionBoardUuid
	
	deletteAuctionRequest := new(pb.AuctionRequest)
	deletteAuctionRequest.RequestType = pb.AuctionRequest_DELETE_AUCTION_BOARD
	deletteAuctionRequest.Request = &pb.AuctionRequest_DeleteAuctionBoardRequest{DeleteAuctionBoardRequest : deleteAuctionBoardRequest}

	deleteAuctionResponse, err := cmdProcesor.ProcessAuctionRequest(deletteAuctionRequest)
	if err != nil{
		t.Fatal("DeleteAuctionRequest err :: ", err)
	}

	if deleteAuctionResponse.ResponseType != pb.AuctionResponse_DELETE_AUCTION_BOARD ||
		deleteAuctionResponse.GetDeleteAuctionBoardResponse() == nil  || 
		deleteAuctionResponse.GetDeleteAuctionBoardResponse().Success == false{
		t.Fatal("DeleteAucitonBoard : Bad response")
	}

	// -------- FETCH ---------

	fetchAuctionBoardRequest := new(pb.FetchAuctionBoardRequest)
	fetchAuctionBoardRequest.AuctionBoardUuid = auctionResponse.GetCreateAuctionBoardResponse().AuctionBoardUuid

	fetchAuctionRequest := new(pb.AuctionRequest)
	fetchAuctionRequest.RequestType = pb.AuctionRequest_FETCH_AUCTION_BOARD_INFO
	fetchAuctionRequest.Request = &pb.AuctionRequest_FetchAuctionBoardRequest{FetchAuctionBoardRequest : fetchAuctionBoardRequest}

	fetchAuctionResponse, err := cmdProcesor.ProcessAuctionRequest(fetchAuctionRequest)
	if err != nil{
		t.Fatal("FetchAuctionRequest err :: ", err)
	}

	if fetchAuctionResponse.ResponseType != pb.AuctionResponse_FETCH_AUCTION_BOARD ||
		fetchAuctionResponse.GetFetchAuctionBoardResponse() == nil {
		t.Fatal("FetchAuctionBoard : Bad response")
	}

	fetchAuctionBoardInfoPb := fetchAuctionResponse.GetFetchAuctionBoardResponse()
	if (fetchAuctionBoardInfoPb.AuctionBoardName != newBoardName){
		t.Fatal("FetchAuctionBoard : unexpected response boardname, \n actual: ", fetchAuctionBoardInfoPb.AuctionBoardName, "\n expected:", newBoardName)
	}
	if fetchAuctionBoardInfoPb.PurseMoney != newPurse {
		t.Fatal("FetchAuctionBoard : unexpected response purse money, \n actual: ", fetchAuctionBoardInfoPb.PurseMoney, "\n expected:", newPurse)
	}
	if (fetchAuctionBoardInfoPb.PurseCcy != newCcy) {
		t.Fatal("FetchAuctionBoard : unexpected response purse ccy, \n actual: ", fetchAuctionBoardInfoPb.PurseCcy, "\n expected:", newCcy)
	}
	if (fetchAuctionBoardInfoPb.AuctioneerPlayerUuid != player.PlayerID.String()) {
		t.Fatal("FetchAuctionBoard : unexpected response player id, \n actual: ", fetchAuctionBoardInfoPb.AuctioneerPlayerUuid, "\n expected:", player.PlayerID.String())
	}
	if (fetchAuctionBoardInfoPb.IsActive != false) {
		t.Fatal("FetchAuctionBoard : unexpected response active status, \n actual: ", fetchAuctionBoardInfoPb.IsActive, "\n expected:", false)
	}
	if (!utils.CompareCategoryListPb(fetchAuctionBoardInfoPb.GetPlayerCategoryList(), auctionboardRequest.PlayerCategoryList)) {
		t.Fatal("FetchAuctionBoard : unexpected response player category list, \n actual : ", fetchAuctionBoardInfoPb.GetPlayerCategoryList(), "\n expected :", auctionboardRequest.PlayerCategoryList)
	}
	if (!utils.CompareDateTimePb(fetchAuctionBoardInfoPb.ScheduleTime, auctionboardRequest.ScheduleTime))  {
			t.Fatal("FetchAuctionBoard : unexpected response schedule time, \n\n actual: ", fetchAuctionBoardInfoPb.ScheduleTime, "\n\n expected:", auctionboardRequest.ScheduleTime)
	}

}
