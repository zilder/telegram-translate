package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"log"
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
	query := fmt.Sprintf("https://translate.yandex.net/api/v1.5/tr.json/translate?key=%s&lang=%s&text=%s",
		yandex_key, "en-ru", text)
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
	query := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%v&text=%s",
		telegram_key, chat.Id, text)
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

