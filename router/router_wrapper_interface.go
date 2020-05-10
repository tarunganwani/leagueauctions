package router

import(
	"net/http"
)

//Wrapper - interface/contract to create/derive any router for the app
type Wrapper interface{
	Init(routerCfg Config) error
	Serve() error
	HandleRoute(route string, 
				httpmethod string, 
				handler func (w http.ResponseWriter, r *http.Request)) error
	FetchRequestVar(r *http.Request, varname string) string
	ServeHTTP(w http.ResponseWriter, r *http.Request) error
}