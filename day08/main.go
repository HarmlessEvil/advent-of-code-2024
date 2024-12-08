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
	m, antennas, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	antinodes := map[Vec2D]struct{}{}

	for _, nodes := range antennas {
		for i := range nodes {
			for j := i + 1; j < len(nodes); j++ {
				distance := nodes[i].Sub(nodes[j])

				antinode := nodes[i].Add(distance)
				if m.InBounds(antinode) {
					antinodes[antinode] = struct{}{}
				}

				antinode = nodes[j].Sub(distance)
				if m.InBounds(antinode) {
					antinodes[antinode] = struct{}{}
				}
			}
		}
	}

	return len(antinodes), nil
}

func part2() (int, error) {
	m, antennas, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	antinodes := map[Vec2D]struct{}{}

	for _, nodes := range antennas {
		for i := range nodes {
			antinodes[nodes[i]] = struct{}{}

			for j := i + 1; j < len(nodes); j++ {
				distance := nodes[i].Sub(nodes[j])

				antinode := nodes[i].Add(distance)
				for m.InBounds(antinode) {
					antinodes[antinode] = struct{}{}

					antinode = antinode.Add(distance)
				}

				antinode = nodes[j].Sub(distance)
				for m.InBounds(antinode) {
					antinodes[antinode] = struct{}{}

					antinode = antinode.Sub(distance)
				}
			}
		}
	}

	return len(antinodes), nil
}

type Map [][]byte

func (m Map) InBounds(v Vec2D) bool {
	return v.Col >= 0 && v.Col < len(m) && v.Row >= 0 && v.Row < len(m[v.Col])
}

type Vec2D struct {
	Row int
	Col int
}

func (v Vec2D) Add(other Vec2D) Vec2D {
	return Vec2D{
		Row: v.Row + other.Row,
		Col: v.Col + other.Col,
	}
}

func (v Vec2D) Sub(other Vec2D) Vec2D {
	return Vec2D{
		Row: v.Row - other.Row,
		Col: v.Col - other.Col,
	}
}

func readInput() (Map, map[byte][]Vec2D, error) {
	f, err := os.Open("input/day-08.txt")
	if err != nil {
		return nil, nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var antennasMap Map
	antennas := map[byte][]Vec2D{}

	scanner := bufio.NewScanner(f)
	for row := 0; scanner.Scan(); row++ {
		line := scanner.Bytes()
		for col, cell := range line {
			if cell >= '0' && cell <= '9' || cell >= 'A' && cell <= 'Z' || cell >= 'a' && cell <= 'z' {
				antennas[cell] = append(antennas[cell], Vec2D{Row: row, Col: col})
			}
		}

		antennasMap = append(antennasMap, line)
	}

	return antennasMap, antennas, scanner.Err()
}
