package database

import (
	"errors"
	"database/sql"
)
//LeagueAuctionDatastore - database store to hold references to all databases
type LeagueAuctionDatastore struct{
	LAUserstore 	UserStore
	LAPlayerStore 	PlayerStore
	LAAuctionStore 	AuctionStore
}


//GetLeagueAuctionDatastore -- initializes and gets a new datastore
func GetLeagueAuctionDatastore(dbObject *sql.DB) (*LeagueAuctionDatastore, error){
	if dbObject == nil{
		return nil, errors.New("db object nil")
	}
	ds := new(LeagueAuctionDatastore)
	ds.LAUserstore 		= GetUserDBStore(dbObject)
	ds.LAPlayerStore 	= GetPlayerDBStore(dbObject)
	ds.LAAuctionStore 	= GetAuctionDBStore(dbObject)
	return ds, nil
}


//GetLeagueAuctionMockstore -- initializes and gets a new mock store
func GetLeagueAuctionMockstore() (*LeagueAuctionDatastore, error){
	ds := new(LeagueAuctionDatastore)
	ds.LAUserstore 		= GetMockUserStore()
	ds.LAPlayerStore 	= GetPlayerMockStore()
	ds.LAAuctionStore 	= GetAuctionMockStore()
	return ds, nil
}