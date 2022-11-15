package main

import (
	"sync"
)

type World struct {
	Width  int
	Height int

	EatPercent float64
	EatBader   float64

	Eat  map[Location]Eat
	fEat []Eat
	cEat []Eat
}

func CreateWorld(W int, H int) World {
	world := World{Width: W, Height: H}
	world.Clear()
	return world
}

func (e *World) Clear() {
	e.Eat = make(map[Location]Eat, int(float64(MainWorld.Height*MainWorld.Width)*MainWorld.EatPercent*2))
	e.fEat = make([]Eat, 0)
	e.cEat = make([]Eat, 0)
}

func (e *World) FillEat() {
	eatCount := int(float64(e.Width*e.Height) * e.EatPercent)
	for i := 0; i < eatCount; i++ {
		for true {
			location := NewLocation(Random(e.Width), Random(e.Height))
			if e.Get(location) == nil {
				eat := CreateEat(location, "e")
				eat.Value = 1
				e.Eat[eat.pos] = eat
				break
			}

		}

	}
}

var EatMutex sync.RWMutex
var fEatMutex sync.RWMutex
var cEatMutex sync.RWMutex

var wg sync.WaitGroup

func (e *World) Do() {
	// start := time.Now()
	//trashs := make([][]Eat, len(e.Eat))
	// fmt.Println(len(e.fEat))
	wg.Add(len(e.Eat))
	for _, eat := range e.Eat {
		go func(eat Eat) {
			defer wg.Done()
			if eat.Value > 1 {
				eat.tick()
			}
			if eat.Value == 1 && Random(1000) == 0 {
				location := NewLocation(eat.pos.X, eat.pos.Y).Bias(Random(3)-1, Random(3)-1)
				if pole := MainWorld.Get(location); pole == nil {
					fEatMutex.Lock()
					e.fEat = append(e.fEat, CreateEat(location, "e"))
					fEatMutex.Unlock()
				}
			}
			if eat.Value > 0 && eat.typo != "" {
				fEatMutex.Lock()
				e.fEat = append(e.fEat, eat)
				fEatMutex.Unlock()
			}
		}(eat)
	}
	wg.Wait()

	e.Eat = make(map[Location]Eat, int(float64(MainWorld.Height*MainWorld.Width)*MainWorld.EatPercent*2))
	for _, v := range e.fEat {
		e.Eat[v.pos] = v
	}

	cEatMutex.Lock()
	e.cEat = make([]Eat, 0)
	e.cEat = e.fEat
	cEatMutex.Unlock()

	e.fEat = make([]Eat, 0)

	// elapsed := time.Since(start)
	// log.Printf("Binomial took %s", elapsed)
}

func (e *World) Get(cord Location) *Eat {
	EatMutex.RLock()
	defer EatMutex.RUnlock()
	if eat, ok := e.Eat[cord]; ok {
		return &eat
	}
	return nil
}

func (e *World) Delete(cord Location) {
	EatMutex.Lock()
	delete(e.Eat, cord)
	EatMutex.Unlock()
	// location := NewLocation(Random(MainWorld.Width), Random(MainWorld.Height))
	// if pole := MainWorld.Get(location); pole == nil {
	// 	fEatMutex.Lock()
	// 	e.fEat = append(e.fEat, CreateEat(location, "e"))
	// 	fEatMutex.Unlock()
	// }
}

func (e *World) Add(cord Location, typo string, genom []int) {
	fEatMutex.Lock()
	eat := CreateEat(cord, typo)
	if len(genom) > 0 {
		eat.genom = append(eat.genom, genom...)
	}
	e.fEat = append(e.fEat, eat)
	fEatMutex.Unlock()
}

func (e *World) GetAllMap() []Eat {
	cEatMutex.RLock()
	defer cEatMutex.RUnlock()
	return e.cEat
}
