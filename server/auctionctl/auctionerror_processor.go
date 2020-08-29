package auctionctl

import(
	pb "github.com/leagueauctions/server/auctioncmd" 
)

func generateAuctionResponse(errorStr string) (*pb.AuctionResponse){
	auctionResponse := new(pb.AuctionResponse)
	auctionResponse.ResponseType = pb.AuctionResponse_ERROR
	auctionResponse.Errormsg = errorStr
	return auctionResponse
}