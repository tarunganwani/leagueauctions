package test

import (
	"database/sql"
	"testing"
	"log"
	"time"
	"github.com/leagueauctions/server/database"
	"github.com/google/uuid"
	"reflect"
)


func TestParticipantPoolForNonExisitingAuctionBoard(t *testing.T) {

	for participantstoreName, particpantstore := range participantstoresMap {
		log.Println("Test non existing particpant", participantstoreName)
		auctionUUID := uuid.New() // random uuid search
		_, err := particpantstore.FetchAllParticipants(auctionUUID)
		if (err != sql.ErrNoRows) {
			t.Error(participantstoreName, "Invalid error:", err)
		}
	}
}


func TestParticipantPoolCRUD(t *testing.T) {

	for storename, particpantstore := range participantstoresMap {
		log.Println("Test CREATE particpant", storename)

		//create user and player - for database foriegn key constraint
		usr1 := database.User{EmailID:"player511@leagueauctions.com$$$$"}
		_ = userstoresMap[storename].CreateUser(&usr1)
		plyr1  := new(database.Player)
		plyr1.PlayerName = "PLAYER511$$$$"
		_ = playerstoresMap[storename].UpdatePlayerInfoForUser(plyr1, usr1.UserID)
		
		cat1 := &database.Category{CategoryName: "A$$$$", BasePrice : 100000,}
		cat2 := &database.Category{CategoryName: "B$$$$", BasePrice : 50000,}
		catList := make([]*database.Category, 0)
		catList = append(catList, cat1, cat2)

		auctionBoard := database.NewAuctionBoardObject(uuid.Nil, plyr1.PlayerID, "FCL 7 Auctions$$$$", 
											time.Now(), 1000, "coins", catList)
		_ = auctionstoresMap[storename].CreateAuctionBoard(auctionBoard)

		auctionBoardFetched, _ := auctionstoresMap[storename].GetAuctionBoardInfo(auctionBoard.AuctionBoardUUID)

		auctionUUID := auctionBoard.AuctionBoardUUID

		participant1 := database.Participant{
			AuctionBoardUUID : auctionUUID,
			ParticipantUUID : uuid.Nil,
			ParticipantRole : database.ViewerRole,
			PlayerUUID : plyr1.PlayerID,
		}
		err := particpantstore.CreateParticipant(&participant1)
		if (err != nil) {
			t.Fatal(storename, "Create error:", err)
		}

		//UPDATE
		participant1.ParticipantRole = database.PlayerRole
		participant1.CategoryUUID = auctionBoardFetched.CategorySet[0].CategoryUUID

		err = particpantstore.UpdateParticipant(&participant1)
		if (err != nil) {
			t.Fatal(storename, "Update error:", err)
		}

		participantList, err := particpantstore.FetchAllParticipants(auctionUUID)
		if (err != nil) {
			t.Fatal(storename, "Fetch error:", err)
		}
		if reflect.DeepEqual(&participant1, participantList[0]) == false{
			t.Fatalf("\nexpected = %#v \n\nactual = %#v", &participant1, participantList[0])
		}
	}
}