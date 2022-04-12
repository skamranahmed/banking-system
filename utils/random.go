package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	INR = "INR"
	USD = "USD"
)

const alphabets = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt : generate a random int between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max+1-min) // min -> max
}

// RandomString : generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabets)

	for i := 0; i < n; i++ {
		c := alphabets[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner : returns a random Owner string
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney : returns a random money amount
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency : returns a random `defined` currency
func RandomCurrency() string {
	currencies := []string{INR, USD}
	n := len(currencies)
	return currencies[rand.Intn(n)]

}

// RandomEmail : returns a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
