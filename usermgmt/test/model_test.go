package test


import (
	"testing"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/leagueauctions/usermgmt"
)

func CheckError(t *testing.T, err error){
	if (err != nil){
		t.Fatal(err)
	}
}

func openPostgreDatabase(user, pwd, dbname string) (*sql.DB, error){
	connectionString := usermgmt.GetPostgresConnectionString(user, pwd, dbname)
	return sql.Open("postgres", connectionString)
}



func clearUserTable(t *testing.T, db *sql.DB) error{
	_, err := db.Exec("DELETE FROM la_schema.la_user")
	if (err != nil){
		return err
	}
    _, err = db.Exec("ALTER SEQUENCE la_schema.la_user_user_id_seq RESTART WITH 1")
	if (err != nil){
		return err
	}
	return nil
}

func createuser(db *sql.DB, email, pwdhash, pwdsalt string) (usermgmt.User, error){
	usr1 := usermgmt.User{EmailID:email, PasswordHash : pwdhash, PasswordSalt : pwdsalt}
	err := usr1.CreateUser(db)
	return usr1, err
}

func TestDBConnection(t *testing.T) {
	_, err := openPostgreDatabase("postgres", "", "leagueauction")
	CheckError(t, err)
}

func TestNonExistentUser(t *testing.T) {

	db, err := openPostgreDatabase("postgres", "postgres", "leagueauction")
	CheckError(t, err)
	err = clearUserTable(t, db)
	if (err != nil){
		t.Fatal(err)
	}	
	usr1 := usermgmt.User{EmailID:"piyush@gmail.com"}
	err = usr1.GetUser(db)
	if (err != sql.ErrNoRows){
		t.Fatal("No rows expected to be fetched. err:", err)
	}
}


func TestCreateFreshUsers(t *testing.T) {

	db, err := openPostgreDatabase("postgres", "postgres", "leagueauction")
	CheckError(t, err)
	err = clearUserTable(t, db)
	if (err != nil){
		t.Fatal(err)
	}	
	
	usr1, err := createuser(db, "a@b.com", "hash1", "salt")
	if (err != nil){
		t.Fatal(err)
	}
	if (usr1.UserID != 1){
		t.Fatal("User Id should be 1 since this is the first entry after clearing the sequence")
	}

	usr2, err := createuser(db, "c@d.com", "hash2", "salt2")
	if (err != nil){
		t.Fatal(err)
	}
	if (usr2.UserID != 2){
		t.Fatal("User Id should be 2 since this is the second entry after clearing the sequence")
	}
}


func TestUpdateUser(t *testing.T) {

	db, err := openPostgreDatabase("postgres", "postgres", "leagueauction")
	CheckError(t, err)
	err = clearUserTable(t, db)
	if (err != nil){
		t.Fatal(err)
	}	
	usr1, err := createuser(db, "a@b.com", "hash1", "salt")
	if (err != nil){
		t.Fatal(err)
	}

	//change password - update hash
	usr1.PasswordHash = "hash2"
	err = usr1.UpdateUser(db)
	if (err != nil){
		t.Fatal(err)
	}

	//fetch record by email id
	usr2 := usermgmt.User{EmailID:usr1.EmailID}
	err = usr2.GetUser(db)
	if (err != nil){
		t.Fatal(err)
	}
	if(usr2.PasswordHash != "hash2"){
		t.Fatal("Expected password hash= \"hash2\" Actual password hash", usr2.PasswordHash)
	}

}


func TestDeleteUser(t *testing.T) {

	db, err := openPostgreDatabase("postgres", "postgres", "leagueauction")
	CheckError(t, err)
	err = clearUserTable(t, db)
	if (err != nil){
		t.Fatal(err)
	}	
	usr1, err := createuser(db, "a@b.com", "hash1", "salt")
	if (err != nil){
		t.Fatal(err)
	}

	//user deactivated
	err = usr1.DeleteUser(db)
	if (err != nil){
		t.Fatal(err)
	}

	//fetch record by email id
	usr2 := usermgmt.User{EmailID:usr1.EmailID}
	err = usr2.GetUser(db)
	if (err != nil){
		t.Fatal(err)
	}
	if(usr2.IsActive != false){
		t.Fatal("Expected IsActive = false Actual IsActive", usr2.IsActive)
	}
}