package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if res, err := part1(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}

func part1() (int, error) {
	locks, keys, maxHeight, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("read input: %w", err)
	}

	count := 0

	for _, lock := range locks {
		for _, key := range keys {
			if doesKeyFitLock(lock, key, maxHeight) {
				count++
			}
		}
	}

	return count, nil
}

func doesKeyFitLock(lock []int, key []int, maxHeight int) bool {
	for i, height := range lock {
		if key[i]+height > maxHeight {
			return false
		}
	}

	return true
}

func readInput() ([][]int, [][]int, int, error) {
	f, err := os.Open("input/day-25.txt")
	if err != nil {
		return nil, nil, 0, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var locks, keys [][]int
	maxHeight := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var schematics []string

		line := scanner.Text()
		for line != "" {
			schematics = append(schematics, line)

			if !scanner.Scan() {
				break
			}

			line = scanner.Text()
		}

		maxHeight = len(schematics) - 2

		heights := make([]int, len(schematics[0]))
		for col := 0; col < len(schematics[0]); col++ {
			for row := 1; row < len(schematics)-1; row++ {
				if schematics[row][col] == '#' {
					heights[col]++
				}
			}
		}

		if strings.ContainsRune(schematics[0], '#') {
			locks = append(locks, heights)
		} else {
			keys = append(keys, heights)
		}
	}

	return locks, keys, maxHeight, scanner.Err()
}
