package main

import (
	"fmt"
	"sort"
	"sync"
)

type Bots struct {
	list      []Bot
	fList     []Bot
	pList     []Bot
	botCount  int
	iteration int
	cycles    int
}

var listMutex sync.RWMutex

func CreateBots(count int) Bots {
	return Bots{botCount: count, iteration: 1}
}

func (e *Bots) FillBot() {
	var list []Bot
	for i := 0; i < e.botCount; i++ {
		bot := NewBot(nil)
		//bot.Genome = []int{17,3,11,13,6,14,5,13,11,19,13,17,0,18,10,0,6,4,4,14,9,4,18,0,8,3,2,14,3,19,14,8,7,4,12,5,8,10,6,9,4,16,11,7,14,16,17,10,8,13,6,17,10,19,8,1,3,17,14,3,15,18,7,14}
		bot.GenerateGenome()

		list = append(list, bot)
	}
	e.list = list
}

func (e *Bots) findBot(cord Location, exclude *Bot) *Bot {
	for k, v := range e.list {
		if v.pos == cord && exclude.numb != e.list[k].numb {
			return &e.list[k]
		}
	}
	return nil
}

func (e *Bots) pair(bot1 *Bot, bot2 *Bot) {
	var NewCord = bot1.pos
	bot := NewBot(&NewCord)
	bot.Genome = append(make([]int, 0), bot1.Genome...)
	for i := 0; i < len(bot.Genome); i += 2 {
		bot.Genome[i] = bot2.Genome[i]
	}
	listMutex.Lock()
	e.list = append(e.list, bot)
	listMutex.Unlock()
}

var max = 0
var total = make([]int, 100)
var upp = 1
var down = 1

func (e *Bots) Updater() {
	for true {
		//time.Sleep(time.Second / 100)

		e.cycles++
		var wg sync.WaitGroup
		wg.Add(len(e.list))
		for k := range e.list {
			go func(k int) {

				e.list[k].Action()
				defer wg.Done()
			}(k)
		}
		wg.Wait()
		MainWorld.Do()

		var botsToNext []Bot

		sort.Slice(e.list, func(i, j int) bool {
			return e.list[i].Health > e.list[j].Health
		})

		for _, bot := range e.list {
			if bot.Health > 0 {
				botsToNext = append(botsToNext, bot)
			} else {
				MainWorld.Add(bot.pos, "p", make([]int, 0))
			}
		}

		if len(botsToNext) <= 8 {
			MainWorld.Clear()
			MainWorld.FillEat()
			if len(botsToNext) == 0 {
				botsToNext = append(botsToNext, e.list[0])
			}
			fmt.Println(botsToNext[0].Genome)
			if len(botsToNext) < 8 {
				lene := len(botsToNext)
				for i := 0; i < 8-lene; i++ {
					botsToNext = append(botsToNext, botsToNext[0])
				}
			}
			if e.iteration%20 == 0 {
				upp = 1
				down = 1
			}
			coup := 0
			for _, v := range total {
				coup += v
			}
			cut := coup / 100
			e.iteration++
			total = append(total[1:], e.cycles)
			coup = 0
			for _, v := range total {
				coup += v
			}
			if cut < coup/100 {
				upp++
			} else {
				down++
			}
			if max < e.cycles {
				fmt.Printf("Поколение: %v - %v                                   \n",
					e.iteration, e.cycles)
				max = e.cycles
			} else {

				fmt.Printf("Поколение: %v - %v(%v)(%v|%v)                       \n",
					e.iteration, e.cycles, max, coup/100, float32(upp)/float32(down))
			}
			e.cycles = 0
			e.list = make([]Bot, 0)
			for _, b := range botsToNext {
				bot := b
				for i := 0; i < e.botCount/8; i++ {
					if i < e.botCount/8/2 {
						bott := NewBot(&Location{Random(MainWorld.Width), Random(MainWorld.Width)})
						bott.Genome = make([]int, 0)
						bott.Genome = append(bott.Genome, bot.Genome...)
						bott.ActualGenome = make([]int, 0)
						bott.ActualGenome = append(bot.ActualGenome)
						bott.tt = bot.tt
						for kk := 0; kk < i+1; kk++ {
							bott.Mutate()
						}
						bott.Reset()
						e.list = append(e.list, bott)
					} else {
						bot.pos = Location{Random(MainWorld.Width), Random(MainWorld.Width)}
						bot.Reset()
						bot.tt = "o"
						e.list = append(e.list, bot)
					}
				}
			}
		} else {
			e.list = make([]Bot, 0)
			for _, bot := range botsToNext {
				e.list = append(e.list, bot)
			}
		}

		listMutex.Lock()
		e.pList = e.list
		listMutex.Unlock()
	}
}
