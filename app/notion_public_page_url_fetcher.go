package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type NotionAPI struct {
	DatabaseID string
	APIKey     string
}

type notionResponse struct {
	Results []struct {
		URL      string `json:"url"`
		Property struct {
			Title struct {
				Title []struct {
					PlainText string `json:"plain_text"`
				} `json:"title"`
			} `json:"名前"`
		} `json:"properties"`
	} `json:"Results"`
}

func (n *NotionAPI) ReadPageID() ([]string, []string, error) {
	dbUrl := "https://api.notion.com/v1/databases/" + n.DatabaseID + "/query"

	payload := strings.NewReader(`{
    "page_size": 100
}`)

	req, err := http.NewRequest("POST", dbUrl, payload)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+n.APIKey)
	req.Header.Add("Notion-Version", "2022-06-28")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	var notionRes notionResponse
	if err := json.NewDecoder(res.Body).Decode(&notionRes); err != nil {
		return nil, nil, err
	}

	if len(notionRes.Results) == 0 {
		return nil, nil, fmt.Errorf("No results")
	}

	var urls []string
	var titles []string
	for _, result := range notionRes.Results {
		urls = append(urls, result.URL)
		titles = append(titles, result.Property.Title.Title[0].PlainText)
	}

	return urls, titles, nil
}
