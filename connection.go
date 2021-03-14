package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func setupConnection() (c *websocket.Conn, err error){
	token, err := getToken()
	if err != nil {
		log.Fatal(err)
	}

	//TODO Will be used to close the connection later
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	//link := "tmp.mama.sh"
	//u := url.URL{Scheme: "wss", Host: link, Path: "/api/game"}

	u := "ws://localhost:8056/join?peer=morten"
	//u := "ws://localhost:8056/create"
	//u := url.URL{Scheme: "ws", Host: link, Path: "/create"}
	log.Printf("connecting to %s", u)

	header := http.Header{}
	header.Add("Authorization", "bearer "+token)

	c, _, err = websocket.DefaultDialer.Dial(u, header)
	return
}
