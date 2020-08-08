package router


import(
	"net/http"
	"net/http/httptest"
)


func executeRequest(r Wrapper, req *http.Request) *httptest.ResponseRecorder {
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    return rr
}