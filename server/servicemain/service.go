package servicemain

import(
	"database/sql"
	"github.com/leagueauctions/server/libs/router"
	"github.com/leagueauctions/server/usermgmt"
	"github.com/leagueauctions/server/utils"
	"github.com/leagueauctions/server/auctionctl"
	"github.com/leagueauctions/server/database"
)

//LeagueAuction - main app structure
type LeagueAuction struct{
	routerCfg router.Config
	router *router.MuxWrapper
	dbObject *sql.DB
}

func (la *LeagueAuction)initDatabase() (*database.LeagueAuctionDatastore, error){

	//TODO: move hardcoded values to database config object
	dbObject, err := utils.OpenPostgreDatabase("postgres", "postgres", "leagueauction")
	if err != nil{
		return nil, err
	}
	datastore, err := database.GetLeagueAuctionDatastore(dbObject)
	return datastore, nil

}

func (la *LeagueAuction)initUserMgmtRoutes(r *router.MuxWrapper, 
										laDatastore *database.LeagueAuctionDatastore) error{

	usrMgmtRouter := new(usermgmt.Router)
	return usrMgmtRouter.Init(r, laDatastore.LAUserstore)
}

// func (la *LeagueAuction)initAuctionRoutes(r *router.MuxWrapper, dbObject *sql.DB, conPool *auctionctl.UserConnectionPool) error{

// 	auctionRouter := new(auctionctl.Router)
// 	return auctionRouter.Init(r, dbObject,conPool)
// }

func (la *LeagueAuction)setupEndpoints(r *router.MuxWrapper, 
										laDatastore *database.LeagueAuctionDatastore, 
										conPool *auctionctl.UserConnectionPool) error{
	
	if err := la.initUserMgmtRoutes(r, laDatastore); err != nil{
		return err
	}
	// if err := la.initAuctionRoutes(r, dbObject, conPool); err != nil{
	// 	return err
	// }
	return nil
}

//InitApp - initialize league auction server
func (la *LeagueAuction)InitApp(routerCfg router.Config) (err error){
	la.router = new(router.MuxWrapper)
	if err = la.router.Init(routerCfg); err != nil{
		return
	}
	var laDatastore *database.LeagueAuctionDatastore
	if laDatastore, err = la.initDatabase(); err != nil{
		return
	}
	
	
	// initialize user connection pool
	userConnPool := new(auctionctl.UserConnectionPool)
	userConnPool.Init()

	if err = la.setupEndpoints(la.router, laDatastore, userConnPool); err != nil{
		return
	}
	return
}


//RunLeagueAuctionServer - run league auction server
func (la *LeagueAuction)RunLeagueAuctionServer() error{
	return la.router.Serve()
}