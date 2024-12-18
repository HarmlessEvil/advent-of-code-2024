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

const simulationSteps = 1024
const gridSize = 70

func part1() (int, error) {
	positions, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	corrupted := map[Vec2D]struct{}{}
	for _, position := range positions[:simulationSteps] {
		corrupted[position] = struct{}{}
	}

	steps := bfs(Vec2D{Row: gridSize, Col: gridSize}, corrupted)

	return steps, nil
}

func part2() (string, error) {
	positions, err := readInput()
	if err != nil {
		return "", fmt.Errorf("readInput: %w", err)
	}

	pos := findPositionThatDisconnectsGraph(Vec2D{Row: gridSize, Col: gridSize}, positions)

	return fmt.Sprintf("%d,%d", pos.Row, pos.Col), nil
}

func findPositionThatDisconnectsGraph(exit Vec2D, positions []Vec2D) Vec2D {
	corrupted := map[Vec2D]struct{}{}
	for _, position := range positions {
		corrupted[position] = struct{}{}

		if bfs(exit, corrupted) == -1 {
			return position
		}
	}

	panic("unreachable")
}

func bfs(exit Vec2D, corrupted map[Vec2D]struct{}) int {
	visited := map[Vec2D]struct{}{}

	q := []PathNode{{}}
	for len(q) > 0 {
		node := q[0]
		q = q[1:]

		if _, ok := visited[node.Pos]; ok {
			continue
		}
		visited[node.Pos] = struct{}{}

		if node.Pos == exit {
			return node.Steps
		}

		for _, dir := range []Dir2D{DirectionUp, DirectionDown, DirectionLeft, DirectionRight} {
			next := node.Pos.Add(dir)
			if next.Row < 0 || next.Row > exit.Row || next.Col < 0 || next.Col > exit.Col {
				continue
			}
			if _, ok := corrupted[next]; ok {
				continue
			}

			q = append(q, PathNode{Pos: next, Steps: node.Steps + 1})
		}
	}

	return -1
}

type PathNode struct {
	Pos   Vec2D
	Steps int
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
	DirectionUp    = Dir2D{Row: -1, Col: 0}
	DirectionDown  = Dir2D{Row: 1, Col: 0}
	DirectionLeft  = Dir2D{Row: 0, Col: -1}
	DirectionRight = Dir2D{Row: 0, Col: 1}
)

func readInput() ([]Vec2D, error) {
	f, err := os.Open("input/day-18.txt")
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var positions []Vec2D

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var pos Vec2D
		if _, err := fmt.Sscanf(scanner.Text(), "%d,%d", &pos.Row, &pos.Col); err != nil {
			return nil, fmt.Errorf("parse line: %w", err)
		}

		positions = append(positions, pos)
	}

	return positions, scanner.Err()
}
