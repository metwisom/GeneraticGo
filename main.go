package main

var MainWorld World
var MainBots Bots

func main() {

	MainWorld = CreateWorld(70, 70)
	MainWorld.EatPercent = 0.00
	MainWorld.EatBader = 0.005
	MainWorld.FillEat()

	MainBots = CreateBots(64)
	MainBots.FillBot()

	go startServ()

	MainBots.Updater()

}
