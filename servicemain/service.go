package servicemain

import(
	"github.com/leagueauctions/router"
	"github.com/leagueauctions/usermgmt"
)

//LeagueAuction - main app structure
type LeagueAuction struct{
	routerCfg router.Config
	router router.Wrapper
}

func (la *LeagueAuction)initUserMgmtRoutes(r router.Wrapper) error{
	usrMgmtRouter := new(usermgmt.Router)
	return usrMgmtRouter.Init(r)
}

//InitApp - initialize league auction server
func (la *LeagueAuction)InitApp(routerCfg router.Config) error{
	la.router = router.MuxWrapper{}
	err := la.router.Init(routerCfg)
	if err != nil{
		return err
	}
	err = la.initUserMgmtRoutes(la.router)
	if err != nil{
		return err
	}
	return nil
}


//RunLeagueAuctionServer - run league auction server
func (la *LeagueAuction)RunLeagueAuctionServer() error{
	return la.router.Serve()
}