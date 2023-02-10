package functions

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	G "github.com/Brum3ns/firefly/pkg/functions/globalVariables"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

/** Random functions Makes random values in string and int with "x" length*/
func RandomInsert(s string) string {
	var (
		l = []string{"#RANDOM#", "#RANDOMNUM#"}
		r string
	)

	for key, i := range G.M_Random {
		switch key {
		case "s":
			r = l[0]
		case "n":
			r = l[1]
		}

		if strings.Contains(s, r) {
			s = strings.ReplaceAll(s, r, RandomCreate(key, i))
		}
	}

	return s
}

/** Craft a random str/int value with "x" length - "[t]ype:[l]ength" Return the random value */
func RandomCreate(t string, l int) string {
	switch t {
	case "n":
		return RandNumber(l)
	default: //Default = [s]tring
		return RandString(l)
	}
}

/**Return random number with given length as string*/
func RandNumber(n int) string {
	rand.Seed(time.Now().UnixNano())
	if n <= 0 {
		n = 1
	}
	randint := (rand.Float64() - 0.01) * (math.Pow(1*10, (float64)(n)))
	return fmt.Sprintf("%.0f", randint)
}

/**Return random string with given length*/
func RandString(n int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
