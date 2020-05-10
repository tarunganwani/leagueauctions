package usermgmt

import(
	"net/http"
	"encoding/json"
	"errors"
	"time"
	"database/sql"
	"github.com/leagueauctions/router"
	"github.com/leagueauctions/utils"
)

//Router - user management router object
type Router struct {
	router 		router.Wrapper
	modelDB 	*sql.DB
}


//Init - Init user management router
func (u *Router)Init(r router.Wrapper, db *sql.DB) error{
	if r == nil{
		return errors.New("router wrapper object can not be nil")
	}
	if db == nil{
		return errors.New("database object can not be nil")
	}
	u.router = r
	err := u.router.HandleRoute("/user/register", "POST", u.RegisterUserHandler)
	if err != nil{
		return err
	}
	err = u.router.HandleRoute("/user/activation", "POST", u.ActivateUserHandler)
	if err != nil{
		return err
	}
	err = u.router.HandleRoute("/user/login", "POST", u.UserLoginHandler)
	if err != nil{
		return err
	}
	err = u.router.HandleRoute("/user/info", "GET", u.GetUserInfoHandler)
	if err != nil{
		return err
	}
	u.modelDB = db
	return nil
}


func (u *Router)respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}

func (u *Router)respondWithError(w http.ResponseWriter, code int, message string) {
    u.respondWithJSON(w, code, map[string]string{"error": message})
}


func (u *Router)createLoginToken(userObj User) (LogInResponse, error){
	// Declare the expiration time of the token
	expirationTime := time.Now().Add(time.Duration(utils.LoginTokenExpiryTimeInMins) * time.Minute)
	
	//Create jwt token on login
	tokenString, err :=  CreateJWTForUserLogin(userObj, expirationTime)
	if err != nil {
		// u.respondWithError(w, http.StatusInternalServerError, "login error: " + err.Error())
		return LogInResponse{}, err
	}
	loginResponse := LogInResponse{Token:tokenString, Expiry:expirationTime.String()}
	return loginResponse, nil
}

//RegisterUserHandler - route to register user
func (u *Router)RegisterUserHandler(w http.ResponseWriter, r* http.Request){
	
	var registerRequestObj RegisterRequest
	err  := json.NewDecoder(r.Body).Decode(&registerRequestObj)
	if err != nil{
		u.respondWithError(w, http.StatusBadRequest, "invalid register json request")
		return
	}

	//check if user already exists
	user := User{EmailID : registerRequestObj.UserID}
	err = user.GetUser(u.modelDB)
	if(err != nil && err == sql.ErrNoRows){	//new user (user doesnt exist) 

		//Create a random salt and password hash store it in same user database
		user.PasswordHash, err = utils.HashPassword(registerRequestObj.Password)
		if err != nil{
			u.respondWithError(w, http.StatusInternalServerError, "error processing password")
			return
		}

		//Create a random activation code and store it in database
		//TODO: can think of using an in-memory database in future for storing OTP information
		user.ActivationCode = utils.GenerateRandomNumber(6)	//6 digit code
		err = user.CreateUser(u.modelDB)
		if err != nil{
			u.respondWithError(w, http.StatusInternalServerError, "error creating user")
			return
		}

		//Send activation code to verify email ID
		//TODO: send email

		resp := RegisterResponse{Status : "awaiting verification"}
		u.respondWithJSON(w, http.StatusOK, resp)
		return
	} else if (err != nil){
		u.respondWithError(w, http.StatusInternalServerError, "internal error " + err.Error())
		return
	} 
	//else error is nil - user already exists
	//if is_active or is_verified is true simply ignore the request
	status := ""
	if(user.IsActive == true){
		status += "user already verified and active"
	}
	if status != "" {
		resp := RegisterResponse{Status : "user already active"}
		u.respondWithJSON(w, http.StatusOK, resp)
		return
	}

	//send activation code if user is not verified yet
	//TODO: send email

	resp := RegisterResponse{Status : "awaiting verification"}
	u.respondWithJSON(w, http.StatusOK, resp)

	return
	
}

