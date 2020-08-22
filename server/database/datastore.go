package database

import (
	"errors"
	"database/sql"
)
//LeagueAuctionDatastore - database store to hold references to all databases
type LeagueAuctionDatastore struct{
	userstore UserStore
}

//GetUserstore - get user store
func (lads * LeagueAuctionDatastore)GetUserstore() UserStore{
	return lads.userstore
}


//GetLeagueAuctionDatastore -- initializes and gets a new datastore
func GetLeagueAuctionDatastore(dbObject *sql.DB) (*LeagueAuctionDatastore, error){
	if dbObject == nil{
		return nil, errors.New("db object nil")
	}
	ds := new(LeagueAuctionDatastore)
	ds.userstore = GetUserDBStore(dbObject)
	return ds, nil
}

