package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

func main() {
	rootPageId := loadEnv("NOTION_ROOT_PAGE_ID")

	notionAPI := &NotionAPI{
		APIKey: loadEnv("NOTION_API_KEY"),
	}
	blockInfo, err := notionAPI.ReadRootPageBlocks(rootPageId)
	if err != nil {
		log.Printf("ReadRootPageBlocks error: %v\n", err)
	}

	publicPages, err := notionAPI.FilterPublicBlocks(blockInfo)
	if err != nil {
		log.Printf("FilterPublicBlocks error: %v\n", err)
	}

	var content string

	//filteredBlocksが空の場合は
	if len(publicPages) == 0 {
		content = "公開中の記事はありません"
	} else {
		var titlesAndUrls []string
		for _, block := range publicPages {
			titlesAndUrls = append(titlesAndUrls, fmt.Sprintf("Title: %s, URL: %s", block.Title, block.URL))
		}

		mention := "<@&1052152206750662656>"
		content = mention + "公開中の記事:\n" + strings.Join(titlesAndUrls, "\n")
	}

	whURL := loadEnv("DISCORD_WEBHOOK_URL")

	// Split content into multiple messages if it exceeds 1900 characters
	contentParts := splitContent(content, 1900)

	for _, contentPart := range contentParts {
		dw := NewDiscordWebhook("NotificationPublicArticles", "", contentPart, nil, false)

		if err := dw.SendWebhook(whURL); err != nil {
			//Discordに送信できなかった場合は、エラーを出力します。
			log.Printf("SendWebhook error: %v\n", err)
		}
	}
}

func splitContent(content string, maxChars int) []string {
	var contentParts []string
	contentLength := len(content)

	if contentLength <= maxChars {
		return []string{content}
	}

	// Split content into parts
	for start := 0; start < contentLength; start += maxChars {
		end := start + maxChars
		if end > contentLength {
			end = contentLength
		}
		contentParts = append(contentParts, content[start:end])
	}

	return contentParts
}

func loadEnv(keyName string) string {
	err := godotenv.Load("../.env")
	// もし err がnilではないなら、"読み込み出来ませんでした"が出力されます。
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	}
	// .envの SAMPLE_MESSAGEを取得して、messageに代入します。
	message := os.Getenv(keyName)

	return message
}
