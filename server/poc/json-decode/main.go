package main

import(
	"fmt"
	"encoding/json"
)

type WSRequest struct{
	Type string `json:"type"`
}

func main(){
	requestWSString := []byte(`{"type":"bidder", "kind":"info", "bidder_id":"v152g-ahjsh222-ashj"}`)
	var wsRequest WSRequest
	err := json.Unmarshal(requestWSString, &wsRequest)
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Printf("request object %#v \n", wsRequest)
}