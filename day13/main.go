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
	machines, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	tokens := 0
	for _, machine := range machines {
		presses := optimalPrizePresses(machine)
		if presses != nil {
			tokens += presses.X*3 + presses.Y
		}
	}

	return tokens, nil
}

func part2() (int, error) {
	machines, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	tokens := 0
	for _, machine := range machines {
		presses := optimalPrizePresses(Machine{
			A: machine.A,
			B: machine.B,
			Prize: Vec2D{
				X: machine.Prize.X + 10000000000000,
				Y: machine.Prize.Y + 10000000000000,
			},
		})
		if presses != nil {
			tokens += presses.X*3 + presses.Y
		}
	}

	return tokens, nil
}

func optimalPrizePresses(machine Machine) *Vec2D {
	xNominator := machine.Prize.X*machine.B.Y - machine.B.X*machine.Prize.Y
	xDenominator := machine.A.X*machine.B.Y - machine.B.X*machine.A.Y

	if xDenominator == 0 || xNominator%xDenominator != 0 {
		return nil
	}

	yNominator := machine.A.X*machine.Prize.Y - machine.Prize.X*machine.A.Y
	yDenominator := machine.A.X*machine.B.Y - machine.B.X*machine.A.Y

	if yDenominator == 0 || yNominator%yDenominator != 0 {
		return nil
	}

	return &Vec2D{
		X: xNominator / xDenominator,
		Y: yNominator / yDenominator,
	}
}

type Vec2D struct {
	X int
	Y int
}

type Machine struct {
	A     Vec2D
	B     Vec2D
	Prize Vec2D
}

func readInput() ([]Machine, error) {
	f, err := os.Open("input/day-13.txt")
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var machines []Machine

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var m Machine

		if _, err := fmt.Sscanf(scanner.Text(), "Button A: X+%d, Y+%d", &m.A.X, &m.A.Y); err != nil {
			return nil, fmt.Errorf("scan button a: %w", err)
		}

		if !scanner.Scan() {
			break
		}

		if _, err := fmt.Sscanf(scanner.Text(), "Button B: X+%d, Y+%d", &m.B.X, &m.B.Y); err != nil {
			return nil, fmt.Errorf("scan button b: %w", err)
		}

		if !scanner.Scan() {
			break
		}

		if _, err := fmt.Sscanf(scanner.Text(), "Prize: X=%d, Y=%d", &m.Prize.X, &m.Prize.Y); err != nil {
			return nil, fmt.Errorf("scan button b: %w", err)
		}

		if !scanner.Scan() && scanner.Err() != nil {
			break
		}

		machines = append(machines, m)
	}

	return machines, scanner.Err()
}
