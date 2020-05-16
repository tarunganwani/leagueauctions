package servicemain

import(
	"database/sql"
	"github.com/leagueauctions/server/router"
	"github.com/leagueauctions/server/usermgmt"
	"github.com/leagueauctions/server/utils"
)

//LeagueAuction - main app structure
type LeagueAuction struct{
	routerCfg router.Config
	router router.Wrapper
	dbObject *sql.DB
}

func (la *LeagueAuction)initDatabase() (*sql.DB, error){

	//TODO: move hardcoded values to database config object
	dbObject, err := utils.OpenPostgreDatabase("postgres", "postgres", "leagueauction")
	if err != nil{
		return nil, err
	}
	return dbObject, nil

}

func (la *LeagueAuction)initUserMgmtRoutes(r router.Wrapper, dbObject *sql.DB) error{

	usrMgmtRouter := new(usermgmt.Router)
	return usrMgmtRouter.Init(r, dbObject)
}

//InitApp - initialize league auction server
func (la *LeagueAuction)InitApp(routerCfg router.Config) error{
	la.router = new(router.MuxWrapper)
	err := la.router.Init(routerCfg)
	if err != nil{
		return err
	}
	dbObj, err := la.initDatabase()
	if err != nil{
		return err
	}
	err = la.initUserMgmtRoutes(la.router, dbObj)
	if err != nil{
		return err
	}
	return nil
}


//RunLeagueAuctionServer - run league auction server
func (la *LeagueAuction)RunLeagueAuctionServer() error{
	return la.router.Serve()
}