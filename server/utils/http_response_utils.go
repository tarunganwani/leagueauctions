package utils

import(
	"net/http"
	"encoding/json"
)

//RespondWithJSON - encode to json and send repsonse
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}

//RespondWithError - send error response
func RespondWithError(w http.ResponseWriter, code int, message string) {
    RespondWithJSON(w, code, map[string]string{"error": message})
}


//RespondWithErrorWS - send websocket error response for 
func RespondWithErrorWS(w http.ResponseWriter, message string) {
    m := map[string]string{"error": message}
    resp, _ := json.Marshal(m)
    w.Write(resp)
}