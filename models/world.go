package models

type World struct {
	width  int
	height int
}

var MainWorld World

func CreateWorld(W int, H int) World {
	return World{width: W, height: H}
}

func (w *World) Width() int {
	return w.width
}

func (w *World) Height() int {
	return w.height
}
