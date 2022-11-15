package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func startServ() {
	http.HandleFunc("/", webs)

	err := http.ListenAndServe(":8765", nil)
	if err != nil {
		log.Fatal(err)
	}
}

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func webs(w http.ResponseWriter, req *http.Request) {

	upgrade.CheckOrigin = func(r *http.Request) bool { return true }

	w.Header().Set("Access-Control-Allow-Origin", "*")

	ws, err := upgrade.Upgrade(w, req, nil)
	if err != nil {
		log.Println(err)
	}

	reader(ws)
}

func reader(conn *websocket.Conn) {
	go sender(conn)
}

type Line struct {
	Type string        `json:"type"`
	Data []interface{} `json:"data"`
}

func sender(conn *websocket.Conn) {
	for true {
		err := sendBot(conn)
		if err != nil {
			conn.Close()
			return
		}
		err = sendMap(conn)
		if err != nil {
			conn.Close()
			return
		}

		time.Sleep(time.Second / 100)
	}
}

func sendMap(conn *websocket.Conn) error {
	var respBot Line
	respBot.Type = "map"

	for _, v := range MainWorld.GetAllMap() {

		if v.typo == "e" {

			if v.Value == 1 {
				slice2 := []interface{}{v.pos.X, v.pos.Y, 1}
				respBot.Data = append(respBot.Data, slice2)
			}
			if v.Value == 2 {
				slice2 := []interface{}{v.pos.X, v.pos.Y, 2}
				respBot.Data = append(respBot.Data, slice2)
			}
			if v.Value == 3 {
				slice2 := []interface{}{v.pos.X, v.pos.Y, 3}
				respBot.Data = append(respBot.Data, slice2)
			}
			if v.Value == 4 {
				slice2 := []interface{}{v.pos.X, v.pos.Y, 4}
				respBot.Data = append(respBot.Data, slice2)
			}
		} else {
			slice2 := []interface{}{v.pos.X, v.pos.Y, 5}
			respBot.Data = append(respBot.Data, slice2)
		}
	}

	jsonStr, err := json.Marshal(respBot)
	if err != nil {
		fmt.Println("err")
		return err
	}

	if err := conn.WriteMessage(1, jsonStr); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func sendBot(conn *websocket.Conn) error {
	var respMap Line
	respMap.Type = "bot"

	listMutex.RLock()
	for _, bot := range MainBots.pList {
		slice2 := []int{bot.pos.X, bot.pos.Y, bot.Health}
		respMap.Data = append(respMap.Data, slice2)
	}

	listMutex.RUnlock()
	jsonStr, err := json.Marshal(respMap)
	if err != nil {
		return err
	}

	if err := conn.WriteMessage(1, jsonStr); err != nil {
		return err
	}
	return nil
}
