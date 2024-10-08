package models

import (
	"main/utils"
)

type Eats struct {
	list    []Eat
	Percent float64
}

var MainEats = Eats{}

func CreateEats() Eats {
	return Eats{
		list: make([]Eat, int(float64(MainWorld.Height()*MainWorld.Width())*MainEats.Percent)),
	}
}

func (e *Eats) Get() []Eat {
	return e.list
}
func (e *Eats) Set(list []Eat) {
	e.list = list
}

func (e *Eats) Update() {
	newList := make([]Eat, 0)
	for _, eat := range e.list {
		eat.Grow++
		//fmt.Println(eat.Grow)
		if eat.Grow < 1000 {
			newList = append(newList, eat)
		}
	}
	e.list = newList
}

func (e *Eats) Clear() {
	e.list = make([]Eat, int(float64(MainWorld.Height()*MainWorld.Width())*MainEats.Percent))
}

func (e *Eats) FillEat() {
	mapArea := MainWorld.Width() * MainWorld.Height()
	eatCount := int(float64(mapArea) * e.Percent)
	for i := 0; i < eatCount; i++ {
		e.addRandomEat()
	}
}

func (e *Eats) addRandomEat() {
	for true {
		location := NewLocation(utils.Random(MainWorld.Width()), utils.Random(MainWorld.Height()))
		eat, _ := e.get(location)
		if eat == nil {
			e.Create(location)
			break
		}
	}

}

func (e *Eats) Delete(cord Location) {
	eat, index := e.get(cord)
	if eat != nil {
		e.list = append(e.list[:index], e.list[index+1:]...)
	}

}

func (e *Eats) get(cord Location) (*Eat, int) {
	for i, eat := range e.list {
		if eat.Pos.X == cord.X && eat.Pos.Y == cord.Y {
			return &eat, i
		}
	}
	return nil, 0
}

func (e *Eats) Create(location Location) {
	eat := Eat{Pos: location}
	e.addIn(eat)
}
func (e *Eats) addIn(eat Eat) {
	e.list = append(e.list, eat)
}
