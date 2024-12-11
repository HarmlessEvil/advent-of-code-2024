package main

import (
	"bytes"
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
	stones, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	for range 25 {
		stones = blink(stones)
	}

	return len(stones), nil
}

func part2() (int, error) {
	stones, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	count := countStones(map[Key]int{}, stones, 75)

	return count, nil
}

func countStones(memo map[Key]int, stones []int, blinks int) int {
	if blinks == 0 {
		return len(stones)
	}

	total := 0

	for _, stone := range stones {
		count, ok := memo[Key{
			Number: stone,
			Steps:  blinks,
		}]

		if !ok {
			next := transform(stone)
			count = countStones(memo, next, blinks-1)

			memo[Key{
				Number: stone,
				Steps:  blinks,
			}] = count
		}

		total += count
	}

	return total
}

type Key struct {
	Number int
	Steps  int
}

func blink(stones []int) []int {
	next := make([]int, 0, len(stones))
	for _, stone := range stones {
		next = append(next, transform(stone)...)
	}

	return next
}

func transform(number int) []int {
	if number == 0 {
		return []int{1}
	}

	if digits, factor := countDigits(number); digits%2 == 0 {
		return []int{number / factor, number % factor}
	}

	return []int{number * 2024}
}

func countDigits(number int) (int, int) {
	count := 0
	factor := 1

	for ; number != 0; count++ {
		number /= 10

		if count%2 == 0 {
			factor *= 10
		}
	}

	return count, factor
}

func readInput() ([]int, error) {
	arrangement, err := os.ReadFile("input/day-11.txt")
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile: %w", err)
	}

	numbers := bytes.Split(arrangement, []byte{' '})

	stones := make([]int, len(numbers))
	for i, number := range numbers {
		n := 0
		for _, digit := range number {
			n = n*10 + int(digit-'0')
		}

		stones[i] = n
	}

	return stones, nil
}
