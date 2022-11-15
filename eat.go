package main

type Eat struct {
	typo   string
	Value  int
	ticks  float64
	placed bool
	pos    Location
	genom  []int
}

func (e *Eat) tick() {
	e.ticks += MainWorld.EatBader
	if e.ticks >= 1 {
		e.Value--
		e.ticks = 0
	}
}

func CreateEat(cord Location, typo string) Eat {
	return Eat{pos: cord, Value: 4, placed: true, typo: typo, genom: make([]int, 0)}
}
