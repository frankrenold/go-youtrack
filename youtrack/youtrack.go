package youtrack

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
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
	// fmt.Printf("%+v\n", *v)
	return nil
}

func (i *Issue) GetStringByCustomfieldName(fieldName string) string {
	for _, cf := range i.CustomFields {
		if cf.Name == fieldName {
			// fmt.Println(cf.Values.Name)
			if len(cf.Values.Text) > 0 {
				return cf.Values.Text
			} else {
				return cf.Values.Name
			}
		}
	}
	return ""
}

func (i *Issue) GetFloatByCustomfieldName(fieldName string) float64 {
	for _, cf := range i.CustomFields {
		if cf.Name == fieldName && cf.Values.Number != 0 {
			return cf.Values.Number
		}
	}
	return 0
}

func (i *Issue) GetSprints() []string {
	var sprints []string
	for _, cf := range i.CustomFields {
		if cf.Name == "Sprints" && len(cf.Values.Multi) > 0 {
			for _, s := range cf.Values.Multi {
				sprints = append(sprints, s.Name)
			}
		}
	}
	return sprints
}
