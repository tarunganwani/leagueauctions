export LEAGUE_AUCTIONS_PATH=/home/tg/go/src/github.com/leagueauctions
cd $LEAGUE_AUCTIONS_PATH/server/auctioncmd
protoc --go_out=$LEAGUE_AUCTIONS_PATH/server/auctioncmd proto/*.proto
mv $LEAGUE_AUCTIONS_PATH/server/auctioncmd/proto/*.pb.go $LEAGUE_AUCTIONS_PATH/server/auctioncmd