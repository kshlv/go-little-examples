package main

// Grade ...
type Grade struct {
	Name                  string
	homework, test, paper float64
}

// Grade ...
func (g *Grade) Grade() float64 {
	return (g.homework + g.test + g.paper) / 3.0
}
