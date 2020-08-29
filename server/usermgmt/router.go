package usermgmt

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
<<<<<<< HEAD

	"github.com/leagueauctions/server/router"
=======
	"database/sql"
	"github.com/leagueauctions/server/libs/router"
>>>>>>> upstream/master
	"github.com/leagueauctions/server/utils"
	"github.com/leagueauctions/server/database"
)

//Router - user management router object
type Router struct {
<<<<<<< HEAD
	router  router.Wrapper
	modelDB *sql.DB
=======
	router 		*router.MuxWrapper
	userstore 	database.UserStore
>>>>>>> upstream/master
}

// Mail -
type Mail struct {
	Server   string `json:"server"`
	Port     int    `json:"port"`
	Sender   string `json:"sender"`
	Password string `json:"password"`
}

//Init - Init user management router
<<<<<<< HEAD
func (u *Router) Init(r router.Wrapper, db *sql.DB) error {
	if r == nil {
		return errors.New("router wrapper object can not be nil")
	}
	if db == nil {
		return errors.New("database object can not be nil")
=======
func (u *Router)Init(r *router.MuxWrapper, userstore database.UserStore) error{
	if r == nil{
		return errors.New("router wrapper object can not be nil")
	}
	if userstore == nil{
		return errors.New("userstore object can not be nil")
>>>>>>> upstream/master
	}
	u.router = r
	err := u.router.HandleRoute("/user/register", "POST", u.RegisterUserHandler)
	if err != nil {
		return err
	}
	err = u.router.HandleRoute("/user/activation", "POST", u.ActivateUserHandler)
	if err != nil {
		return err
	}
	err = u.router.HandleRoute("/user/login", "POST", u.UserLoginHandler)
	if err != nil {
		return err
	}
	err = u.router.HandleRoute("/user/info", "GET", u.GetUserInfoHandler)
	if err != nil {
		return err
	}
	u.userstore = userstore
	return nil
}

<<<<<<< HEAD
func (u *Router) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (u *Router) respondWithError(w http.ResponseWriter, code int, message string) {
	u.respondWithJSON(w, code, map[string]string{"error": message})
}

func (u *Router) createLoginToken(userObj User) (LogInResponse, error) {
=======

func (u *Router)createLoginResponse(userObj database.User) (LogInResponse, error){
>>>>>>> upstream/master
	// Declare the expiration time of the token
	expirationTime := time.Now().Add(time.Duration(utils.LoginTokenExpiryTimeInMins) * time.Minute)

	//Create jwt token on login
<<<<<<< HEAD
	tokenString, err := CreateJWTForUserLogin(userObj, expirationTime)
=======
	tokenString, err :=  utils.CreateJWTForUserLogin(userObj.EmailID, expirationTime)
>>>>>>> upstream/master
	if err != nil {
		// utils.RespondWithError(w, http.StatusInternalServerError, "login error: " + err.Error())
		return LogInResponse{}, err
	}
<<<<<<< HEAD
	loginResponse := LogInResponse{Token: tokenString, Expiry: expirationTime.String()}
=======
	loginResponse := LogInResponse{Token:tokenString, Expiry:expirationTime.String(), UserUUID : userObj.UserID.String()}
>>>>>>> upstream/master
	return loginResponse, nil
}

//RegisterUserHandler - route to register user
func (u *Router) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {

	var registerRequestObj RegisterRequest
<<<<<<< HEAD
	err := json.NewDecoder(r.Body).Decode(&registerRequestObj)
	if err != nil || registerRequestObj.UserID == "" || registerRequestObj.Password == "" {
		u.respondWithError(w, http.StatusBadRequest, "invalid register json request")
=======
	err  := json.NewDecoder(r.Body).Decode(&registerRequestObj)
	if err != nil || registerRequestObj.UserID == "" || registerRequestObj.Password == ""{
		utils.RespondWithError(w, http.StatusBadRequest, "invalid register json request")
>>>>>>> upstream/master
		return
	}

	//check if user already exists
<<<<<<< HEAD
	user := User{EmailID: registerRequestObj.UserID}
	err = user.GetUser(u.modelDB)
	if err != nil && err == sql.ErrNoRows { //new user (user doesnt exist)
=======
	dbuser, err := u.userstore.GetUserByEmailID(registerRequestObj.UserID)
	if(err != nil && err == sql.ErrNoRows){	//new user (user doesnt exist) 
>>>>>>> upstream/master

		newuser := database.User{EmailID:registerRequestObj.UserID}
		//Create a random salt and password hash store it in same user database
<<<<<<< HEAD
		user.PasswordHash, err = utils.HashPassword(registerRequestObj.Password)
		if err != nil {
			u.respondWithError(w, http.StatusInternalServerError, "error processing password")
=======
		newuser.PasswordHash, err = utils.HashPassword(registerRequestObj.Password)
		if err != nil{
			utils.RespondWithError(w, http.StatusInternalServerError, "error processing password")
>>>>>>> upstream/master
			return
		}

		//Create a random activation code and store it in database
		//TODO: can think of using an in-memory database in future for storing OTP information
<<<<<<< HEAD
		user.ActivationCode = utils.GenerateRandomNumber(6) //6 digit code
		err = user.CreateUser(u.modelDB)
		if err != nil {
			u.respondWithError(w, http.StatusInternalServerError, "error creating user")
=======
		newuser.ActivationCode = utils.GenerateRandomNumber(6)	//6 digit code
		err = u.userstore.CreateUser(&newuser)
		if err != nil{
			utils.RespondWithError(w, http.StatusInternalServerError, "error creating user")
>>>>>>> upstream/master
			return
		}

		//Send activation code to verify email ID
		//TODO: send email
		sendEmailToUser(user.EmailID, user.ActivationCode)

<<<<<<< HEAD
		resp := RegisterResponse{Status: "awaiting verification"}
		u.respondWithJSON(w, http.StatusOK, resp)
		return
	} else if err != nil {
		u.respondWithError(w, http.StatusInternalServerError, "internal error "+err.Error())
=======
		resp := RegisterResponse{Status : "awaiting verification"}
		utils.RespondWithJSON(w, http.StatusOK, resp)
		return
	} else if (err != nil){
		utils.RespondWithError(w, http.StatusInternalServerError, "internal error " + err.Error())
>>>>>>> upstream/master
		return
	}
	//else error is nil - user already exists
	//if is_active or is_verified is true simply ignore the request
<<<<<<< HEAD
	if user.IsActive == true {
		resp := RegisterResponse{Status: "user already registered"}
		u.respondWithJSON(w, http.StatusCreated, resp)
=======
	if(dbuser.IsActive == true){
		resp := RegisterResponse{Status : "user already registered"}
		utils.RespondWithJSON(w, http.StatusCreated, resp)
>>>>>>> upstream/master
		return
	}

	//send activation code if user is not verified yet
	//TODO: send email

<<<<<<< HEAD
	resp := RegisterResponse{Status: "awaiting verification"}
	u.respondWithJSON(w, http.StatusOK, resp)
=======
	resp := RegisterResponse{Status : "awaiting verification"}
	utils.RespondWithJSON(w, http.StatusOK, resp)
>>>>>>> upstream/master

	return

}

//sendEmailToUser -
func sendEmailToUser(recipient string, otp int) {
	mail, err := LoadConfiguration("../resources/mail.json")
	if err != nil {
		senderName := mail.Sender
		password := mail.Password
		server := mail.Server
		port := strconv.Itoa(mail.Port)

		sender := utils.NewSender(senderName, password, server, port)

		Receiver := []string{recipient}
		Subject := "League Auctions OTP"
		message := `
		<!DOCTYPE HTML PULBLIC "-//W3C//DTD HTML 4.01 Transitional//EN">
		<html>
		<head>
		<meta http-equiv="content-type" content="text/html"; charset=ISO-8859-1">
		</head>
		<body>OTP is ` + strconv.Itoa(otp) + ` <br>
		<div class="moz-signature"><i><br>
		<br>
		Regards<br>
		League Auctions App<br>
		<i></div>
		</body>
		</html>
		`
		bodyMessage := sender.WriteHTMLEmail(Receiver, Subject, message)

		sender.SendMail(Receiver, Subject, bodyMessage)
	}
}

// LoadConfiguration - from resources
func LoadConfiguration(file string) (Mail, error) {
	var mail Mail
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&mail)

	return mail, err
}

//ActivateUserHandler - activate user
func (u *Router) ActivateUserHandler(w http.ResponseWriter, r *http.Request) {

	var activationRequest UserActivationRequest
	err := json.NewDecoder(r.Body).Decode(&activationRequest)
<<<<<<< HEAD
	if err != nil || activationRequest.UserID == "" || activationRequest.ActivationCode == "" {
		u.respondWithError(w, http.StatusBadRequest, "invalid activation request")
		return
	}

	user := User{EmailID: activationRequest.UserID}
	err = user.GetUser(u.modelDB)
	if err != nil { //user doesnt exist
		if err == sql.ErrNoRows {
			u.respondWithError(w, http.StatusNotFound, "activation error : user not found")
			return
		}
		u.respondWithError(w, http.StatusInternalServerError, "activation error : "+err.Error())
=======
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
>>>>>>> upstream/master
		return
	}

	if isActivationValid := ValidateActivationCode(activationRequest); isActivationValid {

		//Set is verified and is active to true
<<<<<<< HEAD
		user.IsActive = true
		err = user.UpdateUser(u.modelDB)
		if err != nil {
			u.respondWithError(w, http.StatusInternalServerError, "activation error while update "+err.Error())
=======
		dbuser.IsActive = true
		err = u.userstore.UpdateUser(dbuser)
		if(err != nil) {	
			utils.RespondWithError(w, http.StatusInternalServerError, "activation error while update " + err.Error())
>>>>>>> upstream/master
			return
		}
		//create a jwt and pass it
		loginResponse, err := u.createLoginResponse(*dbuser)
		if err != nil {
<<<<<<< HEAD
			u.respondWithError(w, http.StatusInternalServerError, "login error: "+err.Error())
=======
			utils.RespondWithError(w, http.StatusInternalServerError, "login error: " + err.Error())
>>>>>>> upstream/master
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, loginResponse)
		return
	}
	utils.RespondWithError(w, http.StatusUnauthorized, "invalid activation code")
	return
}

//UserLoginHandler - handle user login
func (u *Router) UserLoginHandler(w http.ResponseWriter, r *http.Request) {

	var loginRequest LogInRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
<<<<<<< HEAD
	if err != nil || loginRequest.UserID == "" || loginRequest.Password == "" {
		u.respondWithError(w, http.StatusBadRequest, "invalid login request")
=======
	if err != nil || loginRequest.UserID == "" || loginRequest.Password == ""{
		utils.RespondWithError(w, http.StatusBadRequest, "invalid login request")
>>>>>>> upstream/master
		return
	}

	//Fetch user from database
<<<<<<< HEAD
	userObj := User{EmailID: loginRequest.UserID}
	err = userObj.GetUser(u.modelDB)
	if err != nil {
		if err == sql.ErrNoRows {
			u.respondWithError(w, http.StatusNotFound, "user not found")
=======

	dbuser, err := u.userstore.GetUserByEmailID(loginRequest.UserID)
	if err != nil {
		if err == sql.ErrNoRows{
			utils.RespondWithError(w, http.StatusNotFound, "user not found")
>>>>>>> upstream/master
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "unexpected error while fetching user - login")
		return
	}

	//Check if user activation is already done
<<<<<<< HEAD
	if userObj.IsActive == false {
		u.respondWithError(w, http.StatusForbidden, "user activation required")
=======
	if dbuser.IsActive == false{
		utils.RespondWithError(w, http.StatusForbidden, "user activation required")
>>>>>>> upstream/master
		return
	}

	//validate Password
<<<<<<< HEAD
	passwordsMatch, err := utils.ComparePasswords(userObj.PasswordHash, loginRequest.Password)
	if passwordsMatch == false {
=======
	passwordsMatch, err := utils.ComparePasswords(dbuser.PasswordHash, loginRequest.Password)
	if passwordsMatch == false{
>>>>>>> upstream/master
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "passwords do not match")
			return
		}
		utils.RespondWithError(w, http.StatusUnauthorized, "invalid password")
		return
	}
	loginResponse, err := u.createLoginResponse(*dbuser)
	if err != nil {
<<<<<<< HEAD
		u.respondWithError(w, http.StatusInternalServerError, "login error: "+err.Error())
=======
		utils.RespondWithError(w, http.StatusInternalServerError, "login error: " + err.Error())
>>>>>>> upstream/master
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, loginResponse)
	return
}

//GetUserInfoHandler - Get user info
func (u *Router) GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {

	var userInfoRequest UserInfoRequest
	err := json.NewDecoder(r.Body).Decode(&userInfoRequest)
	if err != nil {
<<<<<<< HEAD
		u.respondWithError(w, http.StatusBadRequest, "get user error :"+err.Error())
=======
		utils.RespondWithError(w, http.StatusBadRequest, "get user error :" + err.Error())
>>>>>>> upstream/master
		return
	}

	//Validate JSON web token
<<<<<<< HEAD
	tokenValidationStatus := ValidateJWTToken(userInfoRequest.TokenString)
	if tokenValidationStatus != http.StatusOK {
		u.respondWithError(w, tokenValidationStatus, "invalid token")
		return
	}

	user := User{EmailID: userInfoRequest.UserID}
	err = user.GetUser(u.modelDB)
	if err != nil {
		if err == sql.ErrNoRows {
			u.respondWithError(w, http.StatusNotFound, "user not found")
=======
	tokenValidationStatus := utils.ValidateJWTToken(userInfoRequest.TokenString)
	if tokenValidationStatus != http.StatusOK{
		utils.RespondWithError(w, tokenValidationStatus, "invalid token")
		return
	}


	dbuser, err := u.userstore.GetUserByEmailID(userInfoRequest.UserID)
	if(err != nil){
		if (err == sql.ErrNoRows){
			utils.RespondWithError(w, http.StatusNotFound, "user not found")
>>>>>>> upstream/master
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

<<<<<<< HEAD
	userInfo := UserInfoResponse{IsActive: user.IsActive, UserSerialID: user.UserID}
	u.respondWithJSON(w, http.StatusOK, userInfo)
=======
	userInfo := UserInfoResponse{IsActive : dbuser.IsActive, UserSerialID : dbuser.UserID.String()}
	utils.RespondWithJSON(w, http.StatusOK, userInfo)
>>>>>>> upstream/master
}
