package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"strings"
)

func main() {
	notionAPI := NotionAPI{
		DatabaseID: loadEnv("NOTION_SURVEY_DATABASE_ID"),
		APIKey:     loadEnv("NOTION_API_KEY"),
	}

	urls, titles, err := notionAPI.ReadPageID()
	if err != nil {
		fmt.Println(err)
	}

	validUrls := []string{}
	validTitles := []string{}

	for i, url := range urls {
		newLink := replaceLink(url)
		isValid, validLink, validTitle, err := checkLink(newLink, titles[i])

		if err != nil {
			fmt.Println(err)
			continue
		}

		if isValid {
			validUrls = append(validUrls, validLink)
			validTitles = append(validTitles, validTitle)
		}
	}

	fmt.Println(validUrls)
	fmt.Println(validTitles)
}

func replaceLink(link string) string {
	newLink := strings.Replace(link, "www.notion.so", "midra-lab.notion.site", 1)
	//fmt.Print(newLink + "\n")
	return newLink
}

func checkLink(link string, title string) (bool, string, string, error) {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return false, "", "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, "", "", err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return false, "", "", nil
	}

	if res.StatusCode == http.StatusOK {
		return true, link, title, nil
	}

	return false, "", "", fmt.Errorf("unexpected status code: %d", res.StatusCode)
}

func loadEnv(keyName string) string {
	err := godotenv.Load(".env")
	// もし err がnilではないなら、"読み込み出来ませんでした"が出力されます。
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	}
	// .envの SAMPLE_MESSAGEを取得して、messageに代入します。
	message := os.Getenv(keyName)

	return message
}
