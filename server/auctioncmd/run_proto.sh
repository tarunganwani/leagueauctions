export LEAGUE_AUCTIONS_PATH=/home/tg/go/src/github.com/leagueauctions
# export TIMESTAMP_INCLUDE_PATH=/home/tg/install/protoc_3_13/include/google/protobuf
cd $LEAGUE_AUCTIONS_PATH/server/auctioncmd
# protoc -I. -I$TIMESTAMP_INCLUDE_PATH --go_out=$LEAGUE_AUCTIONS_PATH/server/auctioncmd proto/*.proto
protoc  --go_out=$LEAGUE_AUCTIONS_PATH/server/auctioncmd proto/*.proto
mv $LEAGUE_AUCTIONS_PATH/server/auctioncmd/proto/*.pb.go $LEAGUE_AUCTIONS_PATH/server/auctioncmd