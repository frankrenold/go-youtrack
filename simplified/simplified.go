package simplified

import (
	"regexp"

	"github.com/frankrenold/go-youtrack/youtrack"
)

type Issue struct {
	Id                 string   `json:"id"`
	Created            int64    `json:"created"`
	Updated            int64    `json:"updated"`
	Type               string   `json:"type"`
	Priority           string   `json:"priority"`
	Summary            string   `json:"summary"`
	Description        string   `json:"description"`
	UserStory          string   `json:"userstory"`
	ReadyIf            string   `json:"ready-if"`
	AcceptanceCriteria string   `json:"acceptance-criteria"`
	StoryPoints        float64  `json:"storypoints"`
	Sprints            []string `json:"sprints"`
}

func NewIssue(yt youtrack.Issue) *Issue {
	i := Issue{
		Id:                 yt.IDReadable,
		Created:            yt.Created,
		Updated:            yt.Updated,
		Type:               yt.Type,
		Priority:           yt.GetStringByCustomfieldName("Priority"),
		Summary:            removeURLsAndEmails(yt.Summary),
		Description:        removeURLsAndEmails(yt.Description),
		UserStory:          removeURLsAndEmails(yt.GetStringByCustomfieldName("User Story")),
		ReadyIf:            removeURLsAndEmails(yt.GetStringByCustomfieldName("Ready if")),
		AcceptanceCriteria: removeURLsAndEmails(yt.GetStringByCustomfieldName("Acceptance Criteria (Done if)")),
		StoryPoints:        yt.GetFloatByCustomfieldName("Story Points"),
		Sprints:            yt.GetSprints(),
	}
	return &i
}

// removeURLsAndEmails entfernt URLs und E-Mail-Adressen aus dem gegebenen Text
func removeURLsAndEmails(text string) string {
	// Regex für URLs
	urlRegex := regexp.MustCompile(`https?://[^\s]+`)
	// Regex für E-Mail-Adressen
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	// Regex für Mentions
	mentionRegex := regexp.MustCompile(`@[^\s]+`)

	// URLs entfernen
	text = urlRegex.ReplaceAllString(text, "")
	// E-Mail-Adressen entfernen
	text = emailRegex.ReplaceAllString(text, "")
	// Mentions entfernen
	text = mentionRegex.ReplaceAllString(text, "")

	return text
}
