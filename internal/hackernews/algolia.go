package hackernews

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"charm.land/log/v2"
)

const (
	userURL      = "https://hn.algolia.com/api/v1/users/%s"
	searchURL    = "https://hn.algolia.com/api/v1/search?query=%s&tags=front_page"
	maxBonusHype = 100
)

type userResponse struct {
	Karma int `json:"karma"`
}

type searchResponse struct {
	NbHits int `json:"nbHits"`
}

var httpClient = &http.Client{Timeout: 5 * time.Second}

// FetchBonusHype looks up a Hacker News username via the Algolia API and
// returns a bonus hype value based on karma and front-page hits.
// Returns 0 if the username is not found or any API call fails.
func FetchBonusHype(username string) int {
	log.Info("FetchBonusHype called", "username", username)
	if username == "" {
		log.Info("FetchBonusHype skipped, empty username")
		return 0
	}

	karma, err := fetchKarma(username)
	if err != nil {
		log.Info("FetchBonusHype failed to fetch karma", "username", username, "error", err)
		return 0
	}

	nbHits, err := fetchFrontPageHits(username)
	if err != nil {
		log.Info("FetchBonusHype failed to fetch front page hits", "username", username, "error", err)
		return 0
	}

	bonus := CalcBonusHype(karma, nbHits)
	log.Info("FetchBonusHype result", "username", username, "karma", karma, "nbHits", nbHits, "bonusHype", bonus)
	return bonus
}

// CalcBonusHype computes bonusHype = min(floor(log10(karma + 1)) * 5 + nbHits * 2, 30).
func CalcBonusHype(karma, nbHits int) int {
	log.Info("CalcBonusHype called", "karma", karma, "nbHits", nbHits)
	bonus := int(math.Floor(math.Log10(float64(karma+1))))*10 + nbHits*2
	if bonus > maxBonusHype {
		bonus = maxBonusHype
	}
	if bonus < 0 {
		bonus = 0
	}
	log.Info("CalcBonusHype result", "bonusHype", bonus)
	return bonus
}

func fetchKarma(username string) (int, error) {
	log.Info("fetchKarma called", "username", username)
	resp, err := httpClient.Get(fmt.Sprintf(userURL, username))
	if err != nil {
		log.Info("fetchKarma request failed", "username", username, "error", err)
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Info("fetchKarma user not found", "username", username, "status", resp.StatusCode)
		return 0, fmt.Errorf("user not found: status %d", resp.StatusCode)
	}

	var u userResponse
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		log.Info("fetchKarma decode failed", "username", username, "error", err)
		return 0, err
	}
	log.Info("fetchKarma result", "username", username, "karma", u.Karma)
	return u.Karma, nil
}

func fetchFrontPageHits(username string) (int, error) {
	log.Info("fetchFrontPageHits called", "username", username)
	resp, err := httpClient.Get(fmt.Sprintf(searchURL, username))
	if err != nil {
		log.Info("fetchFrontPageHits request failed", "username", username, "error", err)
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Info("fetchFrontPageHits bad status", "username", username, "status", resp.StatusCode)
		return 0, fmt.Errorf("search failed: status %d", resp.StatusCode)
	}

	var s searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		log.Info("fetchFrontPageHits decode failed", "username", username, "error", err)
		return 0, err
	}
	log.Info("fetchFrontPageHits result", "username", username, "nbHits", s.NbHits)
	return s.NbHits, nil
}
