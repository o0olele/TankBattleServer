package common

import "math"

type Dot struct {
	X float64
	Y float64
}

type Line struct {
	K float64
	B float64
}

func GetLine(a Dot, B Dot) Line {
	K := (a.Y - B.Y) / (a.X - B.X)
	lb := (a.X*B.Y - B.X*a.Y) / (a.X - B.X)
	ret := Line{K: K, B: lb}
	return ret
}

func GetDDDistance(d1, d2 Dot) float64 {
	dis := math.Sqrt((d1.X-d2.X)*(d1.X-d2.X) + (d1.Y-d2.Y)*(d1.Y-d2.Y))
	return dis
}

func GetDLDistance(line Line, dot Dot) float64 {
	dis := math.Abs(line.K*dot.X-dot.Y+line.B) / math.Sqrt(line.K*line.K+1)
	return dis
}

//cos(bac)
func TriCos(a, b, c Dot) float64 {
	a1 := b.X - a.X
	a2 := b.Y - a.Y
	b1 := c.X - a.X
	b2 := c.Y - a.Y
	res := (a1*b1 + a2*b2) / GetDDDistance(a, b) * GetDDDistance(a, c)
	return res
}
