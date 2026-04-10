package svt

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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	testMode := os.Getenv("TEST_MODE")

	return testMode == "True"
}

func GetRandomInt(ceiling int) int {
	// returns a integer in the closed interval [1, ceiling]

	if checkIfTestMode() {
		// stop program from requesting random.org if it is in test mode
		fmt.Println("TEST MODE: Using random.org is disabled")
		result := rand.Intn(ceiling) + 1
		return result
	}
	fmt.Println("TEST MODE: Using random.org is enabled")
	RandomIntURL := fmt.Sprintf("https://www.random.org/integers/?num=1&min=1&max=%d&col=1&base=10&format=plain&rnd=new", ceiling)

	request := NewGetRandomIntRequest(RandomIntURL)
	response, err := NewGetRandomIntResponseFromClient(request)
	if err != nil {
		// backup
		result := rand.Intn(ceiling)
		return result
	}
	defer response.Body.Close()

	result := ExtractRandomInteger(response)
	return result
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

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
