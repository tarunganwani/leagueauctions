package router


import(
	"net/http"
	"errors"
	"regexp"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"github.com/leagueauctions/server/utils"

	// "fmt"
)

//HTTPMethodsMap - global map of http methods
var HTTPMethodsMap = map[string]interface{} {"GET":nil, "PUT":nil, "POST":nil, "DELETE":nil, "":nil,}

//IsValidHTTPMethod - checks whether method is in the valid set of methods
func IsValidHTTPMethod(methodname string) bool{
	_, found := HTTPMethodsMap[methodname]
	return found
}

//MuxWrapper - wrapper over gorilla mux
type MuxWrapper struct{
	router *mux.Router
	routerconfig Config
}

//Init - initialize Mux wrapper
func (m *MuxWrapper)Init(cfg Config) error{
	m.router = new(mux.Router)
	m.routerconfig = cfg
	return nil
}

//Serve - start the router to serve any requests
func (m *MuxWrapper)Serve() error{
	srvAdd := m.routerconfig.HostAddress + ":" + utils.IntToString(m.routerconfig.PortNo)

	allowedHeaders := []string{"X-Requested-With", "Content-Type", "Authorization"}
	allowedMethods := []string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}
	allowedOrigins := []string{"*"}

	if m.routerconfig.Secure == false{
		// fmt.Println("serve in-secure")
		return http.ListenAndServe(srvAdd, handlers.CORS(
				handlers.AllowedHeaders(allowedHeaders), 
				handlers.AllowedMethods(allowedMethods), 
				handlers.AllowedOrigins(allowedOrigins))(m.router))
	}

	// fmt.Println("serve secured..")
	return http.ListenAndServeTLS(srvAdd, m.routerconfig.CertFilePath, m.routerconfig.KeyPath, 
		handlers.CORS(
			handlers.AllowedHeaders(allowedHeaders), 
				handlers.AllowedMethods(allowedMethods), 
				handlers.AllowedOrigins(allowedOrigins))(m.router))
}

//HandleRoute - handle specific route
func (m *MuxWrapper)HandleRoute (route string, httpmethod string, 
								handler func (w http.ResponseWriter, r *http.Request)) error{
	
	// fmt.Println("Setting handler for route ", route)
	// fmt.Println("m ", m)
	// fmt.Println("m.router ", m.router)

	if m.router == nil {
		return errors.New("mux router object can not be nil")
	}
	if false == IsValidHTTPMethod(httpmethod){
		return errors.New("Invalid http method "+ httpmethod)
	}
	muxRoute, err := m.TranslateRoute(route)
	if (err != nil){
		return err
	}

	if httpmethod != ""{
		m.router.HandleFunc(muxRoute, handler).Methods(httpmethod)
	} else {
		//for cases like websocket upgrade
		m.router.HandleFunc(muxRoute, handler)
	}
	return nil
}
//FetchRequestVar - fetch value of variable from mux vars
func (m *MuxWrapper)FetchRequestVar(r *http.Request, varname string) string{
	vars := mux.Vars(r)
	return vars[varname]
}

//ServeHTTP - serve http
func (m *MuxWrapper)ServeHTTP(w http.ResponseWriter, r *http.Request) error{
	if m.router == nil {
		return errors.New("router object can not be nil")
	}
	m.router.ServeHTTP(w, r)
	return nil
}

//TranslateRoute : translates league auction custom routes to mux route format 
//example: "{id:number}"" --> "{id:[0-9]+}"
func (m *MuxWrapper)TranslateRoute(route string) (string, error){

	varPresenceRe, err := regexp.Compile(":")
	if (err != nil){
		return route, err
	}
	if (varPresenceRe.MatchString(route) == false){
		return route, nil	// No variables present, return the input string as is
	}

	//variables present
	re, err := regexp.Compile("(?i)(number)")
	if (err != nil){
		return route, err
	}

	if (re.MatchString(route) == false){
		return route, errors.New("unsupported datatype")	
	}

	result := re.ReplaceAllStringFunc(route, m.datatypeReplacer)
	return result, nil
}

func (m *MuxWrapper)datatypeReplacer(s string) string {
    d := map[string]string{
        "number":       "[0-9]+",
        "NUMBER":       "[0-9]+",
    }
    r, ok := d[s]
    if ok {
        return r
    }
    return s
}