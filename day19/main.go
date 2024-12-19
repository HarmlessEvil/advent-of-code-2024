package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
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
	available, query, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	total := 0

	slices.Sort(available)
	f := isDesignPossible(available)

	for _, design := range query {
		if f(design) != 0 {
			total++
		}
	}

	return total, nil
}

func part2() (int, error) {
	available, query, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	total := 0

	slices.Sort(available)
	f := isDesignPossible(available)

	for _, design := range query {
		total += f(design)
	}

	return total, nil
}

func isDesignPossible(available []string) (f func(string) int) {
	f = memoize(func(design string) int {
		total := 0

		for i := 0; i < len(design); i++ {
			n := len(design) - i
			if _, ok := slices.BinarySearch(available, design[:n]); ok {
				if i == 0 {
					total++
				}

				total += f(design[n:])
			}
		}

		return total
	})

	return f
}

func memoize[T comparable, R any](f func(T) R) func(T) R {
	memo := map[T]R{}

	return func(t T) R {
		if r, ok := memo[t]; ok {
			return r
		}

		r := f(t)
		memo[t] = r
		return r
	}
}

func readInput() ([]string, []string, error) {
	f, err := os.Open("input/day-19.txt")
	if err != nil {
		return nil, nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	if !scanner.Scan() {
		return nil, nil, scanner.Err()
	}

	available := strings.Split(scanner.Text(), ", ")

	if !scanner.Scan() {
		return nil, nil, scanner.Err()
	}

	var query []string
	for scanner.Scan() {
		query = append(query, scanner.Text())
	}

	return available, query, scanner.Err()
}
