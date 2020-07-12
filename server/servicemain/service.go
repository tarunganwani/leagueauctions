package servicemain

import(
	"database/sql"
	"github.com/leagueauctions/server/libs/router"
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

func (la *LeagueAuction)setupEndpoints(r router.Wrapper, dbObject *sql.DB) error{
	
	if err := la.initUserMgmtRoutes(r, dbObject); err != nil{
		return err
	}
	return nil
}

//InitApp - initialize league auction server
func (la *LeagueAuction)InitApp(routerCfg router.Config) (err error){
	la.router = new(router.MuxWrapper)
	if err = la.router.Init(routerCfg); err != nil{
		return
	}
	var dbObj *sql.DB
	if dbObj, err = la.initDatabase(); err != nil{
		return
	}
	if err = la.setupEndpoints(la.router, dbObj); err != nil{
		return
	}
	return
}


//RunLeagueAuctionServer - run league auction server
func (la *LeagueAuction)RunLeagueAuctionServer() error{
	return la.router.Serve()
}