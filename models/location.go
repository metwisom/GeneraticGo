package models

type Location struct {
	X, Y int
}

func (c Location) Set(X int, Y int) Location {
	c.X = (MainWorld.Width() + X) % MainWorld.Width()
	c.Y = (MainWorld.Height() + Y) % MainWorld.Height()
	return c
}
func (c Location) SetLoc(newPos Location) Location {
	c.X = newPos.X % MainWorld.Width()
	c.Y = newPos.Y % MainWorld.Height()
	return c
}

func (c Location) Bias(biasX int, biasY int) Location {
	c.Set(c.X+biasX, c.Y+biasY)
	return c
}

func NewLocation(X int, Y int) Location {
	location := Location{X: 0, Y: 0}
	location = location.Set(X, Y)
	return location
}
