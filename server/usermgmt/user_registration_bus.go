package usermgmt

//ValidateActivationCode valicate activation code received from user
func ValidateActivationCode(ar UserActivationRequest) bool{
	//fetch user registration info from (in-memory) registration database
	//check if the request activation code matches with database activation code
	if ar.ActivationCode == "123456"{
		return true
	}
	return false
}