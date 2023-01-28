package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/term"
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
		Features    []string `json:"features"`
	}

	type Emojis []struct {
		Name          string        `json:"name"`
		Roles         []interface{} `json:"roles"`
		ID            string        `json:"id"`
		RequireColons bool          `json:"require_colons"`
		Managed       bool          `json:"managed"`
		Animated      bool          `json:"animated"`
		Available     bool          `json:"available"`
	}

	fmt.Print("Enter your discord token: ")
	hidden, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil{
		fmt.Print(err)
	}
	fmt.Print(string(hidden))
	token := string(hidden)
	base_url := "https://discord.com/api/v10"
	guildsurl := base_url + "/users/@me/guilds"

	body := getrequest(guildsurl, token)
	var serverlist Server
	json.Unmarshal(body, &serverlist)
	var result []map[string]string

	for _, server := range serverlist {
		emojiurl := fmt.Sprintf(base_url+"/guilds/%v/emojis", server.ID)
		body := getrequest(emojiurl, token)
		var emojilist Emojis
		err := json.Unmarshal(body, &emojilist)
		if err != nil {
			fmt.Print("error in json marshall" + server.ID)
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
	}
	jsonstr, _ := json.Marshal(result)

	f, _ := os.Create("data.json")

	defer f.Close()

	f.WriteString(string(jsonstr))

}
