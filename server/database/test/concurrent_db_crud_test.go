package test

import(
	"time"
	"testing"
	"github.com/leagueauctions/server/database"
	"github.com/google/uuid"
	"fmt"
	"log"
	"sync"
)

func assertNoPanic(t *testing.T, f func(*testing.T)) {
    defer func() {
        if r := recover(); r != nil {
            t.Errorf("The code panic-ed")
        }
    }()
    f(t)
}

func GetAuctionBoard(i int) *database.AuctionBoard{
	cat1 := &database.Category{CategoryName: "A$$$$", BasePrice : 100000,}
	cat2 := &database.Category{CategoryName: "B$$$$", BasePrice : 50000,}
	catList := make([]*database.Category, 0)
	catList = append(catList, cat1, cat2)
	
	auctionName := fmt.Sprintf("FCL %d Auctions$$$$", i)
	auctionBoard := database.NewAuctionBoardObject(uuid.Nil, uuid.Nil, auctionName, 
										time.Now(), uint64(i*1000), "coins", catList)

	playerEmail := fmt.Sprintf("player%d@leagueauctions.com$$$$", i)
	playerName := fmt.Sprintf("PLAYER%d$$$$", i)
	plyr, _ := createUserAndPlayer(dbStoreName, playerEmail, playerName)
	auctionBoard.AuctioneerUUID = plyr.PlayerID
	return auctionBoard
}

func DoAuctionBoardCRUD(i int, t *testing.T){
	auctionBoard := GetAuctionBoard(i)
	//Create
	err := auctionstoresMap[dbStoreName].CreateAuctionBoard(auctionBoard)
	if (err != nil) {
		t.Error("Create error:", err)
	}
	//update
	auctionBoard.AuctionName = auctionBoard.AuctionName + " - UPDATED"
	err = auctionstoresMap[dbStoreName].UpdateAuctionBoardInfo(auctionBoard)
	if (err != nil) {
		t.Error("update error:", err)
	}
	//fetch
	auctionBoardGot, err := auctionstoresMap[dbStoreName].GetAuctionBoardInfo(auctionBoard.AuctionBoardUUID)
	if (err != nil) {
		t.Error("fetch error:", err)
	}
	if (!CompareAuctionBoards(auctionBoardGot, auctionBoard)){
		t.Error("\n\nactual ", auctionBoardGot, " \n\nexpected ", auctionBoard)
	}
	//delete
	err = auctionstoresMap[dbStoreName].DeleteAuctionBoardInfo(auctionBoard.AuctionBoardUUID)
	if (err != nil) {
		t.Error("delete error:", err)
	}
}

func ConcurrentAuctionBoards(t *testing.T){

	ConcurrentAuctionBoardsRangeFn := func (lo int, hi int, wg *sync.WaitGroup) {
		
		defer wg.Done() 
		
		log.Println("Creating auction board (db store) - range ", lo, "-", hi)
		for i := lo; i <= hi; i++ {
			DoAuctionBoardCRUD(i, t)
		}
	}

	var wg sync.WaitGroup
	for gi := 0; gi < 5; gi++{
		go ConcurrentAuctionBoardsRangeFn(gi*500 + 1, (gi + 1)*500, &wg)
		wg.Add(1)
	}

	wg.Wait()
}

func TestConcurrentDBOps(t *testing.T){
	assertNoPanic(t, ConcurrentAuctionBoards)
}