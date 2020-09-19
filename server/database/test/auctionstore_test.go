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

func CompareDateTime(dt1, dt2 time.Time) bool{
	dt1Round := dt1.Round(time.Second)
	dt2Round := dt2.Round(time.Second)
	return dt1Round.Equal(dt2Round)
}

func CompareCategories(c1, c2 []*database.Category) bool{
	m1 := make(map [string]uint64)
	m2 := make(map [string]uint64)
	for _, c := range c1 {
		m1[c.CategoryName] = c.BasePrice
	}
	for _, c := range c2 {
		m2[c.CategoryName] = c.BasePrice
	}
	return reflect.DeepEqual(m1, m2)
}

func CompareAuctionBoards(ab1, ab2 *database.AuctionBoard) bool{
	if ab1 == nil && ab2 == nil{
		return true
	}
	if ab2 != nil && ab1 != nil{
		if ab1.AuctionBoardUUID == ab2.AuctionBoardUUID &&
			ab1.AuctionCode == ab2.AuctionCode &&
			ab1.AuctionName == ab2.AuctionName &&
			ab1.AuctioneerUUID == ab2.AuctioneerUUID &&
			ab1.IsActive == ab2.IsActive &&
			ab1.PurceCcy == ab2.PurceCcy &&
			ab1.Purse == ab2.Purse &&
			CompareDateTime(ab1.ScheduleTime, ab2.ScheduleTime) &&
			CompareCategories(ab1.CategorySet, ab2.CategorySet) {
				return true
			}
	}
	return false
}

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

	auctionBoard := database.NewAuctionBoardObject(uuid.Nil, uuid.Nil, "FCL 7 Auctions$$$$", 
										time.Now(), 1000, "coins", catList)

	for storeName, auctionstore := range auctionstoresMap {

		plyr, _ := createUserAndPlayer(storeName, "player100@leagueauctions.com$$$$", "PLAYER100$$$$")
		auctionBoard.AuctioneerUUID = plyr.PlayerID
		err := auctionstore.CreateAuctionBoard(auctionBoard)
		if (err != nil) {
			t.Error(storeName, "Invalid error:", err)
		}
	}
}


func TestExistingAuctionBoard(t *testing.T) {

	cat1 := &database.Category{CategoryName: "A$$$$", BasePrice : 100000,}
	cat2 := &database.Category{CategoryName: "B$$$$", BasePrice : 50000,}
	catList := make([]*database.Category, 0)
	catList = append(catList, cat1, cat2)

	auctionBoard := database.NewAuctionBoardObject(uuid.Nil, uuid.Nil, "FCL 7 Auctions$$$$", 
										time.Now().UTC(), 1000, "coins", catList)

	for storeName, auctionstore := range auctionstoresMap {

		plyr, err := createUserAndPlayer(storeName, "player101@leagueauctions.com$$$$", "PLAYER101$$$$")
		if (err != nil) {
			t.Error(storeName, "Unexpected create user player error:", err)
		}
		auctionBoard.AuctioneerUUID = plyr.PlayerID
		err = auctionstore.CreateAuctionBoard(auctionBoard)
		if (err != nil) {
			t.Error(storeName, "Unexpected create auction board error:", err)
		}
		auctionBoardActual, err := auctionstore.GetAuctionBoardInfo(auctionBoard.AuctionBoardUUID)
		if (err != nil) {
			t.Error(storeName, "Unexpected fetch error:", err, " auctionBoard.AuctioneerUUID ", auctionBoard.AuctioneerUUID )
		}
		if CompareAuctionBoards( auctionBoardActual, auctionBoard) == false{
			t.Fatal(storeName, "auctionBoardActual ", auctionBoardActual, "auctionBoard ", auctionBoard)
		}
	}
}


func TestUpdateAuctionBoard(t *testing.T) {

	cat1 := &database.Category{CategoryName: "AA$$$$", BasePrice : 100000,}
	cat2 := &database.Category{CategoryName: "BB$$$$", BasePrice : 50000,}
	catList := make([]*database.Category, 0)
	catList = append(catList, cat1, cat2)

	auctionBoard := database.NewAuctionBoardObject(uuid.Nil, uuid.Nil, "FCL 7 Auctions$$$$", 
										time.Now().UTC(), 1000, "coins", catList)

	for storeName, auctionstore := range auctionstoresMap {

		plyr, err := createUserAndPlayer(storeName, "player102@leagueauctions.com$$$$", "PLAYER102$$$$")
		if (err != nil) {
			t.Error(storeName, "Unexpected create user player error:", err)
		}
		auctionBoard.AuctioneerUUID = plyr.PlayerID
		err = auctionstore.CreateAuctionBoard(auctionBoard)
		if (err != nil) {
			t.Error(storeName, "Unexpected create auction board error:", err)
		}
		auctionBoard.AuctionName = "UPDATED FCL 7 auctions$$$$"
		err = auctionstore.UpdateAuctionBoardInfo(auctionBoard)
		if (err != nil) {
			t.Error(storeName, "Unexpected update error:", err, " auctionBoard ", auctionBoard )
		}
		auctionBoardActual, err := auctionstore.GetAuctionBoardInfo(auctionBoard.AuctionBoardUUID)
		if (err != nil) {
			t.Error(storeName, "Unexpected fetch error:", err, " auctionBoard.AuctioneerUUID ", auctionBoard.AuctioneerUUID )
		}
		if CompareAuctionBoards( auctionBoardActual, auctionBoard) == false{
			t.Fatal(storeName, "auctionBoardActual ", auctionBoardActual, "auctionBoard ", auctionBoard)
		}
	}
}


func TestDeleteAuctionBoard(t *testing.T) {

	cat1 := &database.Category{CategoryName: "AA$$$$", BasePrice : 100000,}
	cat2 := &database.Category{CategoryName: "BB$$$$", BasePrice : 50000,}
	catList := make([]*database.Category, 0)
	catList = append(catList, cat1, cat2)

	auctionBoard := database.NewAuctionBoardObject(uuid.Nil, uuid.Nil, "FCL 7 Auctions$$$$", 
										time.Now().UTC(), 1000, "coins", catList)

	for storeName, auctionstore := range auctionstoresMap {

		plyr, err := createUserAndPlayer(storeName, "player103@leagueauctions.com$$$$", "PLAYER103$$$$")
		if (err != nil) {
			t.Error(storeName, "Unexpected create user player error:", err)
		}
		auctionBoard.AuctioneerUUID = plyr.PlayerID
		err = auctionstore.CreateAuctionBoard(auctionBoard)
		if (err != nil) {
			t.Error(storeName, "Unexpected create auction board error:", err)
		}
		err = auctionstore.DeleteAuctionBoardInfo(auctionBoard.AuctionBoardUUID)
		if (err != nil) {
			t.Error(storeName, "Unexpected delete error:", err, " auctionBoard ", auctionBoard )
		}
		auctionBoardActual, err := auctionstore.GetAuctionBoardInfo(auctionBoard.AuctionBoardUUID)
		if (err != nil) {
			t.Error(storeName, "Unexpected fetch error:", err, " auctionBoard.AuctioneerUUID ", auctionBoard.AuctioneerUUID )
		}
		if (auctionBoardActual.AuctionBoardUUID != auctionBoard.AuctionBoardUUID) || auctionBoardActual.IsActive == true{ 
			t.Fatal(storeName, "auctionBoardActual: Expecting board active status to be false - ", auctionBoardActual, "auctionBoard ", auctionBoard)
		}
	}
}