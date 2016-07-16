package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

type ToneAnalyzer struct {
	Creds string
}

func (a ToneAnalyzer) Analyze(text string) (string, error) {
	body := strings.NewReader(`{"text":"` + text + `"}`)
	url := fmt.Sprintf("https://%s@gateway.watsonplatform.net/tone-analyzer/api/v3/tone?version=2016-05-18", a.Creds)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		return "", errors.New(string(data))
	}

	type category struct {
		Tones Tones  `json:"tones"`
		ID    string `json:"category_id"`
		Name  string `json:"category_name"`
	}

	var tones struct {
		Document struct {
			Categories []category `json:"tone_categories"`
		} `json:"document_tone"`
	}

	if err := json.Unmarshal(data, &tones); err != nil {
		return "", err
	}

	var c *category
	for i := range tones.Document.Categories {
		if tones.Document.Categories[i].ID == "emotion_tone" {
			c = &tones.Document.Categories[i]
			break
		}
	}

	if c == nil {
		return "", errors.New("category emotion_tone not returned")
	}

	sort.Sort(sort.Reverse(c.Tones))

	if len(c.Tones) == 0 {
		return "", errors.New("no tones returned")
	}

	return c.Tones[0].ID, nil
}

type Tone struct {
	Score float64 `json:"score"`
	ID    string  `json:"tone_id"`
	Name  string  `json:"tone_name"`
}

type Tones []Tone

func (ts Tones) Len() int           { return len(ts) }
func (ts Tones) Less(i, j int) bool { return ts[i].Score < ts[j].Score }
func (ts Tones) Swap(i, j int)      { ts[i], ts[j] = ts[j], ts[i] }
