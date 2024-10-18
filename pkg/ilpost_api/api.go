package ilpostapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const BaseURL = "https://api-prod.ilpost.it"

const LoginURL = "https://www.ilpost.it/wp-login.php"

func Login(username, password string) ([]*http.Cookie, error) {
	// Create the POST payload
	form := url.Values{}
	form.Add("log", username)
	form.Add("pwd", password)
	form.Add("wp-submit", "Accedi")
	form.Add("redirect_to", "https://www.ilpost.it/")
	form.Add("testcookie", "1")

	resp, err := http.Post(LoginURL,
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
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

	req, err := http.NewRequest("GET", BaseURL+"/frontend/podcast/list", nil)
	if err != nil {
		log.Println("Error creating GET request:", err)
		return response, err
	}

	req.Header.Set("Cookie", getCookieString(cookies))

	resp, err := http.DefaultClient.Do(req)
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

	url := fmt.Sprintf("%s%s%s?&pg=%v&hits=%v", BaseURL, "/podcast/v1/podcast/", slug, page, perPage)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return response, fmt.Errorf("can't create GET request %s: %s", url, err.Error())
	}

	req.Header.Set("Cookie", getCookieString(cookies))

	resp, err := http.DefaultClient.Do(req)
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

func getCookieString(cookies []*http.Cookie) string {
	var cookieHeader strings.Builder
	for _, cookie := range cookies {
		cookieHeader.WriteString(cookie.String() + "; ")
	}

	return strings.TrimRight(cookieHeader.String(), "; ")
}
