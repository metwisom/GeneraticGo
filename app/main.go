package main

import (
	"encoding/json"
	"fmt"
	"main/models"
	"main/web"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Print("\033[2J") // устанавливаем курсор в начало первой строки

	models.MainWorld = models.CreateWorld(200, 200)
	models.MainBots = models.CreateBots(24)

	//models.MainEats.Damage = 1
	//models.MainEats.Growth = 35
	//models.MainEats.Rot = 100
	models.MainEats.Percent = 0.01
	//models.MainEats.FillEat()

	dat, _ := os.ReadFile("eats.txt")
	var inventory []models.Eat
	if err := json.Unmarshal([]byte(dat), &inventory); err != nil {
		models.MainEats.FillEat()
	} else {
		models.MainEats.Set(inventory)
	}

	dat, _ = os.ReadFile("bots.txt")
	var savedBots []models.Bot
	if err := json.Unmarshal([]byte(dat), &savedBots); err != nil {
		models.MainBots.FillBot()
	} else {
		models.MainBots.Set(savedBots)
	}

	go web.StartServ()

	go models.MainBots.Updater()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	// Ожидаем получение сигнала
	<-signalChannel
	//
	fmt.Println("Received signal, shutting down.")

	//file, _ := os.Create("hello.txt")
	//j, _ := json.Marshal(models.MainBots.Get())
	//fmt.Println(models.MainBots.LastGenome)
	//file.WriteString(string(j))
	//file.Close()

	file1, _ := os.Create("bots.txt")
	j1, _ := json.Marshal(models.MainBots.Get())
	file1.WriteString(string(j1))
	file1.Close()

	file2, _ := os.Create("eats.txt")
	j2, _ := json.Marshal(models.MainEats.Get())
	file2.WriteString(string(j2))
	file2.Close()

}
