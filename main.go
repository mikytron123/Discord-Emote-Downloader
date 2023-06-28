package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/joho/godotenv"
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

	if err != nil {
		fmt.Print(err)
	}

	token := string(os.Getenv("AUTH"))
	base_url := "https://discord.com/api/v10"
	var result []map[string]string

	guildsurl := base_url + "/users/@me/guilds"

	body := getrequest(guildsurl, token)
	var serverlist Server
	json.Unmarshal(body, &serverlist)

	for _, server := range serverlist {
		emojiurl := fmt.Sprintf(base_url+"/guilds/%v/emojis", server.ID)
		time.Sleep(1 * time.Second)
		body := getrequest(emojiurl, token)
		var emojilist Emojis
		err := json.Unmarshal(body, &emojilist)
		if err != nil {
			fmt.Print("error in json marshall emote" + server.ID)
		}
		var tempresult []map[string]string
		for _, emote := range emojilist {
			tempmap := make(map[string]string)
			var ext string
			tempmap["name"] = emote.Name
			if emote.Animated {
				ext = ".gif"
			} else {
				ext = ".png"
			}
			tempmap["url"] = "https://cdn.discordapp.com/emojis/" + emote.ID + ext
			tempmap["tags"] = ""
			tempmap["type"] = "emote"
			tempresult = append(tempresult, tempmap)
		}
		result = append(result, tempresult...)

		stickerurl := fmt.Sprintf(base_url+"/guilds/%v/stickers", server.ID)
		time.Sleep(1 * time.Second)
		stickerbody := getrequest(stickerurl, token)

		var stickerlist StickerList
		err2 := json.Unmarshal(stickerbody, &stickerlist)
		if err2 != nil {
			fmt.Print("error in json marshall sticker " + server.ID)
			fmt.Print(err2)
			fmt.Print(string(stickerbody))
			break
		}
		var tempstickerresult []map[string]string
		for _, sticker := range stickerlist {
			tempmap := make(map[string]string)
			if sticker.FormatType == 1 {
				tempmap["name"] = sticker.Name
				tempmap["url"] = "https://cdn.discordapp.com/stickers/" + sticker.ID + ".png"
				tempmap["tags"] = sticker.Tags
				tempmap["type"] = "sticker"
				tempstickerresult = append(tempstickerresult, tempmap)
				continue
			}
			if sticker.FormatType == 2 {
				tempmap["name"] = sticker.Name
				tempmap["url"] = "https://cdn.discordapp.com/stickers/" + sticker.ID + ".png"
				tempmap["tags"] = sticker.Tags
				tempmap["type"] = "apng" 
				tempstickerresult = append(tempstickerresult, tempmap)
			}

		}
		result = append(result, tempstickerresult...)

	}

	jsonstr, _ := json.Marshal(result)

	f, _ := os.Create("data.json")

	defer f.Close()

	f.WriteString(string(jsonstr))

}
