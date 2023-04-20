package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DiscordWebhook struct {
	UserName  string         `json:"username"`
	AvatarURL string         `json:"avatar_url"`
	Content   string         `json:"content"`
	Embeds    []DiscordEmbed `json:"embeds"`
	TTS       bool           `json:"tts"`
}

type DiscordImg struct {
	URL string `json:"url"`
	H   int    `json:"height"`
	W   int    `json:"width"`
}

type DiscordAuthor struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Icon string `json:"icon_url"`
}

type DiscordField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type DiscordEmbed struct {
	Title  string         `json:"title"`
	Desc   string         `json:"description"`
	URL    string         `json:"url"`
	Color  int            `json:"color"`
	Image  DiscordImg     `json:"image"`
	Thum   DiscordImg     `json:"thumbnail"`
	Author DiscordAuthor  `json:"author"`
	Fields []DiscordField `json:"fields"`
}

func NewDiscordWebhook(userName, avatarURL, content string, embeds []DiscordEmbed, tts bool) *DiscordWebhook {
	return &DiscordWebhook{
		UserName:  userName,
		AvatarURL: avatarURL,
		Content:   content,
		Embeds:    embeds,
		TTS:       tts,
	}
}

func (dw *DiscordWebhook) AddEmbeds(embeds ...DiscordEmbed) {
	dw.Embeds = append(dw.Embeds, embeds...)
}

func (dw *DiscordWebhook) SendWebhook(whURL string) error {
	j, err := json.Marshal(dw)
	if err != nil {
		return fmt.Errorf("json err: %s", err.Error())
	}

	req, err := http.NewRequest("POST", whURL, bytes.NewBuffer(j))
	if err != nil {
		return fmt.Errorf("new request err: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("client err: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		fmt.Println("sent", dw) //成功
	} else {
		return fmt.Errorf("%#v\n", resp) //失敗
	}

	return nil
}
