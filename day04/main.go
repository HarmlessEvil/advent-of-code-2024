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
	puzzle, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	columns := len(puzzle[0])

	wordCount := 0

	// horizontal
	for _, s := range puzzle {
		for col := 3; col < columns; col++ {
			word := s[col-3 : col+1]

			if word == "XMAS" || word == "SAMX" {
				wordCount++
			}
		}
	}

	// vertical
	for row := 3; row < len(puzzle); row++ {
		for col := 0; col < columns; col++ {
			word := string([]byte{
				puzzle[row-3][col],
				puzzle[row-2][col],
				puzzle[row-1][col],
				puzzle[row][col],
			})

			if word == "XMAS" || word == "SAMX" {
				wordCount++
			}
		}
	}

	// main diagonals
	for row := 3; row < len(puzzle); row++ {
		for col := 3; col < columns; col++ {
			word := string([]byte{
				puzzle[row-3][col-3],
				puzzle[row-2][col-2],
				puzzle[row-1][col-1],
				puzzle[row][col],
			})

			if word == "XMAS" || word == "SAMX" {
				wordCount++
			}
		}
	}

	// minor diagonals
	for row := 3; row < len(puzzle); row++ {
		for col := 3; col < columns; col++ {
			word := string([]byte{
				puzzle[row-3][columns-1-col+3],
				puzzle[row-2][columns-1-col+2],
				puzzle[row-1][columns-1-col+1],
				puzzle[row][columns-1-col],
			})

			if word == "XMAS" || word == "SAMX" {
				wordCount++
			}
		}
	}

	return wordCount, nil
}

func part2() (int, error) {
	puzzle, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	count := 0

	for row := 2; row < len(puzzle); row++ {
		for col := 2; col < len(puzzle[row]); col++ {
			leftTop := puzzle[row-2][col-2]
			rightTop := puzzle[row-2][col]
			leftBottom := puzzle[row][col-2]
			rightBottom := puzzle[row][col]

			if puzzle[row-1][col-1] == 'A' &&
				(leftTop == 'M' && rightBottom == 'S' || leftTop == 'S' && rightBottom == 'M') &&
				(rightTop == 'M' && leftBottom == 'S' || rightTop == 'S' && leftBottom == 'M') {
				count++
			}
		}
	}

	return count, nil
}

func readInput() ([]string, error) {
	f, err := os.Open("input/day-04.txt")
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var lines []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}