//ActivateUserHandler - activate user
func (u *Router)ActivateUserHandler(w http.ResponseWriter,r *http.Request){

	var activationRequest UserActivationRequest
	err := json.NewDecoder(r.Body).Decode(&activationRequest)
	if err != nil{
		u.respondWithError(w, http.StatusBadRequest, "invalid activation request")
		return
	}
	
	if isActivationValid := ValidateActivationCode(activationRequest); isActivationValid{
		user := User{EmailID: activationRequest.UserID}
		err = user.GetUser(u.modelDB)
		if(err != nil){	//user doesnt exist
			if err == sql.ErrNoRows{
				u.respondWithError(w, http.StatusNotFound, "activation error : user not found")
				return
			}
			u.respondWithError(w, http.StatusInternalServerError, "activation error : " + err.Error())
			return
		}
		//Set is verified and is active to true
		user.IsActive = true
		err = user.UpdateUser(u.modelDB)
		if(err != nil) {	
			u.respondWithError(w, http.StatusInternalServerError, "activation error while update " + err.Error())
			return
		}
		//create a jwt and pass it
		loginResponse, err := u.createLoginToken(user)
		if err != nil {
			u.respondWithError(w, http.StatusInternalServerError, "login error: " + err.Error())
			return
		}
		u.respondWithJSON(w, http.StatusOK, loginResponse)
		return
	}
	u.respondWithError(w, http.StatusUnauthorized, "invalid activation code")
	return
}


//UserLoginHandler - handle user login
func (u * Router)UserLoginHandler(w http.ResponseWriter, r *http.Request){

	var loginRequest LogInRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil{
		u.respondWithError(w, http.StatusBadRequest, "invalid login request")
		return
	}

	//Fetch user from database
	userObj := User{EmailID : loginRequest.UserID}
	err = userObj.GetUser(u.modelDB)
	if err != nil {
		if err == sql.ErrNoRows{
			u.respondWithError(w, http.StatusNotFound, "user not found")
			return
		}
		u.respondWithError(w, http.StatusInternalServerError, "unexpected error while fetching user - login")
		return
	} 
	//validate Password
	passwordsMatch, err := utils.ComparePasswords(userObj.PasswordHash, loginRequest.Password)
	if err != nil{
		u.respondWithError(w, http.StatusInternalServerError, "error while matching passwords - login")
		return
	}
	if passwordsMatch == false {
		u.respondWithError(w, http.StatusUnauthorized, "invalid password")
		return
	}
	loginResponse, err := u.createLoginToken(userObj)
	if err != nil {
		u.respondWithError(w, http.StatusInternalServerError, "login error: " + err.Error())
		return
	}
	u.respondWithJSON(w, http.StatusOK, loginResponse)
	return
}

//GetUserInfoHandler - Get user info
func (u *Router)GetUserInfoHandler(w http.ResponseWriter, r* http.Request){
	
	var userInfoRequest UserInfoRequest
	err := json.NewDecoder(r.Body).Decode(&userInfoRequest)
	if err != nil {
		u.respondWithError(w, http.StatusBadRequest, "get user error :" + err.Error())
		return
	}

	//Validate JSON web token
	tokenValidationStatus := ValidateJWTToken(userInfoRequest.TokenString)
	if tokenValidationStatus != http.StatusOK{
		u.respondWithError(w, tokenValidationStatus, "invalid token")
		return
	}

	user := User{EmailID:userInfoRequest.UserID}
	err = user.GetUser(u.modelDB)
	if(err != nil){
		if (err == sql.ErrNoRows){
			u.respondWithError(w, http.StatusNotFound, "user not found")
			return
		}
		u.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userInfo := UserInfoResponse{IsActive : user.IsActive, UserSerialID : user.UserID}
	u.respondWithJSON(w, http.StatusOK, userInfo)
}
