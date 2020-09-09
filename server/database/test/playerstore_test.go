package test

import (
	"database/sql"
	"testing"
	"log"
	"github.com/leagueauctions/server/database"
	"github.com/google/uuid"
	"reflect"
)

func TestNonExisitingPlayers(t *testing.T) {

	for playerstoreName, playerstore := range playerstoresMap{
		log.Println("Test non existing players", playerstoreName)
		_, err := playerstore.GetPlayerByUserUUID(uuid.New())
		if (err != sql.ErrNoRows) {
			t.Fatal("Invalid error:", err)
		}
	}
}


func TestCreateAndGetExisitingPlayers(t *testing.T) {

	for playerstoreName, playerstore := range playerstoresMap{
		log.Println("Test create and get existing players", playerstoreName)
		usr1 := database.User{EmailID:"player1@leagueauctions.com$$$$"}
		err := userstoresMap[playerstoreName].CreateUser(&usr1)
		if (err != nil) {
			t.Fatal("User creation failed err:", err)
		}
		plyr1  := new(database.Player)
		plyr1.PlayerName = "PLAYER1$$$$"
		err = playerstore.UpdatePlayerInfoForUser(plyr1, usr1.UserID)
		if (err != nil) {
			t.Fatal("Player creation failed err:", err)
		}
		plyr1.PlayerPicture = []uint8{}	// default value gotten from database
		plyrGot, err := playerstore.GetPlayerByUserUUID(usr1.UserID)
		if (err != nil) {
			t.Fatal("Player FETCH failed err:", err)
		}
		if reflect.DeepEqual(*plyrGot, *plyr1) == false{
			t.Fatalf("Player FETCH failed\nexpected  = %#v\nactual = %#v\n", *plyr1, *plyrGot)
		}

		//Test update
		plyr1.PlayerBio = "PLAYER1-BIO"
		err = playerstore.UpdatePlayerInfoForUser(plyr1, usr1.UserID)
		if (err != nil) {
			t.Fatal("Player creation failed err:", err)
		}
		plyrGot, err = playerstore.GetPlayerByUserUUID(usr1.UserID)
		if (err != nil) {
			t.Fatal("Player FETCH failed err:", err)
		}
		plyrGotByPlayerUUID, err := playerstore.GetPlayerByPlayerUUID(plyr1.PlayerID)
		if (err != nil) {
			t.Fatal("Player FETCH failed err:", err)
		}
		plyr1.PlayerPicture = []uint8{}	// default value gotten from database
		if reflect.DeepEqual(*plyrGot, *plyr1) == false{
			t.Fatalf("Player FETCH failed\nexpected  = %#v\nactual = %#v\n", *plyr1, *plyrGot)
		}
		if reflect.DeepEqual(*plyrGot, *plyrGotByPlayerUUID) == false{
			t.Fatalf("Player FETCH failed\nexpected  = %#v\nactual = %#v\n", *plyr1, *plyrGot)
		}
	}
}
