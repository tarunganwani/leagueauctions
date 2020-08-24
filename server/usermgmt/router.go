package usermgmt

import(
	"net/http"
	"encoding/json"
	"errors"
	"time"
	"database/sql"
	"github.com/leagueauctions/server/libs/router"
	"github.com/leagueauctions/server/utils"
	"github.com/leagueauctions/server/database"
)

//Router - user management router object
type Router struct {
	router 		*router.MuxWrapper
	userstore 	database.UserStore
}


//Init - Init user management router
func (u *Router)Init(r *router.MuxWrapper, userstore database.UserStore) error{
	if r == nil{
		return errors.New("router wrapper object can not be nil")
	}
	if userstore == nil{
		return errors.New("userstore object can not be nil")
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
	u.userstore = userstore
	return nil
}


func (u *Router)createLoginResponse(userObj database.User) (LogInResponse, error){
	// Declare the expiration time of the token
	expirationTime := time.Now().Add(time.Duration(utils.LoginTokenExpiryTimeInMins) * time.Minute)
	
	//Create jwt token on login
	tokenString, err :=  utils.CreateJWTForUserLogin(userObj.EmailID, expirationTime)
	if err != nil {
		// utils.RespondWithError(w, http.StatusInternalServerError, "login error: " + err.Error())
		return LogInResponse{}, err
	}
	loginResponse := LogInResponse{Token:tokenString, Expiry:expirationTime.String(), UserUUID : userObj.UserID.String()}
	return loginResponse, nil
}

//RegisterUserHandler - route to register user
func (u *Router)RegisterUserHandler(w http.ResponseWriter, r* http.Request){
	
	var registerRequestObj RegisterRequest
	err  := json.NewDecoder(r.Body).Decode(&registerRequestObj)
	if err != nil || registerRequestObj.UserID == "" || registerRequestObj.Password == ""{
		utils.RespondWithError(w, http.StatusBadRequest, "invalid register json request")
		return
	}

	//check if user already exists
	dbuser, err := u.userstore.GetUserByEmailID(registerRequestObj.UserID)
	if(err != nil && err == sql.ErrNoRows){	//new user (user doesnt exist) 

		newuser := database.User{EmailID:registerRequestObj.UserID}
		//Create a random salt and password hash store it in same user database
		newuser.PasswordHash, err = utils.HashPassword(registerRequestObj.Password)
		if err != nil{
			utils.RespondWithError(w, http.StatusInternalServerError, "error processing password")
			return
		}

		//Create a random activation code and store it in database
		//TODO: can think of using an in-memory database in future for storing OTP information
		newuser.ActivationCode = utils.GenerateRandomNumber(6)	//6 digit code
		err = u.userstore.CreateUser(&newuser)
		if err != nil{
			utils.RespondWithError(w, http.StatusInternalServerError, "error creating user")
			return
		}

		//Send activation code to verify email ID
		//TODO: send email

		resp := RegisterResponse{Status : "awaiting verification"}
		utils.RespondWithJSON(w, http.StatusOK, resp)
		return
	} else if (err != nil){
		utils.RespondWithError(w, http.StatusInternalServerError, "internal error " + err.Error())
		return
	} 
	//else error is nil - user already exists
	//if is_active or is_verified is true simply ignore the request
	if(dbuser.IsActive == true){
		resp := RegisterResponse{Status : "user already registered"}
		utils.RespondWithJSON(w, http.StatusCreated, resp)
		return
	}

	//send activation code if user is not verified yet
	//TODO: send email

	resp := RegisterResponse{Status : "awaiting verification"}
	utils.RespondWithJSON(w, http.StatusOK, resp)

	return
	
}

//ActivateUserHandler - activate user
func (u *Router)ActivateUserHandler(w http.ResponseWriter,r *http.Request){

	var activationRequest UserActivationRequest
	err := json.NewDecoder(r.Body).Decode(&activationRequest)
	if err != nil || activationRequest.UserID == "" || activationRequest.ActivationCode == ""{
		utils.RespondWithError(w, http.StatusBadRequest, "invalid activation request")
		return
	}
	
	
	dbuser, err := u.userstore.GetUserByEmailID(activationRequest.UserID)
	// err = user.GetUser(u.modelDB)
	if(err != nil){	//user doesnt exist
		if err == sql.ErrNoRows{
			utils.RespondWithError(w, http.StatusNotFound, "activation error : user not found")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "activation error : " + err.Error())
		return
	}

	if isActivationValid := ValidateActivationCode(activationRequest); isActivationValid{
		
		//Set is verified and is active to true
		dbuser.IsActive = true
		err = u.userstore.UpdateUser(dbuser)
		if(err != nil) {	
			utils.RespondWithError(w, http.StatusInternalServerError, "activation error while update " + err.Error())
			return
		}
		//create a jwt and pass it
		loginResponse, err := u.createLoginResponse(*dbuser)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "login error: " + err.Error())
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, loginResponse)
		return
	}
	utils.RespondWithError(w, http.StatusUnauthorized, "invalid activation code")
	return
}


//UserLoginHandler - handle user login
func (u * Router)UserLoginHandler(w http.ResponseWriter, r *http.Request){

	var loginRequest LogInRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil || loginRequest.UserID == "" || loginRequest.Password == ""{
		utils.RespondWithError(w, http.StatusBadRequest, "invalid login request")
		return
	}

	//Fetch user from database

	dbuser, err := u.userstore.GetUserByEmailID(loginRequest.UserID)
	if err != nil {
		if err == sql.ErrNoRows{
			utils.RespondWithError(w, http.StatusNotFound, "user not found")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "unexpected error while fetching user - login")
		return
	} 

	//Check if user activation is already done
	if dbuser.IsActive == false{
		utils.RespondWithError(w, http.StatusForbidden, "user activation required")
		return
	}

	//validate Password
	passwordsMatch, err := utils.ComparePasswords(dbuser.PasswordHash, loginRequest.Password)
	if passwordsMatch == false{
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "passwords do not match")
			return
		}
		utils.RespondWithError(w, http.StatusUnauthorized, "invalid password")
		return
	}
	loginResponse, err := u.createLoginResponse(*dbuser)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "login error: " + err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, loginResponse)
	return
}

//GetUserInfoHandler - Get user info
func (u *Router)GetUserInfoHandler(w http.ResponseWriter, r* http.Request){
	
	var userInfoRequest UserInfoRequest
	err := json.NewDecoder(r.Body).Decode(&userInfoRequest)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "get user error :" + err.Error())
		return
	}

	//Validate JSON web token
	tokenValidationStatus := utils.ValidateJWTToken(userInfoRequest.TokenString)
	if tokenValidationStatus != http.StatusOK{
		utils.RespondWithError(w, tokenValidationStatus, "invalid token")
		return
	}


	dbuser, err := u.userstore.GetUserByEmailID(userInfoRequest.UserID)
	if(err != nil){
		if (err == sql.ErrNoRows){
			utils.RespondWithError(w, http.StatusNotFound, "user not found")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userInfo := UserInfoResponse{IsActive : dbuser.IsActive, UserSerialID : dbuser.UserID.String()}
	utils.RespondWithJSON(w, http.StatusOK, userInfo)
}
