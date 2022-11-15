package main

var GenomeLen = 640
var MaxGenome = 21
var maxHealth = 1000

type Bot struct {
	tt           string
	Genome       []int
	ActualGenome []int
	GenomeStep   int
	pos          Location
	Health       int
	angle        int
	Seeds        int
	puber        int
	age          float64
	numb         int
}

var num int

func (e *Bot) GenNum() {
	num++
	e.numb = num
}

func NewBot(coords *Location) Bot {
	bot := Bot{
		Health:     maxHealth,
		GenomeStep: 0,
		angle:      Random(8),
		Seeds:      20}
	bot.GenNum()
	if coords != nil {
		bot.pos = *coords
	} else {
		bot.pos = Location{Random(MainWorld.Width), Random(MainWorld.Height)}
	}
	return bot
}

func (e *Bot) Reset() {
	e.ActualGenome = make([]int, 0)
	e.angle = Random(8)
	e.GenomeStep = 0
	e.Health = maxHealth
	e.Seeds = 20
	e.puber = 0
	e.age = 0
	e.GenNum()
}

func (e *Bot) Move(cord Location) int {
	e.pos = cord
	return 1
}

var cSin = []int{0, -1, -1, -1, 0, 1, 1, 1}
var cCos = []int{1, 1, 0, -1, -1, -1, 0, 1}

func mSin(ang int) int {
	ang = (8 + ang) % 8
	return cSin[ang]
}

func mCos(ang int) int {
	ang = (8 + ang) % 8
	return cCos[ang]
}

func (e *Bot) front() Location {
	return NewLocation(e.pos.X+mCos(e.angle), e.pos.Y+mSin(e.angle))
}
func (e *Bot) frontleft() Location {
	return NewLocation(e.pos.X+mCos(e.angle+1), e.pos.Y+mSin(e.angle+1))
}
func (e *Bot) frontright() Location {
	return NewLocation(e.pos.X+mCos(e.angle+7), e.pos.Y+mSin(e.angle+7))
}
func (e *Bot) back() Location {
	return NewLocation(e.pos.X+mCos(e.angle+4), e.pos.Y+mSin(e.angle+4))
}
func (e *Bot) backleft() Location {
	return NewLocation(e.pos.X+mCos(e.angle+3), e.pos.Y+mSin(e.angle+3))
}
func (e *Bot) backright() Location {
	return NewLocation(e.pos.X+mCos(e.angle+5), e.pos.Y+mSin(e.angle+5))
}
func (e *Bot) right() Location {
	return NewLocation(e.pos.X+mCos(e.angle+6), e.pos.Y+mSin(e.angle+6))
}
func (e *Bot) left() Location {
	return NewLocation(e.pos.X+mCos(e.angle+2), e.pos.Y+mSin(e.angle+2))
}

func (e *Bot) rotate(angle int) int {
	e.angle = e.angle + angle
	if e.angle >= 8 {
		e.angle -= 8
	}
	if e.angle < 0 {
		e.angle += 8
	}
	return 1
}

func (e *Bot) around() []func() Location {
	funcs := []func() Location{
		e.front,
		e.frontleft,
		e.left,
		e.backleft,
		e.back,
		e.backright,
		e.right,
		e.frontright,
	}
	return funcs
}

func (e *Bot) Action() {
	var done = false
	var refresh = false
	var needFood = 1
	var short_mem = -1
	for !done {
		var step = 1
		done = true
		action := e.Genome[e.GenomeStep]
		e.ActualGenome = append(e.ActualGenome, e.GenomeStep)
		switch action {
		case 0:
			e.Move(e.front())
		case 1:
			e.Move(e.back())
		case 2:
			e.Move(e.right())
		case 3:
			e.Move(e.left())
		case 4:
			if short_mem != -1 {
				e.rotate(short_mem)
			} else {
				e.rotate(-1)
			}
		case 5:
			if short_mem != -1 {
				e.rotate(short_mem)
			} else {
				e.rotate(1)
			}
		case 6:
			ranges := e.around()
			for i, dir := range ranges {
				someFood := MainWorld.Get(dir())
				if someFood != nil && someFood.typo == "e" && someFood.Value == 1 {
					short_mem = i
					break
				}
			}
			if !refresh {
				done = false
				refresh = true
			}
		case 7:
			someFood := MainWorld.Get(e.pos)
			if someFood == nil {
				if e.Seeds > 0 {
					MainWorld.Add(e.pos, "e", e.Genome[:])
					e.Seeds--
					step = 3
				} else {
					step = 2
				}
			}
		case 8:
			ranges := e.around()
			for i, dir := range ranges {
				someFood := MainWorld.Get(dir())
				if someFood != nil {
					if someFood.typo == "p" {
						someFood.typo = "e"
						e.TakeFood(someFood)
						short_mem = i
						if !refresh {
							done = false
							refresh = true
						}
						step = 2
						break
					}
				}
			}
		case 9:
			ranges := e.around()
			for _, dir := range ranges {
				someFood := MainWorld.Get(dir())
				if someFood != nil && someFood.typo == "e" && someFood.Value == 1 {
					e.TakeFood(someFood)
					step = 2
					break
				}
			}
		case 10:
			ranges := e.around()
			for _, dir := range ranges {
				someBot := MainBots.findBot(dir(), e)
				if someBot != nil {
					if e.puber >= 600 && someBot.puber >= 600 {
						MainBots.pair(e, someBot)
						// fmt.Println("SEX")
						e.puber -= 600
						someBot.puber -= 600
						step = 2
						break
					}
				}
			}
		default:
			step = action - 10

			if !refresh {
				done = false
				refresh = true
			}
		}
		e.GenomeStep = (GenomeLen + (e.GenomeStep + step)) % GenomeLen
	}
	var someFood = MainWorld.Get(e.pos)
	if someFood != nil {
		//r := e.Health
		if e.TakeFood(someFood) {
			needFood = 0
		}
		//fmt.Println("FOOD", someFood.Value, r, e.Health)
		//step = 2
	}

	e.Health -= needFood
	e.puber++

	e.age += 1
	if e.age > 1500 {
		e.Health = 0
	}

}

var eatPita = 40

func (e *Bot) TakeFood(eat *Eat) bool {
	if eat.typo == "p" {
		e.Health = 0
	}
	if eat.typo == "e" && eat.placed {
		if eat.Value == 1 {
			e.Health += eatPita
			e.Seeds += 2
		}
		if eat.Value == 2 {
			e.Health += 20
			e.Seeds += 1
		}
		if eat.Value == 3 {
			e.Health += 10
		}
		if e.Health > maxHealth {
			e.Health = maxHealth
		}
		if len(eat.genom) > 0 {
			e.Genome = append(eat.genom)
			e.GenomeStep = 0
		}
		MainWorld.Delete(eat.pos)
		return true
	}
	return false
}

func (e *Bot) Mutate() {
	if len(e.ActualGenome) == 0 {
		e.ActualGenome = append(e.ActualGenome, 0)
	}
	e.ActualGenome = unique(e.ActualGenome)
	ran := Random(len(e.ActualGenome))
	//fmt.Println("MUTATE", e.ActualGenome[ran])
	e.Genome[e.ActualGenome[ran]] = Random(MaxGenome)
}

func (e *Bot) GenerateGenome() {
	e.Genome = make([]int, 0)
	for i := 0; i < GenomeLen; i++ {
		e.Genome = append(e.Genome, 9)
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
