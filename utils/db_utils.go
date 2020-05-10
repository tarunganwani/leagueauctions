package utils

import(
	"fmt"
	"database/sql"
)

//getPostgresConnectionString - gets connection string for postgres
func getPostgresConnectionString(user, password, dbname string) string{
	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
}

//OpenPostgreDatabase - open postgre connection
func OpenPostgreDatabase(user, pwd, dbname string) (*sql.DB, error){
	connectionString := getPostgresConnectionString(user, pwd, dbname)
	return sql.Open("postgres", connectionString)
}