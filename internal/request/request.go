package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	url = "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key="
)

// Response structure to match the expected JSON response
type Response struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func SendRequest(command, queryText, apiKey string) (string, error) {
	apiURL := url + apiKey

	// Create the request body with command and queryText
	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{
						"text": command + queryText,
					},
				},
			},
		},
	}

	// Marshal the body into JSON
	body, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("Error marshalling request body: %v", err)
		return "", err
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return "", err
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making the request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return "", err
	}

	// Check for a successful status code
	if resp.StatusCode == http.StatusOK {
		// Unmarshal the response into the structured format
		var response Response
		if err := json.Unmarshal(respBody, &response); err != nil {
			log.Printf("Error unmarshalling response body: %v", err)
			return "", err
		}

		// Extract the text from the response (first candidate)
		if len(response.Candidates) > 0 && len(response.Candidates[0].Content.Parts) > 0 {
			return response.Candidates[0].Content.Parts[0].Text, nil
		}

		// If no text found in the response, return a default message
		return "Не удалось сгенерировать ответ", nil
	}

	// Handle error status code
	return "", fmt.Errorf("error: %s", resp.Status)
}
