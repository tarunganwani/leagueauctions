package main

import (
	utils "github.com/leagueauctions/server/Go_Adhocs/smtp/utils"
)

func main() {

	resourceManager := utils.ResourceManager{}

	senderName := resourceManager.GetProperty("smtp.sender")
	password := resourceManager.GetProperty("smtp.sender.password")
	recipient := resourceManager.GetProperty("smtp.recipient")
	server := resourceManager.GetProperty("smtp.server")
	port := resourceManager.GetProperty("smtp.port")

	sender := NewSender(senderName, password, server, port)

	//The receiver needs to be in slice as the receive supports multiple receiver
	Receiver := []string{recipient}

	Subject := "Testing HTLML Email from golang"
	message := `
	<!DOCTYPE HTML PULBLIC "-//W3C//DTD HTML 4.01 Transitional//EN">
	<html>
	<head>
	<meta http-equiv="content-type" content="text/html"; charset=ISO-8859-1">
	</head>
	<body>This is the body<br>
	<div class="moz-signature"><i><br>
	<br>
	Regards<br>
	Blake<br>
	<i></div>
	</body>
	</html>
	`
	bodyMessage := sender.WriteHTMLEmail(Receiver, Subject, message)

	sender.SendMail(Receiver, Subject, bodyMessage)
}
