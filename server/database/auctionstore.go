package database

import (
	"log"
	"fmt"
	"time"
	"errors"
	"database/sql"
	"database/sql/driver"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

//AuctionBoard - business object to represent database la_auctionboard record
type AuctionBoard struct{
	AuctionBoardUUID	uuid.UUID 	
	AuctioneerUUID		uuid.UUID 
	AuctionName			string
	ScheduleTime		time.Time
	Purse				uint64
	PurceCcy			string
	IsActive			bool
	AuctionCode			uint32
	CategorySet			[]*Category	//Name to category object map
}

//Category - business object to represent database la_category record
type Category struct{
	CategoryUUID		uuid.UUID
	AuctionBoardUUID	uuid.UUID 
	CategoryName		string
	BasePrice			uint64
}

func (c Category) Value() (driver.Value, error) {
	return fmt.Sprintf("(%s,%s,%s,%d)", c.CategoryUUID, c.AuctionBoardUUID, c.CategoryName, c.BasePrice), nil
  }
  

//AuctionStore - Auction store contract
type AuctionStore interface{
	CreateAuctionBoard(auctionBoard *AuctionBoard) error
	GetAuctionBoardInfo(auctionUUID uuid.UUID) (*AuctionBoard, error)
	UpdateAuctionBoardInfo(auctionBoard *AuctionBoard) error
	DeleteAuctionBoardInfo(auctionUUID uuid.UUID) error
}

//GetAuctionDBStore - get auction database store for tests
func GetAuctionDBStore(_db *sql.DB) AuctionStore{
	auctionstore := new(auctionStoreDbImpl)
	auctionstore.db = _db
	return auctionstore
}

const (
	createAuctionBoardQuery = `SELECT insert_auction_board_info($1, $2, $3, $4, $5, $6, $7::category_item[]) as auction_code`
)

type auctionStoreDbImpl struct{
	db *sql.DB
}

//Courtesy stackoverflow:
//https://stackoverflow.com/questions/47621459/inserting-array-of-custom-types-into-postgres

func (as *auctionStoreDbImpl)CreateAuctionBoard(auctionBoard *AuctionBoard) error{
	if as.db == nil {
		return errors.New("database object can not be nil")
	}
	if auctionBoard == nil {
		return errors.New("auctionBoard object can not be nil")
	}
	auctionBoard.AuctionBoardUUID = uuid.New()
	auctionBoard.IsActive = true
	for _, cat := range auctionBoard.CategorySet{
		cat.CategoryUUID = uuid.New()
		cat.AuctionBoardUUID = auctionBoard.AuctionBoardUUID
	}
	for _, catVal := range auctionBoard.CategorySet{
		log.Println("[DBG]cat val ", *catVal)
	}
	log.Println("[DBG]auctionBoard val ", *auctionBoard)
	err := as.db.QueryRow(createAuctionBoardQuery,auctionBoard.AuctionBoardUUID, auctionBoard.AuctioneerUUID,
												auctionBoard.AuctionName,
												auctionBoard.ScheduleTime.UTC(), auctionBoard.Purse,
												auctionBoard.PurceCcy, pq.Array(auctionBoard.CategorySet)).Scan(&auctionBoard.AuctionCode)
	if err != nil{
		return err
	}
	return nil
}

func (as *auctionStoreDbImpl)GetAuctionBoardInfo(auctionUUID uuid.UUID) (*AuctionBoard, error) {
	return nil, errors.New("Unimplemented")
}

func (as *auctionStoreDbImpl)UpdateAuctionBoardInfo(auctionBoard *AuctionBoard) error{
	return errors.New("Unimplemented")
}

func (as *auctionStoreDbImpl)DeleteAuctionBoardInfo(auctionUUID uuid.UUID) error{
	return errors.New("Unimplemented")
}

//GetAuctionMockStore - get auction database store for tests
func GetAuctionMockStore() AuctionStore{
	auctionstore := new(auctionStoreMockImpl)
	auctionstore.auctionMap = make(map[uuid.UUID]*AuctionBoard)
	return auctionstore
}

type auctionStoreMockImpl struct{
	auctionMap map[uuid.UUID]*AuctionBoard	//auction-id to auction-store map
}


func (as *auctionStoreMockImpl)CreateAuctionBoard(auctionBoard *AuctionBoard) error{
	auctionUUID := uuid.New()
	as.auctionMap [auctionUUID] = auctionBoard
	auctionBoard.IsActive = true
	return nil
}

func (as *auctionStoreMockImpl)GetAuctionBoardInfo(auctionUUID uuid.UUID) (*AuctionBoard, error) {
	if auctionBoard, found := as.auctionMap[auctionUUID]; found == true{
		return auctionBoard, nil
	}
	return nil, sql.ErrNoRows
}

func (as *auctionStoreMockImpl)UpdateAuctionBoardInfo(auctionBoard *AuctionBoard) error{
	if auctionBoard, found := as.auctionMap[auctionBoard.AuctionBoardUUID]; found == true{
		as.auctionMap[auctionBoard.AuctionBoardUUID] = auctionBoard
		return nil
	}
	return sql.ErrNoRows
}

func (as *auctionStoreMockImpl)DeleteAuctionBoardInfo(auctionUUID uuid.UUID) error{
	if auctionBoard, found := as.auctionMap[auctionUUID]; found == true{
		auctionBoard.IsActive = false
		return nil
	}
	return sql.ErrNoRows
}