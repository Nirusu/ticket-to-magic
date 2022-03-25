package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Might need adjustment if you get a 403 error back from the API
const UserAgent = "Mozilla/5.0 (X11; Linux x86_64; rv:98.0) Gecko/20100101 Firefox/98.0"

// These are the names the API response uses for the facilities
const CaliforniaAdventurePark = "DLR_CA"
const DisneylandPark = "DLR_DP"

// URLs to query for tokens & availability
const TokenURL = "https://disneyland.disney.go.com/com-shared/api/get-token/"
const AvailabilityURL = "https://cme-dlr.wdprapps.disney.com/availability/api/v2/availabilities/?sku=66282&sku=66283"

// Translation from the short names to a more human-readable output
var ParkTranslation = map[string]string{
	"DLR_CA": "California Adventure Park",
	"DLR_DP": "Disneyland Park",
}

// The time layout used for time parsing
const DateLayout = "2006-01-02"

// TokenResponse is the response from the API when we get the access token
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   string `json:"expires_in"`
	StatusCode  int    `json:"status_code"`
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	ValidUntil  time.Time
}

// JSON structures for the availability response
type CalenderAvailabilityResponse struct {
	Availabilities []SingleDay `json:"calendar-availabilities"`
}

type SingleDay struct {
	Date         string     `json:"date"`
	Availability string     `json:"availability"`
	Facilities   []Facility `json:"facilities"`
}

type Facility struct {
	FacilityName string `json:"facilityName"`
	Available    bool   `json:"available"`
	Blocked      bool   `json:"blocked"`
}

func main() {
	// Parse target time from command line (YYYY-MM-DD)
	targetDate := os.Args[1]
	targetTime, err := time.Parse(DateLayout, targetDate)
	if err != nil {
		panic(err)
	}

	var accessToken *AccessToken
	for range time.Tick(time.Second * 10) {
		// Check if token is valid. Use 5 second tolerance to avoid errors
		if accessToken == nil || time.Now().After(accessToken.ValidUntil) {
			log.Println("Requesting new token...")
			accessToken, err = GetAccessToken()
			if err != nil {
				panic(err)
			}
		}

		calenderAvailabilityResponse, err := QueryAvailability(accessToken.AccessToken)
		if err != nil {
			panic(err)
		}

		// Loop through the response, check for dates with availability and check if it's before the target date we got.
		for _, singleDay := range calenderAvailabilityResponse.Availabilities {
			if singleDay.Availability != "cms-key-no-availability" {
				availableTime, err := time.Parse(DateLayout, singleDay.Date)
				if err != nil {
					panic(err)
				}
				for _, singleFacility := range singleDay.Facilities {
					if singleFacility.Available {
						if availableTime.Before(targetTime) {
							fmt.Printf("%s: %s is available\n", singleDay.Date, ParkTranslation[singleFacility.FacilityName])
						}
					}
				}
			}
		}
	}
}

func GetAccessToken() (*AccessToken, error) {
	// Setup HTTP client
	client := &http.Client{}
	req, err := http.NewRequest("POST", TokenURL, nil)
	if err != nil {
		return nil, err
	}

	// Fake user-agent to get access token. The default Go user-agent gives us a "Permission Denied" error back.
	req.Header.Set("User-Agent", UserAgent)

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if User Agent was blacklisted / not accepted
	if resp.StatusCode == 403 {
		panic("403 Forbidden - Change user agent?")
	}

	// Parse the response, get the access token
	var tokenResponse TokenResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&tokenResponse); err != nil {
		return nil, err
	}
	resp.Body.Close()

	var accessToken AccessToken
	accessToken.AccessToken = tokenResponse.AccessToken
	expiresIn, err := strconv.Atoi(tokenResponse.ExpiresIn)
	if err != nil {
		return nil, err
	}
	accessToken.ValidUntil = time.Now().Add(time.Duration(expiresIn-5) * time.Second)

	return &accessToken, nil
}

func QueryAvailability(accessTokenString string) (*CalenderAvailabilityResponse, error) {
	// Create request to create availability
	req, err := http.NewRequest("GET", AvailabilityURL, nil)
	if err != nil {
		return nil, err
	}

	// Same as before, fake the user agent and use the authorization token we got before
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Authorization", "Bearer "+accessTokenString)

	client := &http.Client{}
	// Query for availability...
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response
	var calenderAvailabilityResponse CalenderAvailabilityResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&calenderAvailabilityResponse); err != nil {
		return nil, err
	}

	return &calenderAvailabilityResponse, nil
}
