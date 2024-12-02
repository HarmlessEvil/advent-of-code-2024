package main

import (
	"bufio"
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
	leftList, rightList, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	slices.Sort(leftList)
	slices.Sort(rightList)

	totalDistance := 0
	for i, leftNumber := range leftList {
		rightNumber := rightList[i]

		if leftNumber <= rightNumber {
			totalDistance += rightNumber - leftNumber
		} else {
			totalDistance += leftNumber - rightNumber
		}
	}

	return totalDistance, nil
}

func part2() (int, error) {
	leftList, rightList, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	frequency := make(map[int]int, len(rightList))
	for _, number := range rightList {
		frequency[number]++
	}

	similarityScore := 0
	for _, number := range leftList {
		similarityScore += number * frequency[number]
	}

	return similarityScore, nil
}

func readInput() ([]int, []int, error) {
	f, err := os.Open("input/day-01.txt")
	if err != nil {
		return nil, nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var leftList, rightList []int

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var left, right int
		if _, err := fmt.Sscanf(scanner.Text(), "%d %d", &left, &right); err != nil {
			return nil, nil, fmt.Errorf("parse line: %w", err)
		}

		leftList = append(leftList, left)
		rightList = append(rightList, right)
	}

	return leftList, rightList, scanner.Err()
}
