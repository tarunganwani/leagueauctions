package main

import (
	"log"

	"github.com/leagueauctions/server/libs/router"
	"github.com/leagueauctions/server/servicemain"
	_ "github.com/lib/pq"
)

func main() {

	routerCfg := router.Config{
		HostAddress:  "192.168.1.22",
		PortNo:       8080,
		Secure:       true,
		CertFilePath: "../../certs/cert1.cer",
		KeyPath:      "../../certs/key1.cer",
	}
	laService := new(servicemain.LeagueAuction)
	log.Println("Initializing service")
	err := laService.InitApp(routerCfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Running service on port ", routerCfg.PortNo)
	err = laService.RunLeagueAuctionServer()
	if err != nil {
		log.Fatal(err.Error())
	}
}
