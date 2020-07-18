package main

import (
	"log"
	"os"

	"github.com/leagueauctions/server/libs/router"
	"github.com/leagueauctions/server/servicemain"
	_ "github.com/lib/pq"
)

func main() {

	certsdir := ""
	if certsdir = os.Getenv("CERT_DIR"); certsdir == ""{
		log.Fatal("CERT_DIR environment variable not set")
	}

	certfilepath := certsdir + "/cert1.cer"
	keypath := certsdir + "/key1.cer"

	routerCfg := router.Config{
		// HostAddress:  "192.168.1.22",
		HostAddress:  "127.0.0.1",
		PortNo:       8081,
		Secure:       true,
		CertFilePath: certfilepath,
		KeyPath:      keypath,
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
