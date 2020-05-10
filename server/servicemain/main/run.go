package main

import(
	"github.com/leagueauctions/server/servicemain"
	"github.com/leagueauctions/server/router"
	_ "github.com/lib/pq"
	"log"
)

func main(){
	
	routerCfg := router.Config{HostAddress: "localhost", PortNo : 8081}
	laService := new(servicemain.LeagueAuction)
	log.Println("Initializing service")
	err := laService.InitApp(routerCfg)
	if err != nil{
		log.Fatal(err.Error())
	}
	log.Println("Running service on port ", routerCfg.PortNo)
	err = laService.RunLeagueAuctionServer()
	if err != nil{
		log.Fatal(err.Error())
	}
}