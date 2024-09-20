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

func Login(username, password string) (string, error) {
	client := &http.Client{}

	// Create the POST payload
	payload := "log=" + url.QueryEscape(username) + "&pwd=" + url.QueryEscape(password) + "&rememberme=forever&testcookie=1"
	req, err := http.NewRequest("POST", LoginURL, bytes.NewBufferString(payload))
	if err != nil {
		log.Println("Error creating request:", err)
		return "", err
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Perform the login request
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error making login request:", err)
		return "", err
	}
	defer resp.Body.Close()

	// Check for errors in the response
	if resp.StatusCode != http.StatusOK {
		log.Printf("Login failed with status code %d\n", resp.StatusCode)
		return "", fmt.Errorf("login failed: %s", resp.Status)
	}

	// Retrieve all cookies and concatenate them
	var cookieString string
	cookies := resp.Cookies()
	if len(cookies) < 2 {
		return "", fmt.Errorf("login failed, no cookie returned")
	}

	for _, cookie := range cookies {
		cookieString += cookie.String() + "; "
	}

	// Trim the trailing "; "
	if len(cookieString) > 2 {
		cookieString = cookieString[:len(cookieString)-2]
	}

	return cookieString, nil
}

func FetchPodcastList(cookieString *string) (PodcastListResponse, error) {
	var response PodcastListResponse
	client := &http.Client{}

	// Create the GET request
	req, err := http.NewRequest("GET", BaseURL+"/frontend/podcast/list", nil)
	if err != nil {
		log.Println("Error creating GET request:", err)
		return response, err
	}

	// Include the Cookie header
	if cookieString != nil {
		req.Header.Set("Cookie", *cookieString)
	}

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

func FetchPodcastEpisodes(cookieString *string, slug string, page int, perPage int) (PodcastEpisodesResponse, error) {
	var response PodcastEpisodesResponse
	client := &http.Client{}

	// Create the GET request
	url := fmt.Sprintf("%s%s%s?&pg=%v&hits=%v", BaseURL, "/v1/podcast/", slug, page, perPage)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error creating GET request:", err)
		return response, err
	}

	// Include the Cookie header
	if cookieString != nil {
		req.Header.Set("Cookie", *cookieString)
	}

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

	return response, nil
}
