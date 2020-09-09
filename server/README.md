# leagueauctions
create auctions app for sports leagues

# build dependncies
go get -u github.com/gorilla/mux
go get -u golang.org/x/crypto/bcrypt
go get -u github.com/lib/pq
go get -u github.com/dgrijalva/jwt-go
go get -u github.com/gorilla/handlers

# running the service(from code)
$ cd leagueauctions/servicemain/main
$ go run run.go

# Postgres issue

Problem #1 : By default postgres does not have any db password set on ubuntu. Now apparently this is an  issue while conencting to the database 

FAIL: TestNonExistentUser (0.01s)
	/home/tg/go/src/github.com/leagueauctions/server/usermgmt/test/model_test.go:48: pq: password authentication failed for user "postgres"

Solution for problem #1 :

change postgres user password. this is sorta required to connect to league auction database and access tables 
inside psql prompt
>> ALTER USER postgres PASSWORD 'newPassword';
https://stackoverflow.com/questions/7695962/postgresql-password-authentication-failed-for-user-postgres

# Proto buf dependencies

go get "github.com/golang/protobuf/ptypes/timestamp"
