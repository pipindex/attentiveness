package main

import (
	"attentiveness"
	"attentiveness/models"
	"encoding/json"
	"log"
	"time"
)

type Message struct {
	Message string `json:"message"`
}

var (
	// firebaseURL              = Models.FirebaseURL{"https://glaring-heat-9244.firebaseio.com/", "ug40VoPebrjwaHzHmaKqKVCQpjfPF7G83aUmY1RN"}
	firebaseURL         = Models.FirebaseURL{"https://pho-chat-dev.firebaseio.com/", "B7ioKnRoaf6ORCU4p8qUeAETTVcqBVykYIbeSEpF"}
	timeNowMilliseconds = int(time.Now().UnixNano() / int64(time.Millisecond)) // int(time.Now().Unix()*1000)//1463988502284//*1000
)

func main() {
	log.Println("Timestamp:", timeNowMilliseconds)

	activeWebinars := attentiveness.GetActiveWebinars(firebaseURL)

    log.Println(json.Marshal(activeWebinars))

	response := attentiveness.CalculateAverageAttentivenessForActiveWebinars(firebaseURL, activeWebinars, timeNowMilliseconds)

	responseString, _ := json.Marshal(response)

	log.Println("Response JSON:", string(responseString[:]))
}


