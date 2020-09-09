package test

import (
	"database/sql"
	"testing"
	"log"
	"time"
	"github.com/leagueauctions/server/database"
	"github.com/google/uuid"
	_ "reflect"
)


func TestNonExisitingAuctionBoard(t *testing.T) {

	for auctionstoreName, auctionstore := range auctionstoresMap {
		log.Println("Test non existing auction", auctionstoreName)
		auctionBoardUUID := uuid.New() // random uuid search
		_, err := auctionstore.GetAuctionBoardInfo(auctionBoardUUID)
		if (err != sql.ErrNoRows) {
			t.Error(auctionstoreName, "Invalid error:", err)
		}
	}
}

func createUserAndPlayer(storeName string, emailID string, playerName string) (*database.Player, error){

	usr1 := database.User{EmailID:emailID}
	err := userstoresMap[storeName].CreateUser(&usr1)
	if (err != nil) {
		return nil, err
	}

	plyr1  := new(database.Player)
	plyr1.PlayerName = playerName
	err = playerstoresMap[storeName].UpdatePlayerInfoForUser(plyr1, usr1.UserID)
	return plyr1, err
}

func TestCreateAuctionBoard(t *testing.T) {

	cat1 := &database.Category{CategoryName: "A$$$$", BasePrice : 100000,}
	cat2 := &database.Category{CategoryName: "B$$$$", BasePrice : 50000,}
	catList := make([]*database.Category, 0)
	catList = append(catList, cat1, cat2)

	auctionBoard := new(database.AuctionBoard)
	auctionBoard.AuctioneerUUID = uuid.New()
	auctionBoard.AuctionName = "FCL 7 Auctions$$$$" 
	auctionBoard.ScheduleTime = time.Now()
	auctionBoard.Purse = 10000
	auctionBoard.PurceCcy = "coins"
	auctionBoard.CategorySet = catList

	for storeName, auctionstore := range auctionstoresMap {

		plyr, _ := createUserAndPlayer(storeName, "player100@leagueauctions.com$$$$", "PLAYER100$$$$")
		auctionBoard.AuctioneerUUID = plyr.PlayerID
		err := auctionstore.CreateAuctionBoard(auctionBoard)
		if (err != nil) {
			t.Error(storeName, "Invalid error:", err)
		}
	}
}