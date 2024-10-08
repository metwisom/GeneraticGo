package models

import (
	"main/utils"
)

var GenomeLen = 512
var MaxGenome = 16
var maxHealth = 50

type Bot struct {
	Genome       []int
	ActualGenome []int
	GenomeStep   int
	Pos          Location
	Health       int

	Age int

	seed int

	id int
}

var num int

func (e *Bot) ID() int {
	if e.id == 0 {
		num++
		e.id = num
	}
	return e.id
}

func NewBot(coords *Location) Bot {
	bot := Bot{
		Health:     maxHealth,
		GenomeStep: 0,
		seed:       10}

	bot.Move(Location{utils.Random(MainWorld.Width()), utils.Random(MainWorld.Height())})

	return bot
}

func (e *Bot) Reset() {
	//fmt.Println(e.ActualGenome)
	e.ActualGenome = make([]int, 0)
	e.GenomeStep = 0
	e.Health = maxHealth
}

func (e *Bot) Move(cord Location) int {
	e.Pos = cord
	return 1
}

var cSin = []int{0, 1, 0, -1}
var cCos = []int{1, 0, -1, 0}

func mSin(ang int) int {
	ang = (4 + ang) % 4
	return cSin[ang]
}

func mCos(ang int) int {
	ang = (4 + ang) % 4
	return cCos[ang]
}

func (e *Bot) front(dist int) Location {
	return NewLocation(e.Pos.X+0*dist, e.Pos.Y+1*dist)
}
func (e *Bot) right(dist int) Location {
	return NewLocation(e.Pos.X+1*dist, e.Pos.Y+0*dist)
}
func (e *Bot) back(dist int) Location {
	return NewLocation(e.Pos.X+0*dist, e.Pos.Y-1*dist)
}
func (e *Bot) left(dist int) Location {
	return NewLocation(e.Pos.X-1*dist, e.Pos.Y+0*dist)
}
func (e *Bot) place() {
	location := NewLocation(e.Pos.X-1, e.Pos.Y+0)
	MainEats.Create(location)
	e.seed--
}

//4, 8, 0, 6, 0, 1, 0, 7, 9, 0

//0 Шаг							1	2
//1 Поворот						1	2
//2 голоден ли					1	2	3
//3 Есть ли рядом разможиться		1	2
//4 Посадить						1	2	3
//5 Размножиться					1	2
//6 Где еда						1	2
//7 Сьесть						1	2
//8 Посадить					1	2

func (e *Bot) Action() {

	var step = 1

	//if len(e.Genome) == 0 {
	//
	//	fmt.Println(e.Genome, e.ID())
	//	fmt.Println(e.GenomeStep, e.ID())
	//}
	//fmt.Println("action", e.ID())
	action := e.Genome[e.GenomeStep]
	e.ActualGenome = append(e.ActualGenome, e.GenomeStep)
	e.ActualGenome = unique(e.ActualGenome)
	//fmt.Println(e.ActualGenome)
	switch action {
	case 0:
		e.Move(e.front(1))
	case 1:
		e.Move(e.back(1))
	case 2:
		e.Move(e.left(1))
	case 3:
		e.Move(e.right(1))
	case 4:
		if e.seed > 0 {
			e.place()
			step = 2
		}
	default:
		step = action - 3

	}
	e.GenomeStep = (GenomeLen + (e.GenomeStep + step)) % GenomeLen

	eat, _ := MainEats.get(e.Pos)
	if eat != nil && eat.Grow > 50 && e.Health <= 80 {
		e.Health += 20
		MainEats.Delete(e.Pos)
		e.seed++
		//MainEats.addRandomEat()
	}

	e.Health--
	e.Age++

}

func (e *Bot) Mutate() {
	if len(e.ActualGenome) == 0 {
		e.ActualGenome = append(e.ActualGenome, 0)
	}
	e.ActualGenome = unique(e.ActualGenome)
	ran := utils.Random(len(e.ActualGenome))
	e.Genome[e.ActualGenome[ran]] = utils.Random(MaxGenome)
}

func (e *Bot) GenerateGenome() {
	e.Genome = make([]int, 0)
	for i := 0; i < GenomeLen; i++ {
		e.Genome = append(e.Genome, 3)
	}
}

func unique(intSlice []int) []int {
	// return intSlice
	keys := make(map[int]bool)
	var list []int
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
