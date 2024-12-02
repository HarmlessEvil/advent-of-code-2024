package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
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
	reports, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	safeAmount := 0
	for _, report := range reports {
		if report.IsSafe() {
			safeAmount++
		}
	}

	return safeAmount, nil
}

func part2() (int, error) {
	reports, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	safeAmount := 0
	for _, report := range reports {
		if report.IsSafeWithTolerance() {
			safeAmount++
		}
	}

	return safeAmount, nil
}

type Report []int

func (r Report) IsSafe() bool {
	isIncreasing := r[0] < r[1]
	for i := 1; i < len(r); i++ {
		if isIncreasing && r[i-1] > r[i] || !isIncreasing && r[i-1] < r[i] {
			return false
		}

		diff := r[i-1] - r[i]
		if diff == 0 || diff < -3 || diff > 3 {
			return false
		}
	}

	return true
}

func (r Report) IsSafeWithTolerance() bool {
	diffs := make([]int, len(r)-1)
	for i := 1; i < len(r); i++ {
		diffs[i-1] = r[i] - r[i-1]
	}

	for i := range diffs {
		if diffs[i] == 0 || diffs[i] < -3 || diffs[i] > 3 {
			_, withoutCenter, withoutRight := cutLeftCenterRight(r, i)
			return withoutCenter.IsSafe() || withoutRight.IsSafe()
		}

		if i == 0 {
			continue
		}

		sameSign := (diffs[i-1] > 0) == (diffs[i] > 0)
		if !sameSign {
			withoutLeft, withoutCenter, withoutRight := cutLeftCenterRight(r, i)
			return withoutLeft.IsSafe() || withoutCenter.IsSafe() || withoutRight.IsSafe()
		}
	}

	return true
}

func cutLeftCenterRight(r Report, i int) (Report, Report, Report) {
	withoutLeft := slices.Clone(r)
	if i > 0 {
		withoutLeft = slices.Delete(withoutLeft, i-1, i)
	}

	withoutCenter := slices.Clone(r)
	withoutCenter = slices.Delete(withoutCenter, i, i+1)

	withoutRight := slices.Clone(r)
	withoutRight = slices.Delete(withoutRight, i+1, i+2)

	return withoutLeft, withoutCenter, withoutRight
}

func readInput() ([]Report, error) {
	f, err := os.Open("input/day-02.txt")
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var reports []Report

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		levels := strings.Split(scanner.Text(), " ")

		report := make([]int, len(levels))
		for i := range levels {
			level, err := strconv.Atoi(levels[i])
			if err != nil {
				return nil, fmt.Errorf("parse level: %w", err)
			}

			report[i] = level
		}

		reports = append(reports, report)
	}

	return reports, scanner.Err()
}
