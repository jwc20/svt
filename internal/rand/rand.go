package rand

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func checkIfTestMode() bool {
	// Best-effort load; .env may not exist (e.g. in production or CI).
	_ = godotenv.Load()
	testMode := os.Getenv("TEST_MODE") == "True"

	//if testMode {
	//	log.Print("TEST MODE: Using random.org is disabled")
	//} else {
	//	log.Print("TEST MODE: Using random.org is enabled")
	//}
	return testMode
}

func GetRandomInt(ceiling int) int {
	// Returns an integer in the closed interval [1, ceiling]
	if ceiling == 0 {
		return 1
	}

	fallbackURL := fmt.Sprintf("https://rng-api.fastapicloud.dev/random/%d", ceiling)
	localFallback := func() int { return rand.Intn(ceiling) }

	tryRequest := func(url string) (int, bool) {
		result, err := MakeRequest(url)
		if err != nil || result == -1 {
			return 0, false
		}
		return result, true
	}

	if checkIfTestMode() {
		if result, ok := tryRequest(fallbackURL); ok {
			return result
		}
		return localFallback()
	}

	// Production: try random.org → fallback API → local rand
	primaryURL := fmt.Sprintf(
		"https://www.random.org/integers/?num=1&min=1&max=%d&col=1&base=10&format=plain&rnd=new", ceiling,
	)
	if result, ok := tryRequest(primaryURL); ok {
		return result
	}
	if result, ok := tryRequest(fallbackURL); ok {
		return result
	}
	return localFallback()
}

func MakeRequest(url string) (int, error) {
	request := NewGetRandomIntRequest(url)
	response, err := NewGetRandomIntResponseFromClient(request)
	if err != nil {
		fmt.Println("Error: ", err)
		return -1, err
	}
	defer response.Body.Close()

	result := ExtractRandomInteger(response)
	return result, nil
}

func NewGetRandomIntRequest(url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	return req
}

func NewGetRandomIntResponseFromClient(req *http.Request) (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return resp, err
}

func ExtractRandomInteger(resp *http.Response) int {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	result, err := strconv.Atoi(strings.TrimSpace(string(body)))
	if err != nil {
		log.Fatal(err)
	}
	return result
}
