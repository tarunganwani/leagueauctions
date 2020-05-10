package utils

import(
	"math"
	"math/rand"
	"time"
)

//GenerateRandomNumber - 
func GenerateRandomNumber(dig int) int{
	rand.Seed(time.Now().UnixNano())
    min := int(math.Pow10(dig))
    max := int(math.Pow10(dig+1) - 1)
	return rand.Intn(max - min + 1) + min
}