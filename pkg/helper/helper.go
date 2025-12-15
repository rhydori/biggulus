package helper

import "math"

type Vector2 struct {
	X, Y float64
}

func (v Vector2) Normalize() Vector2 {
	l := math.Hypot(v.X, v.Y)
	if l == 0 {
		return Vector2{0, 0}
	}
	return Vector2{v.X / l, v.Y / l}
}
