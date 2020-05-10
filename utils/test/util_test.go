package test

import(
	"fmt"
	"testing"
	"github.com/leagueauctions/utils"
)

var numgenTests = []struct {
	dig		int
	min		int
	max		int
}{
	{3, 100, 999},
	{4, 1000, 9999},
	{6, 100000, 999999},
}

func TestNDigitRandomNumber(t *testing.T){
	
	for _, val := range numgenTests{

		//10 random numbers
		for i := 0; i < 10; i++{
			num := utils.GenerateRandomNumber(val.dig)
			if (num < val.min && num >= val.max){
				t.Fatal("expected number in range (", val.min, " - ", val.max, ") actual = ", num)
			}
			// fmt.Println(num, "dig = ", val.dig, " num = ", num)
		}
	}
}

func TestSaltedHash(t * testing.T){
	pass := string("Password123")
	hash, err := utils.HashPassword(pass)
	if err != nil{
		t.Fatal(err)
	}
	fmt.Println("pwd ", pass)
	fmt.Println("hash ", hash)
	validPass, err := utils.ComparePasswords(hash, pass)
	if err != nil{
		t.Fatal(err)
	}
	if validPass == false{
		t.Fatal("Password and hash should match")
	}
}
