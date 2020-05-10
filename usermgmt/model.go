package usermgmt


import (
	"database/sql"
   	"fmt"
)

//User - business object to represent database la_user object
type User struct{
	UserID			int 	
	EmailID 		string	
	PasswordHash 	string	
	PasswordSalt 	string	
	ActivationCode	int
	IsActive		bool
}


const (
	//SelectUserByEmailIDQuery - query to fetch user attributes from users
	SelectUserByEmailIDQuery = "SELECT user_id, password_hash, password_salt, is_active FROM la_schema.la_user WHERE email_id=$1"
	//CreateUserReturnIDQuery - query to insert a user record RETURNING user_id
	CreateUserReturnIDQuery = "INSERT INTO la_schema.la_user(email_id, password_hash, password_salt) values($1,$2,$3) RETURNING user_id"
	//UpdateUserQuery - given user id, update email, password hash and salt
	UpdateUserQuery = "UPDATE la_schema.la_user SET email_id = $1, password_hash = $2, password_salt = $3, is_active = $4 where user_id = $5"
	//DeleteUserQuery - delete a user, given its user id or email id
	DeleteUserQuery = "UPDATE la_schema.la_user SET is_active = FALSE WHERE user_id = $1 OR email_id = $2"
)

//GetPostgresConnectionString - gets connection string for postgres
func GetPostgresConnectionString(user, password, dbname string) string{
	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
}
//GetUser - fetch user from database
func (u *User)GetUser(db *sql.DB) error{
	return db.QueryRow(SelectUserByEmailIDQuery,u.EmailID).Scan(&u.UserID, &u.PasswordHash, &u.PasswordSalt, &u.IsActive)
}

//CreateUser - create user from database
func (u *User)CreateUser(db *sql.DB) error{
	return db.QueryRow(CreateUserReturnIDQuery,u.EmailID, u.PasswordHash, u.PasswordSalt).Scan(&u.UserID)
}

//UpdateUser - update league auction user 
func (u *User)UpdateUser(db *sql.DB) error{
	_, err := db.Exec(UpdateUserQuery, u.EmailID, u.PasswordHash, u.PasswordSalt, u.IsActive, u.UserID)
	return err
}

//DeleteUser - delete league auction user 
func (u *User)DeleteUser(db *sql.DB) error{
	_, err := db.Exec(DeleteUserQuery, u.UserID, u.EmailID)
	return err
}