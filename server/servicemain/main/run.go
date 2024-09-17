package main

import (
	"log"
<<<<<<< HEAD

	"github.com/leagueauctions/server/router"
=======
	"os"

	"github.com/leagueauctions/server/libs/router"
>>>>>>> upstream/master
	"github.com/leagueauctions/server/servicemain"
	_ "github.com/lib/pq"
)

<<<<<<< HEAD
func main() {

	routerCfg := router.Config{
		HostAddress:  "192.168.1.22",
		PortNo:       8080,
		Secure:       true,
		CertFilePath: "../../certs/cert1.cer",
		KeyPath:      "../../certs/key1.cer",
	}
=======

func main() {

	routerCfg := bakeRouterCfg()
>>>>>>> upstream/master
	laService := new(servicemain.LeagueAuction)
	log.Println("Initializing service...")
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
<<<<<<< HEAD
=======


func bakeRouterCfg() router.Config{

	enableHTTPS := ""
	certsdir := ""
	hostname := "127.0.0.1"

	if enableHTTPS = os.Getenv("HTTPSSECURE"); enableHTTPS != ""{
		//run on https
		if certsdir = os.Getenv("CERT_DIR"); certsdir == ""{
			log.Fatal("CERT_DIR environment variable not set")
		}
		certfilepath := certsdir + "/cert1.cer"
		keypath := certsdir + "/key1.cer"

		return router.Config{
			// HostAddress:  "192.168.1.22",
			HostAddress:  hostname,
			PortNo:       8081,
			Secure:       true,
			CertFilePath: certfilepath,
			KeyPath:      keypath,
		}
	}

	// plain http mode 
	return router.Config{
		HostAddress:  hostname,
		PortNo:       8080,
		Secure:       false,
	}
	
}
>>>>>>> upstream/master
