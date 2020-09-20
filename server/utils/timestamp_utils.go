package utils 

import(
	"time"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

//ConvertTimeToTimestampProto - convert time to timestamp proto
func ConvertTimeToTimestampProto(t time.Time) (*timestamp.Timestamp, error){
	return ptypes.TimestampProto(t)
}