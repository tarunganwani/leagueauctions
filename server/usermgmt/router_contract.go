package usermgmt

//LogInRequest - user log in request
type LogInRequest struct {
	UserID			string	`json:"user_id"`
	Password		string	`json:"user_password"`
}

//LogInResponse - user login(or successful registration) response
type LogInResponse struct{
	Token			string	`json:"login_token"`
	Expiry			string	`json:"token_expiry"`
}

//RegisterRequest - user registration request
type RegisterRequest struct{
	UserID			string	`json:"user_id"`
	Password		string	`json:"user_password"`
}

//RegisterResponse - user registration response
type RegisterResponse struct{
	Status			string	`json:"status"`
}

//UserActivationRequest - user activation request
type UserActivationRequest struct{
	UserID			string	`json:"user_id"`
	ActivationCode	string	`json:"user_activation_code"`
}

//UserInfoRequest - user info request
type UserInfoRequest struct {
	UserID			string 	`json:"user_id"`
	TokenString		string	`json:"token_string"`
}

//UserInfoResponse - user info response
type UserInfoResponse struct {
	UserSerialID	string 	`json:"user_serial_id"`
	IsActive		bool	`json:"is_active"`
}

//TranslateUserContractToUserModelObject - convert router obj to model obj
// func (u *UserContract)TranslateUserContractToUserModelObject() User{
// 	return User{EmailID : u.EmailID}
// }