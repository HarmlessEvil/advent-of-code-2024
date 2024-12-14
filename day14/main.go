package main

import (
	"bufio"
	"cmp"
	"fmt"
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
	robots, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	const seconds = 100
	boundaries := Vec2D{
		Row: 103,
		Col: 101,
	}

	for i := range robots {
		robots[i] = move(robots[i], seconds, boundaries)
	}

	score := computeSafetyScore(robots, boundaries)
	return score, nil
}

func computeSafetyScore(robots []Robot, boundaries Vec2D) int {
	quadrants := [4]int{}
	quadrantSizes := Vec2D{
		Row: boundaries.Row / 2,
		Col: boundaries.Col / 2,
	}

	for _, robot := range robots {
		if robot.Position.Col/quadrantSizes.Col == 0 {
			if robot.Position.Row/quadrantSizes.Row == 0 {
				quadrants[0]++
			} else if robot.Position.Row != quadrantSizes.Row {
				quadrants[2]++
			}
		} else if robot.Position.Col != quadrantSizes.Col {
			if robot.Position.Row/quadrantSizes.Row == 0 {
				quadrants[1]++
			} else if robot.Position.Row != quadrantSizes.Row {
				quadrants[3]++
			}
		}
	}

	return quadrants[0] * quadrants[1] * quadrants[2] * quadrants[3]
}

func part2() (int, error) {
	robots, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	boundaries := Vec2D{
		Row: 103,
		Col: 101,
	}

	states := make([]State, boundaries.Row*boundaries.Col)
	states[0] = State{
		Robots:      robots,
		SafetyScore: computeSafetyScore(robots, boundaries),
		Seconds:     0,
	}

	for seconds := 1; seconds < len(states); seconds++ {
		r := make([]Robot, len(robots))
		for i := range robots {
			r[i] = move(robots[i], seconds, boundaries)
		}

		states[seconds] = State{
			Robots:      r,
			SafetyScore: computeSafetyScore(r, boundaries),
			Seconds:     seconds,
		}
	}

	slices.SortFunc(states, func(a, b State) int {
		return cmp.Compare(a.SafetyScore, b.SafetyScore)
	})

	return states[0].Seconds, nil
}

type State struct {
	Robots      []Robot
	SafetyScore int
	Seconds     int
}

func move(robot Robot, seconds int, boundaries Vec2D) Robot {
	return Robot{
		Position: Vec2D{
			Row: ((robot.Velocity.Row*seconds+robot.Position.Row)%boundaries.Row + boundaries.Row) % boundaries.Row,
			Col: ((robot.Velocity.Col*seconds+robot.Position.Col)%boundaries.Col + boundaries.Col) % boundaries.Col,
		},
		Velocity: robot.Velocity,
	}
}

func printMap(robots []Robot, boundaries Vec2D) {
	rows := make([][]byte, boundaries.Row)
	for i := range rows {
		rows[i] = make([]byte, boundaries.Col)
		for j := range rows[i] {
			rows[i][j] = '.'
		}
	}

	for _, robot := range robots {
		cell := &rows[robot.Position.Row][robot.Position.Col]
		if *cell == '.' {
			*cell = '1'
		} else {
			*cell++
		}
	}

	for _, row := range rows {
		fmt.Println(string(row))
	}
}

type Vec2D struct {
	Row int
	Col int
}

type Robot struct {
	Position Vec2D
	Velocity Vec2D
}

func readInput() ([]Robot, error) {
	f, err := os.Open("input/day-14.txt")
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var robots []Robot

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var robot Robot

		if _, err := fmt.Sscanf(
			scanner.Text(),
			"p=%d,%d v=%d,%d",
			&robot.Position.Col,
			&robot.Position.Row,
			&robot.Velocity.Col,
			&robot.Velocity.Row,
		); err != nil {
			return nil, fmt.Errorf("parse robot: %w", err)
		}

		robots = append(robots, robot)
	}

	return robots, scanner.Err()
}
