package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if res, err := part1(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}

	if res, err := part2(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}

func part1() (int, error) {
	m, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	polygons := computePolygons(m)

	total := 0

	set := map[*Polygon]struct{}{}
	for _, polygon := range polygons {
		if _, ok := set[polygon]; ok {
			continue
		}

		set[polygon] = struct{}{}

		total += polygon.Area * polygon.Perimeter
	}

	return total, nil
}

func part2() (int, error) {
	m, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	polygons := computePolygons(m)

	total := 0

	set := map[*Polygon]struct{}{}
	for _, polygon := range polygons {
		if _, ok := set[polygon]; ok {
			continue
		}

		set[polygon] = struct{}{}

		price := polygon.Area * polygon.Sides

		total += price
	}

	return total, nil
}

func computePolygons(m Map) map[Vec2D]*Polygon {
	polygons := map[Vec2D]*Polygon{}

	for row := 0; row < len(m); row++ {
		for col := 0; col < len(m); col++ {
			start := Vec2D{Row: row, Col: col}
			if polygons[start] == nil {
				polygons[start] = new(Polygon)
				findExtents(m, polygons, start)
			}
		}
	}

	return polygons
}

func findExtents(m Map, polygons map[Vec2D]*Polygon, start Vec2D) {
	plant := m[start.Row][start.Col]
	polygon := polygons[start]

	neighbours := 0
	horizontalNeighbours := 0
	verticalNeighbours := 0

	for _, dir := range []Dir2D{DirectionRight, DirectionDown, DirectionLeft, DirectionUp} {
		next := start.Add(dir)

		if !m.InBounds(next) || m[next.Row][next.Col] != plant || polygons[next] == nil {
			continue
		}

		neighbours++
		if dir.IsHorizontal() {
			horizontalNeighbours++
		} else {
			verticalNeighbours++
		}
	}

	polygon.Area++

	switch neighbours {
	case 0:
		polygon.Perimeter += 4
		polygon.Sides += 4
	case 1:
		polygon.Perimeter += 2

		hasHorizontalNeighbour := horizontalNeighbours == 1

		var other Vec2D
		if hasHorizontalNeighbour {
			other = start.Add(DirectionRight)
			if !m.InBounds(other) || m[other.Row][other.Col] != plant || polygons[other] == nil {
				other = start.Add(DirectionLeft)
			}
		} else {
			other = start.Add(DirectionDown)
			if !m.InBounds(other) || m[other.Row][other.Col] != plant || polygons[other] == nil {
				other = start.Add(DirectionUp)
			}
		}

		var directions []Dir2D
		if hasHorizontalNeighbour {
			directions = []Dir2D{DirectionDown, DirectionUp}
		} else {
			directions = []Dir2D{DirectionRight, DirectionLeft}
		}

		diagonalNeighbours := 0
		for _, dir := range directions {
			next := other.Add(dir)

			if !m.InBounds(next) || m[next.Row][next.Col] != plant || polygons[next] == nil {
				continue
			}

			diagonalNeighbours++
		}

		switch diagonalNeighbours {
		case 1:
			polygon.Sides += 2
		case 2:
			polygon.Sides += 4
		}
	case 2:
		var directions []Dir2D
		var base int // how many sides add or remove from polygon if this cell has 0 diagonal neighbours?

		if horizontalNeighbours == verticalNeighbours {
			if h := start.Add(DirectionRight); !m.InBounds(h) || m[h.Row][h.Col] != plant || polygons[h] == nil {
				directions = append(directions, DirectionLeft)
			} else {
				directions = append(directions, DirectionRight)
			}

			if v := start.Add(DirectionDown); !m.InBounds(v) || m[v.Row][v.Col] != plant || polygons[v] == nil {
				directions = append(directions, DirectionUp)
			} else {
				directions = append(directions, DirectionDown)
			}

			directions = []Dir2D{
				directions[0].Add(directions[1].Neg()),
				directions[1].Add(directions[0].Neg()),
			}

			base = -2
		} else {
			directions = []Dir2D{DirectionRightDown, DirectionLeftDown, DirectionLeftUp, DirectionRightUp}

			base = -4
		}

		diagonalNeighbours := 0
		for _, dir := range directions {
			next := start.Add(dir)

			if !m.InBounds(next) || m[next.Row][next.Col] != plant || polygons[next] == nil {
				continue
			}

			diagonalNeighbours++
		}

		polygon.Sides += base + diagonalNeighbours*2
	case 3:
		polygon.Perimeter -= 2

		hasHorizontalFreeSpace := horizontalNeighbours < verticalNeighbours

		var free Vec2D
		if hasHorizontalFreeSpace {
			free = start.Add(DirectionRight)
			if m.InBounds(free) && m[free.Row][free.Col] == plant && polygons[free] != nil {
				free = start.Add(DirectionLeft)
			}
		} else {
			free = start.Add(DirectionDown)
			if m.InBounds(free) && m[free.Row][free.Col] == plant && polygons[free] != nil {
				free = start.Add(DirectionUp)
			}
		}

		var directions []Dir2D
		if hasHorizontalFreeSpace {
			directions = []Dir2D{DirectionDown, DirectionUp}
		} else {
			directions = []Dir2D{DirectionRight, DirectionLeft}
		}

		diagonalNeighbours := 0
		for _, dir := range directions {
			next := free.Add(dir)

			if !m.InBounds(next) || m[next.Row][next.Col] != plant || polygons[next] == nil {
				continue
			}

			diagonalNeighbours++
		}

		polygon.Sides += -4 + diagonalNeighbours*2
	case 4:
		polygon.Perimeter -= 4
		polygon.Sides -= 4
	}

	for _, dir := range []Dir2D{DirectionRight, DirectionDown, DirectionLeft, DirectionUp} {
		next := start.Add(dir)

		if !m.InBounds(next) || m[next.Row][next.Col] != plant || polygons[next] != nil {
			continue
		}

		polygons[next] = polygon
		findExtents(m, polygons, next)
	}
}

type Map []string

func (m Map) InBounds(v Vec2D) bool {
	return v.Row >= 0 && v.Row < len(m) && v.Col >= 0 && v.Col < len(m[v.Row])
}

type Vec2D struct {
	Row int
	Col int
}

func (v Vec2D) Add(d Dir2D) Vec2D {
	return Vec2D{Row: v.Row + d.Row, Col: v.Col + d.Col}
}

type Dir2D Vec2D

var (
	DirectionRight = Dir2D{Row: 0, Col: 1}
	DirectionDown  = Dir2D{Row: 1, Col: 0}
	DirectionLeft  = Dir2D{Row: 0, Col: -1}
	DirectionUp    = Dir2D{Row: -1, Col: 0}

	DirectionRightDown = Dir2D{Row: 1, Col: 1}
	DirectionLeftDown  = Dir2D{Row: 1, Col: -1}
	DirectionLeftUp    = Dir2D{Row: -1, Col: -1}
	DirectionRightUp   = Dir2D{Row: -1, Col: 1}
)

func (d Dir2D) Neg() Dir2D {
	return Dir2D{Row: -d.Row, Col: -d.Col}
}

func (d Dir2D) Add(other Dir2D) Dir2D {
	return Dir2D{Row: d.Row + other.Row, Col: d.Col + other.Col}
}

func (d Dir2D) IsHorizontal() bool {
	return d == DirectionLeft || d == DirectionRight
}

type Polygon struct {
	Area      int
	Perimeter int
	Sides     int
}

func readInput() (Map, error) {
	f, err := os.Open("input/day-12.txt")
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var m Map

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		m = append(m, scanner.Text())
	}

	return m, scanner.Err()
}
