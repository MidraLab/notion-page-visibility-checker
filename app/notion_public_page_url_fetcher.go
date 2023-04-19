package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type NotionAPI struct {
	APIKey string
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

type block struct {
	Object string `json:"object"`
	ID     string `json:"id"`
	Type   string `json:"type"`
}

type blocksResponse struct {
	Object  string  `json:"object"`
	Results []block `json:"results"`
}

type blockInfo struct {
	Type  string
	ID    string
	Title string
	URL   string
}

func (n *NotionAPI) ReadPageID(databaseId string) ([]string, []string, error) {
	dbUrl := "https://api.notion.com/v1/databases/" + databaseId + "/query"

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

func (n *NotionAPI) ReadRootPageBlocks(rootPageId string) ([]blockInfo, error) {
	var blockInfos []blockInfo

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.notion.com/v1/blocks/%s/children", rootPageId), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+n.APIKey)
	req.Header.Set("Notion-Version", "2022-06-28")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error: %s", string(body))
	}

	var blocksResponse blocksResponse
	err = json.Unmarshal(body, &blocksResponse)
	if err != nil {
		return nil, err
	}

	for _, block := range blocksResponse.Results {
		blockInfo := blockInfo{
			Type: block.Type,
			ID:   block.ID,
		}

		if block.Type == "child_database" || block.Type == "child_page" {
			blockIDWithoutHyphens := strings.ReplaceAll(block.ID, "-", "")
			blockURL := fmt.Sprintf("https://www.notion.so/%s", blockIDWithoutHyphens)
			blockInfo.URL = ReplaceLink(blockURL)

			// Get the title for child_database and child_page blocks
			title, err := n.GetDatabaseTitle(block.ID)
			if err != nil {
				return nil, err
			}
			blockInfo.Title = title
		}

		blockInfos = append(blockInfos, blockInfo)
	}

	return blockInfos, nil
}

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

func (n *NotionAPI) FilterBlocks(blocks []blockInfo) ([]blockInfo, error) {
	var filteredBlocks []blockInfo
	for _, block := range blocks {

		if block.Type == "child_database" {
			isValid, err := CheckLink(ReplaceLink(block.URL))
			if err != nil {
				return nil, err
			}

			if isValid {
				filteredBlocks = append(filteredBlocks, blockInfo{
					Type:  block.Type,
					Title: block.Title,
					URL:   ReplaceLink(block.URL),
				})
			}

			urls, titles, err := n.ReadPageID(block.ID)
			if err != nil {
				return nil, err
			}

			for i, url := range urls {
				isValid, err := CheckLink(ReplaceLink(url))
				if err != nil {
					return nil, err
				}
				if isValid {
					filteredBlocks = append(filteredBlocks, blockInfo{
						Type:  block.Type,
						Title: titles[i],
						URL:   ReplaceLink(url),
					})
				}
			}

		}
	}

	return filteredBlocks, nil
}

func ReplaceLink(link string) string {
	newLink := strings.Replace(link, "www.notion.so", "midra-lab.notion.site", 1)
	return newLink
}

func CheckLink(url string) (bool, error) {
	resp, err := http.Get(url)
	if err != nil {
		return false, nil
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}
