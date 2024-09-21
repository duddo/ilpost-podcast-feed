package ilpostapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

const BaseURL = "https://api-prod.ilpost.it"

const LoginURL = "https://www.ilpost.it/wp-login.php"

func Login(username, password string) ([]*http.Cookie, error) {
	client := &http.Client{}

	// Create the POST payload
	payload := "log=" + url.QueryEscape(username) + "&pwd=" + url.QueryEscape(password) + "&rememberme=forever&testcookie=1"
	req, err := http.NewRequest("POST", LoginURL, bytes.NewBufferString(payload))
	if err != nil {
		log.Println("Error creating request:", err)
		return nil, err
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Perform the login request
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error making login request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Check for errors in the response
	if resp.StatusCode != http.StatusOK {
		log.Printf("Login failed with status code %d\n", resp.StatusCode)
		return nil, fmt.Errorf("login failed: %s", resp.Status)
	}

	var cookies = resp.Cookies()
	if len(cookies) < 2 {
		return nil, fmt.Errorf("login failed, no cookie returned")
	}

	return resp.Cookies(), nil
}

func FetchPodcastList(cookies []*http.Cookie) (PodcastListResponse, error) {
	var response PodcastListResponse
	client := &http.Client{}

	// Create the GET request
	req, err := http.NewRequest("GET", BaseURL+"/frontend/podcast/list", nil)
	if err != nil {
		log.Println("Error creating GET request:", err)
		return response, err
	}

	// Include the Cookie header
	req.Header.Set("Cookie", GetCookieString(cookies))

	// Make the GET request
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error making GET request:", err)
		return response, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("received status code %s", http.StatusText(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return response, err
	}

	// Deserialize the JSON response
	if err := json.Unmarshal(body, &response); err != nil {
		log.Println("Error decoding JSON:", err)
		return response, err
	}

	if response.Data == nil {
		return response, fmt.Errorf("empty data on response")
	}

	return response, nil
}

func FetchPodcastEpisodes(cookies []*http.Cookie, slug string, page int, perPage int) (PodcastEpisodesResponse, error) {
	var response PodcastEpisodesResponse
	client := &http.Client{}

	// Create the GET request
	url := fmt.Sprintf("%s%s%s?&pg=%v&hits=%v", BaseURL, "/v1/podcast/", slug, page, perPage)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return response, fmt.Errorf("can't create GET request %s: %s", url, err.Error())
	}

	// Include the Cookie header
	cookieString := GetCookieString(cookies)
	req.Header.Set("Cookie", cookieString)

	// Make the GET request
	resp, err := client.Do(req)
	if err != nil {
		return response, fmt.Errorf("can't make GET request %s: %s", url, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("received status code %s from %s", http.StatusText(resp.StatusCode), url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("can't GET response body from %s: %s", url, err.Error())
	}

	// Deserialize the JSON response
	if err := json.Unmarshal(body, &response); err != nil {
		return response, fmt.Errorf("can't decode JSON of GET request %s: %s", url, err.Error())
	}

	return response, nil
}

func GetCookieString(cookies []*http.Cookie) string {

	var cookieString string

	for _, cookie := range cookies {
		cookieString += cookie.String() + "; "
	}

	// Trim the trailing "; "
	if len(cookieString) > 2 {
		cookieString = cookieString[:len(cookieString)-2]
	}

	return cookieString
}
