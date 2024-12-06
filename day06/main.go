package main

import (
	"bufio"
	"bytes"
	"fmt"
	"iter"
	"maps"
	"os"
	"slices"
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
	m, start, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	positions := 0

	for pos := range m.Iterator(Position{Pos: start, Dir: DirectionUp}) {
		cell := m.At(pos.Pos)
		if *cell != 'X' {
			positions++
		}

		*cell = 'X'
	}

	return positions, nil
}

func part2() (int, error) {
	m, start, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	obstacleCount := 0

	visits := make(map[Position]struct{})

	for pos := range m.Iterator(Position{Pos: start, Dir: DirectionUp}) {
		nextPos := pos.Pos.Add(pos.Dir)
		if m.InBounds(nextPos) {
			if next := m.At(nextPos); *next == '.' {
				m := m.Clone()
				*m.At(nextPos) = '#'

				visits := maps.Clone(visits)
				if hasCycle(m, pos, visits) {
					obstacleCount++
				}
			}
		}

		*m.At(pos.Pos) = 'X'
		visits[pos] = struct{}{}
	}

	return obstacleCount, nil
}

func hasCycle(m Map, start Position, visits map[Position]struct{}) bool {
	for pos := range m.Iterator(start) {
		if _, ok := visits[pos]; ok {
			return true
		}

		visits[pos] = struct{}{}
	}

	return false
}

type Position struct {
	Pos Vec2D
	Dir Dir2D
}

type Map [][]byte

func (m Map) Iterator(position Position) iter.Seq[Position] {
	return func(yield func(Position) bool) {
		for {
			next := position.Pos.Add(position.Dir)
			for m.InBounds(next) && *m.At(next) == '#' {
				position.Dir = position.Dir.Rotate90DegClockwise()
				next = position.Pos.Add(position.Dir)
			}

			if !yield(position) {
				return
			}

			if !m.InBounds(next) {
				return
			}

			position.Pos = next
		}
	}
}

func (m Map) At(pos Vec2D) *byte {
	return &m[pos.Row][pos.Col]
}

func (m Map) InBounds(pos Vec2D) bool {
	return pos.Row >= 0 && pos.Row < len(m) && pos.Col >= 0 && pos.Col < len(m[0])
}

func (m Map) Clone() Map {
	res := make(Map, len(m))
	for i, s := range m {
		res[i] = slices.Clone(s)
	}

	return res
}

type Vec2D struct {
	Row int
	Col int
}

func (v Vec2D) Add(direction Dir2D) Vec2D {
	return Vec2D{v.Row + direction.Row, v.Col + direction.Col}
}

type Dir2D Vec2D

var (
	DirectionUp    = Dir2D(Vec2D{Row: -1, Col: 0})
	DirectionRight = Dir2D(Vec2D{Row: 0, Col: 1})
	DirectionDown  = Dir2D(Vec2D{Row: 1, Col: 0})
	DirectionLeft  = Dir2D(Vec2D{Row: 0, Col: -1})
)

func (d Dir2D) Rotate90DegClockwise() Dir2D {
	switch d {
	case DirectionUp:
		return DirectionRight
	case DirectionRight:
		return DirectionDown
	case DirectionDown:
		return DirectionLeft
	case DirectionLeft:
		return DirectionUp
	}

	panic("Invalid direction")
}

func readInput() (Map, Vec2D, error) {
	f, err := os.Open("input/day-06.txt")
	if err != nil {
		return nil, Vec2D{}, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var res Map
	var guardPosition Vec2D

	scanner := bufio.NewScanner(f)
	for row := 0; scanner.Scan(); row++ {
		line := slices.Clone(scanner.Bytes())

		col := bytes.IndexByte(line, '^')
		if col != -1 {
			guardPosition = Vec2D{Row: row, Col: col}
		}

		res = append(res, line)
	}

	return res, guardPosition, scanner.Err()
}
