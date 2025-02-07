package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"

	"github.com/joho/godotenv"
)

type Issue struct {
	IDReadable   string        `json:"idReadable"`
	Updated      int64         `json:"updated"`
	Summary      string        `json:"summary"`
	Created      int64         `json:"created"`
	CustomFields []CustomField `json:"customFields"`
	Description  string        `json:"description"`
	Type         string        `json:"$type"`
}

type CustomField struct {
	Values Value  `json:"value"`
	Name   string `json:"name"`
	Type   string `json:"$type"`
}

type Value struct {
	Text   string  `json:"text"`
	Name   string  `json:"name"`
	Type   string  `json:"$type"`
	Number float64 `json:"number"`
	Multi  []Value `json:"multi"`
}

func (v *Value) UnmarshalJSON(data []byte) error {
	// fmt.Printf("\n%v\n", string(data))

	switch data[0] {
	case '{': // one object -> normal unmarshal
		re := regexp.MustCompile(`"text":"([^"\\]*(?:\\.[^"\\]*)*)"`)
		matches := re.FindSubmatch(data)
		if len(matches) > 1 {
			// fmt.Printf("FOUND: %s\n", matches[1])
			v.Text = fmt.Sprintf("%s", matches[1])
		}

		re = regexp.MustCompile(`"name":"([^"\\]*(?:\\.[^"\\]*)*)"`)
		matches = re.FindSubmatch(data)
		if len(matches) > 1 {
			// fmt.Printf("FOUND: %s\n", matches[1])
			v.Name = fmt.Sprintf("%s", matches[1])
		}

		re = regexp.MustCompile(`"\$type":"([^"\\]*(?:\\.[^"\\]*)*)"`)
		matches = re.FindSubmatch(data)
		if len(matches) > 1 {
			// fmt.Printf("FOUND: %s\n", matches[1])
			v.Type = fmt.Sprintf("%s", matches[1])
		}
	case '[': // multiple values -> ignore
		// fmt.Printf("MULTI: %s\n", data)
		v.Multi = []Value{}
		re := regexp.MustCompile(`\{.*?\}`)
		matches := re.FindAll(data, -1)
		for _, match := range matches {
			// for i, match := range matches {
			// fmt.Printf("FOUND %d: %s\n", i, match)
			var multiVal Value
			err := json.Unmarshal(match, &multiVal)
			if err != nil {
				return err
			}
			v.Multi = append(v.Multi, multiVal)
		}
	default: // try raw number
		strData := string(data)
		if strData != "null" {
			number, err := strconv.ParseFloat(strData, 64)
			if err != nil {
				return err
			}
			v.Number = number
		}
	}
	fmt.Printf("%+v\n", *v)
	return nil
}

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
	var issues []Issue
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

	// Remove any urls and email addresses

}

// removeURLsAndEmails entfernt URLs und E-Mail-Adressen aus dem gegebenen Text
func removeURLsAndEmails(text string) string {
	// Regex für URLs
	urlRegex := regexp.MustCompile(`https?://[^\s]+`)
	// Regex für E-Mail-Adressen
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)

	// URLs entfernen
	text = urlRegex.ReplaceAllString(text, "")
	// E-Mail-Adressen entfernen
	text = emailRegex.ReplaceAllString(text, "")

	return text
}
