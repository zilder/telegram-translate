package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"log"
	"net/url"
	"strconv"
)

var (
	yandex_key		string
	telegram_key	string
)

type Update struct {
	UpdateId	int64	`json:"update_id"`
	Message		Message `json:"message"`
}

type Message struct {
	MessageId	int64	`json:"message_id"`
	Text		string	`json:"text"`
	Chat		Chat	`json:"chat"`
}


type Chat struct {
	Id			int64	`json:"id"`
	FirstName	string	`json:"first_name"`
	LastName	string	`json:"last_name"`
}

type Translation struct {
	Text		[]string	`json:"text"`
	Lang		string		`json:"lang"`
}

func translate(text string) []string {
	// make query to translation service
	values := url.Values{}
	values.Add("key", yandex_key)
	values.Add("lang", "en-ru")
	values.Add("text", text)
	query := fmt.Sprintf("https://translate.yandex.net/api/v1.5/tr.json/translate?%s", values.Encode())
	resp, err := http.Get(query)
	if err != nil {
		log.Println("HttpRequest error", err)
	}
	defer resp.Body.Close()

	// decoding response
	var translation Translation
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&translation)
	log.Println("translation", translation)

	return translation.Text
}

func onMessage(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var updateInfo Update
	err := decoder.Decode(&updateInfo)
	if err != nil {
		//panic()
		log.Fatal()
	}
	log.Println(updateInfo)
	variants := translate(updateInfo.Message.Text)
	log.Println(variants)
	sendMessage(updateInfo.Message.Chat, variants[0])
}

func sendMessage(chat Chat, text string) {
	values := url.Values{}
	values.Add("chat_id", strconv.FormatInt(chat.Id, 10))
	values.Add("text", text)
	query := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?%s",
		telegram_key, values.Encode())
	log.Println(query)
	http.Get(query)
	//defer resp.Body.Close()
}


func main() {
	fmt.Println("starting server...")

	yandex_key   = "trnsl.1.1.20150914T160436Z.5a6abb6f3a3807df.c3ec0be6bb41a554ed08d86b91cc2d8d07178eb5"
	telegram_key = "131973765:AAG7qZTqwpes8iZ2ET-NMdTYL4qXm7y7Lr8"

	http.HandleFunc("/onMessage", onMessage)
	http.ListenAndServe(":8080", nil)
}

