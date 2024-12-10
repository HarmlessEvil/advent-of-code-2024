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
	topology, trailHeads, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	sum := 0
	for _, head := range trailHeads {
		sum += countTrailScore(topology, head, map[Vec2D]struct{}{})
	}

	return sum, nil
}

func part2() (int, error) {
	topology, trailHeads, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	sum := 0
	for _, head := range trailHeads {
		sum += countTrailRating(topology, head)
	}

	return sum, nil
}

func countTrailScore(topology Topology, start Vec2D, dest map[Vec2D]struct{}) int {
	height := topology[start.Row][start.Col]

	if height == '9' {
		if _, ok := dest[start]; ok {
			return 0
		}

		dest[start] = struct{}{}
		return 1
	}

	score := 0
	for _, dir := range AllDirections {
		pos := start.Add(dir)

		if !topology.InBounds(pos) {
			continue
		}

		nextHeight := topology[pos.Row][pos.Col]
		if nextHeight == height+1 {
			score += countTrailScore(topology, pos, dest)
		}
	}

	return score
}

func countTrailRating(topology Topology, start Vec2D) int {
	height := topology[start.Row][start.Col]

	if height == '9' {
		return 1
	}

	score := 0
	for _, dir := range AllDirections {
		pos := start.Add(dir)

		if !topology.InBounds(pos) {
			continue
		}

		nextHeight := topology[pos.Row][pos.Col]
		if nextHeight == height+1 {
			score += countTrailRating(topology, pos)
		}
	}

	return score
}

type Vec2D struct {
	Row int
	Col int
}

func (v Vec2D) Add(d Direction) Vec2D {
	return Vec2D{
		Row: v.Row + d.Row,
		Col: v.Col + d.Col,
	}
}

type Direction Vec2D

var (
	DirectionUp    = Direction{Row: -1, Col: 0}
	DirectionDown  = Direction{Row: 1, Col: 0}
	DirectionLeft  = Direction{Row: 0, Col: -1}
	DirectionRight = Direction{Row: 0, Col: 1}

	AllDirections = []Direction{DirectionUp, DirectionDown, DirectionLeft, DirectionRight}
)

type Topology []string

func (t Topology) InBounds(v Vec2D) bool {
	return v.Row >= 0 && v.Row < len(t) && v.Col >= 0 && v.Col < len(t[v.Row])
}

func readInput() (Topology, []Vec2D, error) {
	f, err := os.Open("input/day-10.txt")
	if err != nil {
		return nil, nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var topography Topology
	var trailHeads []Vec2D

	scanner := bufio.NewScanner(f)
	for row := 0; scanner.Scan(); row++ {
		line := scanner.Text()
		for col := range line {
			if line[col] == '0' {
				trailHeads = append(trailHeads, Vec2D{
					Row: row,
					Col: col,
				})
			}
		}

		topography = append(topography, line)
	}

	return topography, trailHeads, scanner.Err()
}
