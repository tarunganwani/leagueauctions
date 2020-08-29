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

const (
	dbStoreName = "db-store"
	mockStoreName = "mock-store"
)

var userstoresMap map[string]database.UserStore
var playerstoresMap map[string]database.PlayerStore

func setupUserStore(db *sql.DB) {
	userstoresMap = make(map[string]database.UserStore)
	userstoresMap[dbStoreName] = database.GetUserDBStore(db)
	userstoresMap[mockStoreName] = database.GetMockUserStore()
}

func setupProfileStore(db *sql.DB) {
	playerstoresMap = make(map[string]database.PlayerStore)
	playerstoresMap[dbStoreName] = database.GetPlayerDBStore(db)
	playerstoresMap[mockStoreName] = database.GetPlayerMockStore()
}

func cleanup(db *sql.DB) {
	db.Exec("DELETE FROM la_schema.la_user WHERE email_id like '%$$$$'")
	db.Exec("DELETE FROM la_schema.la_player WHERE player_name like '%$$$$'")
}

func setup() (*sql.DB, error){

	db, err := utils.OpenPostgreDatabase("postgres", "postgres", "leagueauction")
	if err != nil{
		return nil, err
	}
	cleanup(db)
	setupUserStore(db)
	setupProfileStore(db)
	
	return db, nil
}

func teardown(db *sql.DB){
	cleanup(db)
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
