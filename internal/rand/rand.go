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
	// returns a integer in the closed interval [1, ceiling]

	// edge case: ceiling = 0
	if ceiling == 0 {
		return 1
	}

	var RandomIntURL string
	if checkIfTestMode() {
		// stop program from requesting random.org if it is in test mode
		RandomIntURL = fmt.Sprintf("https://rng-api.fastapicloud.dev/random/%d", ceiling)
	} else {
		RandomIntURL = fmt.Sprintf("https://www.random.org/integers/?num=1&min=1&max=%d&col=1&base=10&format=plain&rnd=new", ceiling)
	}
	result, err := MakeRequest(RandomIntURL)
	if err != nil || result == -1 {
		result := rand.Intn(ceiling)
		return result
	}
	return result
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
