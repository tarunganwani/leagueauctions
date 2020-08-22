package test

import (
	"log"
	"os"
	"testing"
	"github.com/leagueauctions/server/database"
	"github.com/leagueauctions/server/utils"
	_ "github.com/lib/pq"
	"database/sql"
)

var userstoresMap map[string]database.UserStore

func setup() (*sql.DB, error){

	userstoresMap = make(map[string]database.UserStore)

	userstoresMap["mock-userstore"] = database.GetMockUserStore()
	
	db, err := utils.OpenPostgreDatabase("postgres", "postgres", "leagueauction")
	if err != nil{
		return nil, err
	}
	userstoresMap["db-userstore"] = database.GetUserDBStore(db)
	return db, nil
}

func teardown(db *sql.DB){
	db.Exec("DELETE FROM la_schema.la_user WHERE email_id like '%$$$$'")
}

//TestMain - run all tests in this package
func TestMain(m *testing.M){
	log.Println("Setting up..")
	db, err := setup()
	if err != nil{
		log.Println("Setup error ", err)
		os.Exit(1)
	}
	v := m.Run()
	log.Println("Tearing down..")
	teardown(db)
	os.Exit(v)
}
