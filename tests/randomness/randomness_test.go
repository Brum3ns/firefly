package tests

import (
	"fmt"
	"log"
	"testing"

	"github.com/Brum3ns/firefly/pkg/randomness"
	"github.com/brianvoe/gofakeit/v7"
)

func Test_RandomnessAccuracy(t *testing.T) {
	//defaultConfig := randomness.DefaultConfig()
	config := randomness.Config{
		InRow:      randomness.DEFAULT_INROW,
		Vocal:      randomness.DEFAULT_VOCAL,
		Digit:      randomness.DEFAULT_DIGIT,
		Consonant:  randomness.DEFAULT_CONSONANT,
		Blacklist:  randomness.DEFAULT_BLACKLISTS,
		BlackRegex: randomness.DEFAULT_BLACKREGEX,
		Spaces:     []rune{' ', '_', '-', '.'},
	}

	// Setup randomness config
	r, err := randomness.NewRandomness(config)
	if err != nil {
		log.Println(err)
	}

	// Config random strings to test
	//var (
	//	amountToTest         = 10000
	//	lengthOfRandomString = 16
	//)
	//lst := getRandomStrings(lengthOfRandomString, amountToTest)
	lst := getTestItems()
	//lst_valid := getValidStrings()

	// Check values if they are random
	hit := 0
	miss := 0
	for _, i := range lst {
		if r.IsRandom(i) {
			//fmt.Println("RANDOM:", i)
			hit++
		} else {
			//fmt.Println("NORMAL:", i)
			miss++
			fmt.Println(i, randomness.IsRandomByEntropy(i, float64(3500/1000)), randomness.CalcEntropy(i))
		}
	}

	// Show result
	fmt.Printf("\n===Result===\nHit:%d, Miss:%d (%f%%)\n============\n", hit, miss, float64(miss)/float64(hit)*10)
}

func getRandomStrings(nr, amount int) []string {
	var lst []string
	for i := 0; i < amount; i++ {
		// nr, _ := strconv.Atoi(random.RandNumber(2))
		v, err := gofakeit.Generate("{regex:\\w{32,32}}")
		if err != nil {
			log.Fatalln(err)
		}

		lst = append(lst, v)

	}
	return lst
}

func getValidStrings() []string {
	return []string{
		"username",
		"master",
		"testThisstuff",
		"works",
		"cat",
		"PillarTown",
	}
}

func getTestItems() []string {
	return []string{
		"X5Fh2z9M8Q",
		"3f09f39b8a",
		"YWJjMTIzIT8kKiYoKSctPUB+",
		"ae2b1fca515949e5d54fb22b8ed95575",
		"dGhpc2lzYXJhbmRvbXN0cmluZw==",
		"bdf9e2cd-7cf1-41cb-9a89-e327d56f354b",
		"A1B2C3D4E5",
		"xj29vk3sld93",
		"U2FsdGVkX1+g9JZK9g0dP0e3",
		"098f6bcd4621d373cade4e832627b4f6",
		"hello",
		"user_id",
		"status=200",
		"July 4, 2025",
		"/robots.txt",
		"index.html",
		"login",
		"error_message",
		"Connection Timeout",
		"api/v1/users",
		"qwerty123",
		"temp_var",
		"auth_token",
		"admin_user",
		"config_value",
		"Test1234",
		"__init__",
		"MainActivity",
		"fileNotFound",
		"sha256sum",
	}
}
