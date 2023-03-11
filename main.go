package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
    "time"
	"log"
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
	if err != nil{
		fmt.Print("error in readall")
	}
	return body

}

func main() {

	type Server []struct {
		ID          string   `json:"id"`
		Name        string   `json:"name"`
		Icon        string   `json:"icon"`
		Owner       bool     `json:"owner"`
		Permissions string   `json:"permissions"`
	}

	type Emojis []struct {
		Name          string        `json:"name"`
		ID            string        `json:"id"`
		RequireColons bool          `json:"require_colons"`
		Managed       bool          `json:"managed"`
		Animated      bool          `json:"animated"`
		Available     bool          `json:"available"`
	}

	type Stickers []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		FormatType  int    `json:"format_type"`
	}

	err := godotenv.Load()
    if err != nil {
		log.Fatal("Error loading .env file")
    }
	//hidden, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil{
		fmt.Print(err)
	}
	//fmt.Print(string(hidden))
	token := string(os.Getenv("AUTH"))
	base_url := "https://discord.com/api/v10"
	guildsurl := base_url + "/users/@me/guilds"

	body := getrequest(guildsurl, token)
	var serverlist Server
	json.Unmarshal(body, &serverlist)
	var result []map[string]string

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
			tempresult = append(tempresult, tempmap)
		}
		result = append(result, tempresult...)

		stickerurl := fmt.Sprintf(base_url+"/guilds/%v/stickers", server.ID)
		time.Sleep(1 * time.Second)
		stickerbody := getrequest(stickerurl, token)

		var stickerlist Stickers
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
