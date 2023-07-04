package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ReplaceLink replaces the link from www.notion.so to midra-lab.notion.site
func ReplaceLink(link string) string {
	newLink := strings.Replace(link, "www.notion.so", "midra-lab.notion.site", 1)
	return newLink
}

// RemoveNotionLink removes the link from www.notion.so
func RemoveNotionLink(link string) string {
	newLink := strings.Replace(link, "https://www.notion.so/", "", 1)
	return newLink
}

// GenerateNotionPageURL creates the URL of a notion page
func GenerateNotionPageURL(pageId string) string {
	blockIDWithoutHyphens := strings.ReplaceAll(pageId, "-", "")
	return fmt.Sprintf("https://www.notion.so/%s", blockIDWithoutHyphens)
}

// GenerateCustomDomainNotionURL creates the URL of a notion page with custom domain
func GenerateCustomDomainNotionURL(pageId string) string {
	blockURL := GenerateNotionPageURL(pageId)
	return ReplaceLink(blockURL)
}

func CheckLink(url string) (bool, error) {
	client := &http.Client{}
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return false, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

// GetDatabaseTitle returns the title of a database
func (n *NotionAPI) GetDatabaseTitle(databaseId string) (string, error) {
	dbUrl := "https://api.notion.com/v1/databases/" + databaseId

	req, err := http.NewRequest("GET", dbUrl, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+n.APIKey)
	req.Header.Add("Notion-Version", "2022-06-28")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var notionRes struct {
		Title []struct {
			Text struct {
				Content string `json:"content"`
			} `json:"text"`
		} `json:"title"`
	}
	if err := json.NewDecoder(res.Body).Decode(&notionRes); err != nil {
		return "", err
	}

	if len(notionRes.Title) == 0 {
		return "", fmt.Errorf("No title found")
	}

	return notionRes.Title[0].Text.Content, nil
}
