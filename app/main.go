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
	blockInfos, err := notionAPI.ReadRootPageBlocks(rootPageId)
	if err != nil {
		os.Exit(1)
	}

	filteredBlocks, err := notionAPI.FilterBlocks(blockInfos)
	if err != nil {
		os.Exit(1)
	}

	var content string

	//filteredBlocksが空の場合は
	if len(filteredBlocks) == 0 {
		content = "公開中の記事はありません"
	} else {
		var titlesAndUrls []string
		for _, block := range filteredBlocks {
			titlesAndUrls = append(titlesAndUrls, fmt.Sprintf("Title: %s, URL: %s", block.Title, block.URL))
		}

		content = "公開中の記事:\n" + strings.Join(titlesAndUrls, "\n")
	}

	dw := NewDiscordWebhook("NotificationPublicArticles", "", content, nil, false)

	whURL := loadEnv("DISCORD_WEBHOOK_URL")

	if err := dw.SendWebhook(whURL); err != nil {
		log.Fatal(err)
	}
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
