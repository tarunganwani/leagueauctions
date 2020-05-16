package test

import (
	"testing"
	"github.com/leagueauctions/server/servicemain"
	"github.com/leagueauctions/server/router"
)

func TestLaService(t *testing.T) {

	routerCfg := router.Config{HostAddress: "localhost", PortNo : 8081}
	laService := new(servicemain.LeagueAuction)
	err := laService.InitApp(routerCfg)
	if err != nil{
		t.Fatal(err.Error())
	}
	err = laService.RunLeagueAuctionServer()
	if err != nil{
		t.Fatal(err.Error())
	}
}