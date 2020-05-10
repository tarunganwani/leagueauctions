package test

import(
	"testing"
	"github.com/leagueauctions/router"
	// "fmt"
)

var translateTests = []struct {
	in  string
	out string
	errExpected bool
}{
	{"/users", "/users", false},
	{"/user/{userid:number}", "/user/{userid:[0-9]+}", false},
	{"/user/{userid:alphanum}", "", true},
}

func TestTranslateRoute(t *testing.T){

	muxWrapper := router.MuxWrapper{}
	for _, tt := range translateTests{
		translatedRoute, err := muxWrapper.TranslateRoute(tt.in)
		if !tt.errExpected && err != nil{
			t.Fatal("input ", tt.in, "error", err)
			if translatedRoute != tt.out {
				t.Fatal("Input", tt.in, "Expected", tt.out, "Actual ", translatedRoute)
			}
		} else if tt.errExpected && err == nil{
			t.Fatal("Error expected while transalating input ", tt.in)
		}
	}
}
