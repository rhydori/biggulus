package helper

import "math"

type Vec2 struct {
	X, Y float64
}

func (v Vec2) Normalize() Vec2 {
	l := math.Hypot(v.X, v.Y)
	if l == 0 {
		return Vec2{0, 0}
	}
	return Vec2{v.X / l, v.Y / l}
}
