package main

type Location struct {
	X, Y int
}

func (c Location) Set(X int, Y int) Location {
	c.X = (MainWorld.Width + X) % MainWorld.Width
	c.Y = (MainWorld.Height + Y) % MainWorld.Height
	return c
}

func (c Location) Bias(biasX int, biasY int) Location {
	c.Set(c.X+biasX, c.Y+biasY)
	return c
}

func NewLocation(X int, Y int) Location {
	var location Location
	return location.Set(X, Y)
}
