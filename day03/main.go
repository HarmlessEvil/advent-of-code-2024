package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
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
	instructions, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	sum := 0
	for _, instruction := range instructions {
		mul, ok := instruction.(Mul)
		if !ok {
			continue
		}

		sum += mul.Evaluate()
	}

	return sum, nil
}

func part2() (int, error) {
	instructions, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	enabled := true

	sum := 0
	for _, instruction := range instructions {
		switch instruction := instruction.(type) {
		case Mul:
			if enabled {
				sum += instruction.Evaluate()
			}
		case Do:
			enabled = true
		case Dont:
			enabled = false
		}
	}

	return sum, nil
}

type Do struct{}

type Dont struct{}

type Mul struct {
	Left  int
	Right int
}

func (m Mul) Evaluate() int {
	return m.Left * m.Right
}

var mulRegex = regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)|do\(\)|don't\(\)`)

func readInput() ([]any, error) {
	f, err := os.Open("input/day-03.txt")
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	memory, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("read memory: %w", err)
	}

	matches := mulRegex.FindAllSubmatch(memory, -1)

	instructions := make([]any, len(matches))
	for i, match := range matches {
		var instruction any

		if match[0][0] == 'm' {
			left, err := strconv.Atoi(string(match[1]))
			if err != nil {
				return nil, fmt.Errorf("parse left: %w", err)
			}

			right, err := strconv.Atoi(string(match[2]))
			if err != nil {
				return nil, fmt.Errorf("parse right: %w", err)
			}

			instruction = Mul{
				Left:  left,
				Right: right,
			}
		} else if match[0][2] == '(' {
			instruction = Do{}
		} else {
			instruction = Dont{}
		}

		instructions[i] = instruction
	}

	return instructions, nil
}
