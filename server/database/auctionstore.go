package database

import (
	// "log"
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

//NewAuctionBoardObject - create new board obj
func NewAuctionBoardObject(auctionBoardUUID uuid.UUID, auctioneerUUID uuid.UUID,
						name string, schtime time.Time, purse uint64, ccy string,
						catList []*Category)(*AuctionBoard){
	board := new(AuctionBoard)
	board.AuctionBoardUUID = auctionBoardUUID
	board.AuctioneerUUID = auctioneerUUID
	board.AuctionName = name
	board.ScheduleTime = schtime
	board.Purse = purse
	board.PurceCcy = ccy
	board.CategorySet = catList
	return board

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
	createAuctionBoardQuery 		= `SELECT insert_auction_board_info($1, $2, $3, $4, $5, $6, $7::category_item[]) as auction_code`
	fetchAuctionBoardQuery 			= `SELECT auctioneer_id, auction_name, schedule_time, purse, purse_ccy, is_active, auction_code FROM la_schema.la_auctionboard WHERE auction_id = $1`
	fetchAuctionCategoryListQuery 	= `SELECT category_id, category_name, base_price FROM la_schema.la_category WHERE auction_id = $1`
	deleteAuctionBoardQuery			= `UPDATE la_schema.la_auctionboard SET is_active = false WHERE auction_id = $1`
	updateAuctionBoardQuery			= `UPDATE la_schema.la_auctionboard SET auction_name = $1, schedule_time = $2, purse = $3, purse_ccy = $4 WHERE auction_id = $5`
	// updateCategoryQuery				= `UPDATE la_schema.la_category SET category_name = $1, base_price = $2 WHERE category_id = $3`
	// deleteCategoryQuery				= `DELETE FROM la_schema.la_category WHERE category_id = $1`
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
	// for _, catVal := range auctionBoard.CategorySet{
	// 	log.Println("[DBG]cat val ", *catVal)
	// }
	// log.Println("[DBG]auctionBoard val ", *auctionBoard)
	err := as.db.QueryRow(createAuctionBoardQuery,
							auctionBoard.AuctionBoardUUID, auctionBoard.AuctioneerUUID,
                            auctionBoard.AuctionName,
                            auctionBoard.ScheduleTime.UTC(), auctionBoard.Purse,
							auctionBoard.PurceCcy, pq.Array(auctionBoard.CategorySet)).Scan(&auctionBoard.AuctionCode)
	if err != nil{
		return err
	}
	return nil
}

func (as *auctionStoreDbImpl)GetAuctionBoardInfo(auctionUUID uuid.UUID) (*AuctionBoard, error) {
	if as.db == nil {
		return nil, errors.New("database object can not be nil")
	}

	auctionBoard := new(AuctionBoard)
	auctionBoard.AuctionBoardUUID = auctionUUID
	err := as.db.QueryRow(fetchAuctionBoardQuery,auctionUUID).Scan(
					&auctionBoard.AuctioneerUUID, &auctionBoard.AuctionName, 
					&auctionBoard.ScheduleTime, &auctionBoard.Purse, &auctionBoard.PurceCcy,
					&auctionBoard.IsActive,  &auctionBoard.AuctionCode)
	if err != nil{
		return nil, err
	}
	catRows, err := as.db.Query(fetchAuctionCategoryListQuery,auctionUUID)
	if err != nil{
		return nil, err
	}
	defer catRows.Close()
	for catRows.Next() {
		cat := new(Category)
		if err := catRows.Scan(&cat.CategoryUUID, &cat.CategoryName, &cat.BasePrice); err != nil {
			return nil, err
		}
		auctionBoard.CategorySet = append(auctionBoard.CategorySet, cat)
	}
	return auctionBoard, nil
}

func (as *auctionStoreDbImpl)UpdateAuctionBoardInfo(auctionBoard *AuctionBoard) error{
	if as.db == nil {
		return errors.New("database object con not be nil")
	}
	if auctionBoard == nil {
		return errors.New("auctionboard object not be nil")
	}
	_, err := as.db.Exec(updateAuctionBoardQuery,
							auctionBoard.AuctionName, auctionBoard.ScheduleTime, 
							auctionBoard.Purse, auctionBoard.PurceCcy,
							auctionBoard.AuctionBoardUUID)
	return err
}

func (as *auctionStoreDbImpl)DeleteAuctionBoardInfo(auctionUUID uuid.UUID) error{
	if as.db == nil {
		return errors.New("database object con not be nil")
	}
	if auctionUUID == uuid.Nil {
		return errors.New("auctionUUID not be nil")
	}
	_, err := as.db.Exec(deleteAuctionBoardQuery, auctionUUID)
	return err
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
	auctionBoard.AuctionBoardUUID = auctionUUID
	auctionBoard.IsActive = true
	return nil
}

func (as *auctionStoreMockImpl)GetAuctionBoardInfo(auctionUUID uuid.UUID) (*AuctionBoard, error) {
	if auctionBoard, found := as.auctionMap[auctionUUID]; found == true{
		return auctionBoard, nil
	}
	return nil, sql.ErrNoRows
}

func (as *auctionStoreMockImpl)UpdateAuctionBoardInfo(updateAuctionBoard *AuctionBoard) error{
	if auctionBoard, found := as.auctionMap[updateAuctionBoard.AuctionBoardUUID]; found == true{
		boardToUpdate:= as.auctionMap[auctionBoard.AuctionBoardUUID] 
		boardToUpdate.AuctionName = updateAuctionBoard.AuctionName
		boardToUpdate.Purse = updateAuctionBoard.Purse
		boardToUpdate.PurceCcy = updateAuctionBoard.PurceCcy
		boardToUpdate.ScheduleTime = updateAuctionBoard.ScheduleTime
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