package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/frankrenold/go-youtrack/simplified"
	"github.com/frankrenold/go-youtrack/youtrack"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Define the API URL
	apiURL := fmt.Sprintf("https://%s/api/issues?query=%s&customFields=User%%20Story&customFields=Ready%%20if&customFields=Acceptance%%20Criteria%%20(Done%%20if)&customFields=Story%%20Points&customFields=Sprints&customFields=Priority&fields=idReadable,created,updated,summary,description,customFields(name,value(text,name))&$top=%v", os.Getenv("YT_DOMAIN"), url.QueryEscape(os.Getenv("YT_SEARCH_QUERY")), os.Getenv("MAX_ISSUES"))
	// Create client
	client := &http.Client{}
	// Create request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	// Add needed headers
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("YT_API_TOKEN")))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	// Run request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
	// Read the response body into expected struct
	var issues []youtrack.Issue
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}
	os.WriteFile("raw-response.json", body, 0644)
	err = json.Unmarshal(body, &issues)
	if err != nil {
		fmt.Println("Error processing response json:", err)
		return
	}

	// simplify
	var simpleIssues []simplified.Issue
	for _, issue := range issues {
		si := simplified.NewIssue(issue)
		simpleIssues = append(simpleIssues, *si)
	}
	simplJson, err := json.Marshal(simpleIssues)
	if err != nil {
		fmt.Println("Error converting simple issues to json:", err)
		return
	}
	os.WriteFile("issues.json", simplJson, 0644)
}
