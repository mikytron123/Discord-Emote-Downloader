package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func getrequest(url string, token string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("error in new request")
	}
	req.Header.Add("Authorization", token)
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("error in do request")
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Print("error in readall")
	}
	return body

}

func main() {

	type Server []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Icon        string `json:"icon"`
		Owner       bool   `json:"owner"`
		Permissions string `json:"permissions"`
	}

	type Emojis []struct {
		Name          string `json:"name"`
		ID            string `json:"id"`
		RequireColons bool   `json:"require_colons"`
		Managed       bool   `json:"managed"`
		Animated      bool   `json:"animated"`
		Available     bool   `json:"available"`
	}

	type StickerList []struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		FormatType int    `json:"format_type"`
		Tags       string `json:"tags"`
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := string(os.Getenv("AUTH"))
	baseUrl := "https://discord.com/api/v10"
	var result []map[string]string

	guildsurl := baseUrl + "/users/@me/guilds"

	body := getrequest(guildsurl, token)
	var serverList Server
	json.Unmarshal(body, &serverList)

	emojiBaseUrl := "https://cdn.discordapp.com/emojis/"
	stickerBaseUrl := "https://cdn.discordapp.com/stickers/"
	for _, server := range serverList {
		emojiurl := fmt.Sprintf("%v/guilds/%v/emojis", baseUrl, server.ID)
		time.Sleep(1 * time.Second)
		body := getrequest(emojiurl, token)
		var emojiList Emojis
		err = json.Unmarshal(body, &emojiList)
		if err != nil {
			fmt.Print("error in json marshall emote" + server.ID)
		}

		for _, emote := range emojiList {
			tempMap := make(map[string]string)
			var ext string

			if emote.Animated {
				ext = ".gif"
			} else {
				ext = ".png"
			}
			tempMap["name"] = emote.Name
			tempMap["url"] = fmt.Sprintf("%v%v%v", emojiBaseUrl, emote.ID, ext)
			tempMap["tags"] = ""
			tempMap["type"] = "emote"
			result = append(result, tempMap)
		}

		stickerUrl := fmt.Sprintf("%v/guilds/%v/stickers", baseUrl, server.ID)
		time.Sleep(1 * time.Second)
		stickerBody := getrequest(stickerUrl, token)

		var stickerList StickerList
		err = json.Unmarshal(stickerBody, &stickerList)

		if err != nil {
			fmt.Print("error in json marshall sticker " + server.ID)
			fmt.Print(err)
			fmt.Print(string(stickerBody))
			break
		}

		for _, sticker := range stickerList {
			tempMap := make(map[string]string)
			var ext string
			var stickerType string
			if sticker.FormatType == 1 {
				ext = ".png"
				stickerType = "sticker"
			} else if sticker.FormatType == 2 {
				ext = ".png"
				stickerType = "apng"
			} else if sticker.FormatType == 4 {
				ext = ".gif"
				stickerType = "sticker"
			} else {
				continue
			}

			tempMap["name"] = sticker.Name
			tempMap["url"] = fmt.Sprintf("%v%v%v", stickerBaseUrl, sticker.ID, ext)
			tempMap["tags"] = sticker.Tags
			tempMap["type"] = stickerType
			result = append(result, tempMap)
		}

	}

	jsonstr, _ := json.Marshal(result)

	f, _ := os.Create("data.json")

	defer f.Close()

	f.WriteString(string(jsonstr))

}
