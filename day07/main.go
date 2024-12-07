package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	equations, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	total := 0

	for _, equation := range equations {
		if equation.Match(Add, Multiply) {
			total += equation.Result
		}
	}

	return total, nil
}

func part2() (int, error) {
	equations, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	total := 0

	for _, equation := range equations {
		if equation.Match(Add, Multiply, Concatenate) {
			total += equation.Result
		}
	}

	return total, nil
}

type Equation struct {
	Result  int
	Numbers []int
}

type Operation func(a, b int) int

func (e Equation) Match(operations ...Operation) bool {
	if len(e.Numbers) == 1 {
		if e.Result == e.Numbers[0] {
			return true
		} else {
			return false
		}
	}

	if e.Result < e.Numbers[0] {
		return false
	}

	equation := Equation{Result: e.Result, Numbers: e.Numbers[1:]}
	original := e.Numbers[1]

	for _, operation := range operations {
		e.Numbers[1] = operation(e.Numbers[0], e.Numbers[1])

		if equation.Match(operations...) {
			return true
		}

		e.Numbers[1] = original
	}

	return false
}

func Add(a, b int) int {
	return a + b
}

func Multiply(a, b int) int {
	return a * b
}

func Concatenate(a, b int) int {
	factor := 1
	for rem := b; rem > 0; rem /= 10 {
		factor *= 10
	}

	return a*factor + b
}

func readInput() ([]Equation, error) {
	f, err := os.Open("input/day-07.txt")
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var equations []Equation

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), ": ", 2)

		result, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("parse result: %w", err)
		}

		numbersText := strings.Split(parts[1], " ")
		numbers := make([]int, len(numbersText))

		for i, n := range numbersText {
			number, err := strconv.Atoi(n)
			if err != nil {
				return nil, fmt.Errorf("parse number: %w", err)
			}

			numbers[i] = number
		}

		equations = append(equations, Equation{
			Result:  result,
			Numbers: numbers,
		})
	}

	return equations, scanner.Err()
}
