package lib

type Pos struct {
	X float32
	Y float32
}

type Rect struct {
	X, Y, W, H, RR float32
}

func (p Pos) Inside(r Rect) bool {
	return p.X > r.X && p.X < r.X+r.W && p.Y > r.Y && p.Y < r.Y+r.H
}
