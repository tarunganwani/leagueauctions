package database

import (
	"errors"
	"database/sql"
	"github.com/google/uuid"
)

//User - business object to represent database la_user object
type User struct{
	UserID			uuid.UUID 	
	EmailID 		string	
	PasswordHash 	string	
	PasswordSalt 	string	
	ActivationCode	int
	IsActive		bool
}

//UserStore - User db store contract
type UserStore interface{
	GetUserByEmailID(emailID string) (*User, error)
	CreateUser(u *User) error
	UpdateUser(u *User) error
	DeleteUser(u *User) error
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

//userStoreDbImpl - User db store
type userStoreDbImpl struct{
	db *sql.DB
}

//GetUserDBStore - get user database store for tests
func GetUserDBStore(_db *sql.DB) UserStore{
	usrstore := new(userStoreDbImpl)
	usrstore.db = _db
	return usrstore
}

//GetUserByEmailID - fetch user from database
func (us *userStoreDbImpl)GetUserByEmailID(emailID string) (*User, error){
	if us.db == nil {
		return nil, errors.New("database object con not be nil")
	}
	u := User{EmailID : emailID}
	err := us.db.QueryRow(SelectUserByEmailIDQuery,u.EmailID).Scan(&u.UserID, &u.PasswordHash, &u.PasswordSalt, &u.IsActive)
	return &u, err
}

//CreateUser - create user from database
func (us *userStoreDbImpl)CreateUser(u *User) error{
	if us.db == nil {
		return errors.New("database object con not be nil")
	}
	return us.db.QueryRow(CreateUserReturnIDQuery,u.EmailID, u.PasswordHash, u.PasswordSalt).Scan(&u.UserID)
}

//UpdateUser - update league auction user 
func (us *userStoreDbImpl)UpdateUser(u *User) error{
	if us.db == nil {
		return errors.New("database object con not be nil")
	}
	_, err := us.db.Exec(UpdateUserQuery, u.EmailID, u.PasswordHash, u.PasswordSalt, u.IsActive, u.UserID)
	return err
}

//DeleteUser - delete league auction user 
func (us *userStoreDbImpl)DeleteUser(u *User) error{
	if us.db == nil {
		return errors.New("database object con not be nil")
	}
	_, err := us.db.Exec(DeleteUserQuery, u.UserID, u.EmailID)
	return err
}


//GetMockUserStore - get mock user store for tests
func GetMockUserStore() UserStore{
	userstore := new(userStoreMockImpl)
	userstore.userarray = make([]User, 0, 10)
	return userstore
}




// ************************** MockStore **************************

//userStoreMockImpl - User db store
type userStoreMockImpl struct{
	userarray []User
}

//GetUserByEmailID - fetch user from database
func (us *userStoreMockImpl)GetUserByEmailID(emailID string) (*User, error){
	for _, usr := range us.userarray {
		if usr.EmailID == emailID {
			return &usr, nil
		}
	}
	return nil, sql.ErrNoRows
}

//CreateUser - create user from database
func (us *userStoreMockImpl)CreateUser(u *User) error{
	if u == nil{
		return errors.New("nil user")
	}
	_, err := us.GetUserByEmailID(u.EmailID)
	if err == nil{	//User found
		return errors.New("User already exists")
	}
	u.UserID = uuid.New()
	us.userarray = append(us.userarray, *u)
	return nil
}

//UpdateUser - update league auction user 
func (us *userStoreMockImpl)UpdateUser(u *User) error{
	if u == nil{
		return errors.New("nil user")
	}

	searchIdx := -1
	for idx, usr := range us.userarray {
		if usr.EmailID == u.EmailID {
			searchIdx = idx
			break
		}
	}
	
	if searchIdx == -1{	//User not found
		return sql.ErrNoRows	//errors.New("User doesnt exist")
	}

	//udpate 
	us.userarray[searchIdx] = *u;
	return nil
}

//DeleteUser - delete league auction user 
func (us *userStoreMockImpl)DeleteUser(u *User) error{
	if u == nil{
		return errors.New("nil user")
	}

	searchIdx := -1
	for idx, usr := range us.userarray {
		if usr.EmailID == u.EmailID {
			searchIdx = idx
			break
		}
	}
	
	if searchIdx == -1{	//User not found
		return sql.ErrNoRows //errors.New("User doesnt exist")
	}
	us.userarray[searchIdx].IsActive = false
	//_ = RemoveUserFromArray(us.userarray, searchIdx)
	return nil
}
