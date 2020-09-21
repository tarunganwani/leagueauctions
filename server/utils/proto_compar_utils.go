package utils

import(
	"time"
	"reflect"
	pb "github.com/leagueauctions/server/auctioncmd" 
	"github.com/golang/protobuf/ptypes/timestamp"
)

//CompareCategoryListPb - compares category list protobuf objects
func CompareCategoryListPb(pbCatList1, pbCatList2 []*pb.PlayerCategory) bool{
	m1 := make(map [string]uint64)
	m2 := make(map [string]uint64)
	for _, c := range pbCatList1 {
		m1[c.CategoryName] = c.PlayerBasePrice
	}
	for _, c := range pbCatList2 {
		m2[c.CategoryName] = c.PlayerBasePrice
	}
	return reflect.DeepEqual(m1, m2)
}


//CompareDateTimePb - compare protobuf times
func CompareDateTimePb(dt1Pb, dt2Pb *timestamp.Timestamp) bool{
	dt1 := dt1Pb.AsTime()
	dt2 := dt2Pb.AsTime()
	dt1Round := dt1.Round(time.Second)
	dt2Round := dt2.Round(time.Second)
	return dt1Round.Equal(dt2Round)
}