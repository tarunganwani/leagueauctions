package converterutil

import(
	pb "github.com/leagueauctions/server/auctioncmd" 
	"github.com/leagueauctions/server/database" 
	"github.com/leagueauctions/server/utils" 
	"errors" 
	"github.com/google/uuid"
)

//CreateRequestToAuctionBoardDBObject - convert create auction board proto request into auction board db object
func CreateRequestToAuctionBoardDBObject (createAuctionBoardCmd *pb.CreateAuctionBoardRequest) (*database.AuctionBoard, error){

	if createAuctionBoardCmd == nil {
		return nil, errors.New("Nil proto object for CreateAuctionBoardRequest")
	}
	auctioneerUUID, err := utils.GetUUIDFromString(createAuctionBoardCmd.AuctioneerPlayerUuid)
	if err != nil{
		return nil, err
	}
	newBoardUUID := uuid.New()
	auctionBoardDbObj := new(database.AuctionBoard)
	auctionBoardDbObj.AuctionBoardUUID = newBoardUUID
	auctionBoardDbObj.Purse = createAuctionBoardCmd.PurseMoney
	auctionBoardDbObj.PurceCcy = createAuctionBoardCmd.PurseCcy
	auctionBoardDbObj.ScheduleTime = createAuctionBoardCmd.ScheduleTime.AsTime()
	auctionBoardDbObj.AuctioneerUUID = auctioneerUUID
	auctionBoardDbObj.IsActive = false
	auctionBoardDbObj.AuctionName = createAuctionBoardCmd.AuctionBoardName

	for _, catPb := range createAuctionBoardCmd.PlayerCategoryList{
		catDb := new(database.Category)
		catDb.AuctionBoardUUID = newBoardUUID
		catDb.BasePrice = catPb.PlayerBasePrice
		catDb.CategoryName = catPb.CategoryName
		catDb.CategoryUUID = uuid.New()
		auctionBoardDbObj.CategorySet = append(auctionBoardDbObj.CategorySet, catDb)
	}
	
	return auctionBoardDbObj, nil
}

//UpdateRequestToAuctionBoardDBObject - convert update auction board request to auction board database object
func UpdateRequestToAuctionBoardDBObject(updateAuctionBoardReqPb *pb.UpdateAuctionBoardRequest) (*database.AuctionBoard, error){


	auctionBoardUUID, err := utils.GetUUIDFromString(updateAuctionBoardReqPb.AuctionBoardUuid)
	if err != nil{
		return nil, err
	}
	auctionBoardDbObj := new(database.AuctionBoard)
	auctionBoardDbObj.AuctionBoardUUID = auctionBoardUUID
	auctionBoardDbObj.Purse = updateAuctionBoardReqPb.PurseMoney
	auctionBoardDbObj.PurceCcy = updateAuctionBoardReqPb.PurseCcy
	auctionBoardDbObj.ScheduleTime = updateAuctionBoardReqPb.ScheduleTime.AsTime()
	auctionBoardDbObj.AuctionName = updateAuctionBoardReqPb.AuctionBoardName
	return auctionBoardDbObj, nil
}

//GenerateFetchAuctionBoardResponse - convert auction db obj to fetch auction board response 
func GenerateFetchAuctionBoardResponse(auctionDbObj *database.AuctionBoard)(*pb.FetchAuctionBoardResponse, error){
	
	if auctionDbObj == nil{
		return nil, errors.New("Nil auction board object GenerateFetchAuctionBoardResponse")
	}
	timestampProto, err := utils.ConvertTimeToTimestampProto(auctionDbObj.ScheduleTime)
	if err != nil{
		return nil, errors.New("time conversion error GenerateFetchAuctionBoardResponse")
	}
	auctionBoardInfoPb := new(pb.AuctionBoardInfo)
	auctionBoardInfoPb.AuctionBoardName = auctionDbObj.AuctionName
	auctionBoardInfoPb.AuctionCode = auctionDbObj.AuctionCode
	auctionBoardInfoPb.IsActive = auctionDbObj.IsActive
	auctionBoardInfoPb.PurseCcy = auctionDbObj.PurceCcy
	auctionBoardInfoPb.PurseMoney = auctionDbObj.Purse
	auctionBoardInfoPb.ScheduleTime = timestampProto
	auctionBoardInfoPb.AuctioneerPlayerUuid = utils.GetStringFromUUID(auctionDbObj.AuctioneerUUID)

	fetchAuctionBoardResponsePb := new(pb.FetchAuctionBoardResponse)
	fetchAuctionBoardResponsePb.AuctionBoardInfo = auctionBoardInfoPb

	playerCatSlice := make([]*pb.PlayerCategory, len(auctionDbObj.CategorySet))
	sliceIndex := 0
	for _, catDB := range auctionDbObj.CategorySet{
		catPb := new(pb.PlayerCategory)
		catPb.CategoryName = catDB.CategoryName
		catPb.PlayerBasePrice = catDB.BasePrice
		catPb.CategoryUuid = utils.GetStringFromUUID(catDB.CategoryUUID)
		playerCatSlice[sliceIndex] = catPb
		sliceIndex++
	}
	fetchAuctionBoardResponsePb.AuctionBoardInfo.PlayerCategoryList = playerCatSlice

	return fetchAuctionBoardResponsePb, nil
}