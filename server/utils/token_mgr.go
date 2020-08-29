package utils

import(
	// "log"
	"time"
	"net/http"
	"github.com/dgrijalva/jwt-go"
)


// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type claimsT struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var jwtkey = []byte("$3CRet-Fr3Ak1nG-K3Y")

//CreateJWTForUserLogin - create json web token for user login
func CreateJWTForUserLogin(emailid string, expirationTime time.Time) (string, error){

	// Create the JWT claims, which includes the username and expiry time
	claims := &claimsT{
		Username: emailid,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtkey)
	return tokenString, err
}

//ValidateJWTToken - validate jwt token
func ValidateJWTToken(tknStr string) (statuscode int){
	
	claims := &claimsT{}
	
	//init return values
	statuscode = http.StatusOK

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtkey, nil
	})
	// log.Println(tkn, err)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			statuscode = http.StatusUnauthorized
			return
		}
		statuscode = http.StatusBadRequest
		return
	}
	if !tkn.Valid {
		statuscode = http.StatusUnauthorized
		return
	}
	return
}