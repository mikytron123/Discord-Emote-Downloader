package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	//	"golang.org/x/term"
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

	type StickerPackList struct {
		StickerPacks []struct {
			ID       string `json:"id"`
			Stickers StickerList `json:"stickers"`
			Name string `json:"name"`
		} `json:"sticker_packs"`
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	//hidden, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Print(err)
	}
	//fmt.Print(string(hidden))
	token := string(os.Getenv("AUTH"))
	base_url := "https://discord.com/api/v10"
	var result []map[string]string

	stickerpackurl := fmt.Sprintf(base_url + "/sticker-packs")
	stickerpackbody := getrequest(stickerpackurl, token)

	var stickerpacklist StickerPackList
	error := json.Unmarshal(stickerpackbody, &stickerpacklist)
	if error != nil {
		fmt.Print("error in json marshall stickerpack list ")
		fmt.Print(error)
	}
	for _, stickerpack := range stickerpacklist.StickerPacks {
		for _, sticker := range stickerpack.Stickers {
			if sticker.FormatType == 3 {
				name := stickerpack.Name + " " + sticker.Name
				url := "https://cdn.discordapp.com/stickers/" + sticker.ID + ".json"
				result = append(result, map[string]string{"name": name, "url": url,"tags":sticker.Tags})
			
		    }else if sticker.FormatType == 2 {
                 name := stickerpack.Name + " " + sticker.Name
				 url := "https://cdn.discordapp.com/stickers/" + sticker.ID + ".png"
				 result = append(result,map[string]string{"name":name,"url":url,"tags":sticker.Tags})
           
			}

		}
	}

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
				tempstickerresult = append(tempstickerresult, tempmap)
				continue
			}
			//}else if sticker.FormatType == 4 {
			if sticker.FormatType == 2 {
				tempmap["name"] = sticker.Name
				tempmap["url"] = "https://cdn.discordapp.com/stickers/" + sticker.ID + ".png"
				tempmap["tags"] = sticker.Tags
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
