package test


import (
	"testing"
	"log"
	// "database/sql"
	// "github.com/leagueauctions/server/usermgmt"
	"github.com/leagueauctions/server/database"
)

func CheckError(t *testing.T, err error){
	if (err != nil){
		t.Fatal(err)
	}
}

func TestCreateUsers(t *testing.T) {

	for userstoreName, userstore := range userstoresMap{
		log.Println("Test create users ", userstoreName)
		usr1 := database.User{EmailID:"user1@leagueauctions.com$$$$"}
		err := userstore.CreateUser(&usr1)
		if (err != nil){
			t.Fatal("User creation failed err:", err)
		}
		gotuser, err := userstore.GetUserByEmailID("user1@leagueauctions.com$$$$")
		if (err != nil){
			t.Fatal("User fetch failed err:", err)
		}
		if *gotuser != usr1{
			t.Fatal("expected ", usr1, " actual ", gotuser)
		}
	}
}

func TestNonExistentUser(t *testing.T) {
	for userstoreName, userstore := range userstoresMap{
		log.Println("Test non existent users ", userstoreName)
		_, err := userstore.GetUserByEmailID("nonexisitinguser@leagueauctions.com$$$$")
		if (err == nil){
			t.Fatal("Expected not found error: Got nil")
		}
	}
}



func TestUpdateUser(t *testing.T) {
	for userstoreName, userstore := range userstoresMap{
		log.Println("Test update users ", userstoreName)
		usr2 := database.User{EmailID:"user2@leagueauctions.com$$$$"}
		usr2.IsActive = false
		err := userstore.CreateUser(&usr2)
		if (err != nil){
			t.Fatal("User creation failed err:", err)
		}
		usr2.IsActive = true
		err = userstore.UpdateUser(&usr2)
		if (err != nil){
			t.Fatal("User update failed err:", err)
		}
		gotuser, err := userstore.GetUserByEmailID("user2@leagueauctions.com$$$$")
		if *gotuser != usr2{
			t.Fatal("expected ", usr2, " actual ", gotuser)
		}
	}
}


func TestDeleteUser(t *testing.T) {
	for userstoreName, userstore := range userstoresMap{
		log.Println("Test delete users ", userstoreName)
		usr3 := database.User{EmailID:"user3@leagueauctions.com$$$$", IsActive : true}
		err := userstore.CreateUser(&usr3)
		
		// log.Println("userstore CreateUser ", userstore)
		err = userstore.DeleteUser(&usr3)
		if (err != nil){
			t.Fatal("User update failed err:", err)
		}
		// log.Println("userstore DeleteUser", userstore)
		gotuser, err := userstore.GetUserByEmailID("user3@leagueauctions.com$$$$")
		if gotuser.IsActive == true{
			t.Fatal("gotuser.IsActive :: expected  = false actual ", gotuser.IsActive)
		}
	}
}
