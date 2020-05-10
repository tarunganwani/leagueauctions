package utils

import(
	"golang.org/x/crypto/bcrypt"
)

//HashPassword - use bcrypt to salt and hash 
func HashPassword(pwd string) (string, error) {
    // MinCost is just an integer constant provided by the bcrypt
    // package along with DefaultCost & MaxCost. 
    // The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
    hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
    if err != nil {
        return "", err
	}    
	// GenerateFromPassword returns a byte slice so we need to
    // convert the bytes to a string and return it
    return string(hash), nil
}

//ComparePasswords - compare the hash with plain password
func ComparePasswords(hashedPwd string, plainPwd string) (bool, error) {    // Since we'll be getting the hashed password from the DB it
    // will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)    
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(plainPwd))
    if err != nil {
        return false, err
    }
	return true, nil
}