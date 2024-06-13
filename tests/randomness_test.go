package tests

import (
	"fmt"
	"log"
	"testing"

	"github.com/Brum3ns/firefly/pkg/random"
	"github.com/Brum3ns/firefly/pkg/randomness"
)

func Test_RandomnessAccuracy(t *testing.T) {
	// Config
	var (
		amountToTest         = 100
		lengthOfRandomString = 16
	)

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
	lst_random := getRandomStrings(lengthOfRandomString, amountToTest)
	// lst_valid := getValidStrings()

	// Check values if they are random
	hit := 0
	miss := 0
	for _, i := range lst_random {
		if r.IsRandom(i) {
			//fmt.Println("RANDOM:", i)
			hit++
		} else {
			//fmt.Println("NORMAL:", i)
			miss++
		}
	}

	// Show result
	fmt.Printf("\n===Result===\nHit:%d, Miss:%d (%f%%)\n============\n", hit, miss, float64(miss)/float64(hit)*10)
}

func getRandomStrings(nr, amount int) []string {
	var lst []string
	for i := 0; i < amount; i++ {
		// nr, _ := strconv.Atoi(random.RandNumber(2))
		lst = append(lst, random.RandString(nr))

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
