package models

import "fmt"

type Bots struct {
	list       []Bot
	count      int
	LastGenome []int
}

var iteration = 0

var moves = 0

var MainBots = Bots{}

func CreateBots(count int) Bots {
	return Bots{count: count}
}

func (e *Bots) Get() []Bot {
	return e.list
}
func (e *Bots) Set(list []Bot) {
	e.list = list
}

func (e *Bots) FillBot() {
	var list []Bot
	for i := 0; i < e.count; i++ {
		bot := NewBot(nil)
		bot.GenerateGenome()
		list = append(list, bot)
	}
	e.list = list
}

func (e *Bots) Updater() {
	for true {
		//time.Sleep(time.Second / 10)

		//var wg sync.WaitGroup

		//wg.Add(len(e.list))
		for k := range e.list {
			//go func(k int) {
			e.list[k].Action()
			//defer wg.Done()
			//}(k)
		}
		//wg.Wait()

		newBots := make([]Bot, 0)
		//wg.Add(len(e.list))
		for k := range e.list {
			//go func(k int) {
			if e.list[k].Health >= 0 && e.list[k].Age < 1000 {
				newBots = append(newBots, e.list[k])
			}
			//defer wg.Done()
			//}(k)
		}
		//wg.Wait()
		if len(newBots) == 0 {
			newBots = append(newBots, e.list[0])
		}

		//fmt.Println("")

		if len(newBots) <= 2 {
			newWaveBots := make([]Bot, 0)
			for i, bot := range newBots {
				if i == 0 {
					MainBots.LastGenome = append(make([]int, 0), bot.Genome...)
				}
				for i := 0; i < e.count/len(newBots); i++ {
					newBot := NewBot(nil)
					newBot.Genome = make([]int, 0)
					newBot.Genome = append(newBot.Genome, bot.Genome...)

					newBot.ActualGenome = make([]int, 0)
					newBot.ActualGenome = append(bot.ActualGenome)
					newBot.Mutate()
					newWaveBots = append(newWaveBots, newBot)

					if len(newBot.Genome) == 0 {
						//fmt.Println("test")
					}

				}
				bot.Reset()
				newWaveBots = append(newWaveBots, bot)
				if len(bot.Genome) == 0 {
					//fmt.Println("BOT test")
				}
			}
			newBots = newWaveBots
			fmt.Println("iteration ", iteration, " moves", moves)
			iteration++
			moves = 0

		}
		moves++
		//for i, bot := range e.list {
		//fmt.Println(i, bot.ID(), bot.Genome)
		//}

		e.list = newBots

		//for i, bot := range e.list {
		//fmt.Println(i, bot.ID(), bot.Genome)
		//}
		MainEats.Update()
	}
}
