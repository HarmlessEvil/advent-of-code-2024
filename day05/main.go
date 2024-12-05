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
	rules, updates, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	sum := 0

	for _, update := range updates {
		if isUpdateCorrectlyOrdered(update, rules) {
			sum += update[len(update)/2]
		}
	}

	return sum, nil
}

func part2() (int, error) {
	rules, updates, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	invertedRules := rules.Invert()

	sum := 0

	for _, update := range updates {
		if !isUpdateCorrectlyOrdered(update, rules) {
			allowed := make(map[int]struct{}, len(update))
			for _, page := range update {
				allowed[page] = struct{}{}
			}

			visited := make(map[int]struct{})
			var ordered []int

			for _, page := range update {
				if _, ok := visited[page]; !ok {
					ordered = topologicalSort(page, invertedRules, allowed, visited, ordered)
				}
			}

			sum += ordered[len(ordered)/2]
		}
	}

	return sum, nil
}

func topologicalSort(page int, rules PageOrderingRules, allowed, visited map[int]struct{}, res []int) []int {
	if _, ok := visited[page]; ok {
		return res
	}

	for next := range rules[page] {
		if _, ok := allowed[next]; ok {
			res = topologicalSort(next, rules, allowed, visited, res)
		}
	}

	visited[page] = struct{}{}
	return append(res, page)
}

func isUpdateCorrectlyOrdered(update []int, rules PageOrderingRules) bool {
	used := make(map[int]bool, len(update))
	for _, page := range update {
		used[page] = false
	}

	for _, page := range update {
		used[page] = true

		beforeSet := rules[page]
		for before := range beforeSet {
			if beforeUsed, ok := used[before]; ok && !beforeUsed {
				return false
			}
		}
	}

	return true
}

type PageOrderingRules map[int]map[int]struct{}

func (r PageOrderingRules) Invert() PageOrderingRules {
	res := make(map[int]map[int]struct{})

	for after, beforeSet := range r {
		for before := range beforeSet {
			if _, ok := res[before]; !ok {
				res[before] = make(map[int]struct{})
			}

			res[before][after] = struct{}{}
		}
	}

	return res
}

func readInput() (PageOrderingRules, [][]int, error) {
	f, err := os.Open("input/day-05.txt")
	if err != nil {
		return nil, nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	rules := make(PageOrderingRules)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			break
		}

		var before, after int
		if _, err := fmt.Sscanf(text, "%d|%d", &before, &after); err != nil {
			return nil, nil, fmt.Errorf("parse rule: %w", err)
		}

		if _, ok := rules[after]; !ok {
			rules[after] = make(map[int]struct{})
		}
		rules[after][before] = struct{}{}
	}

	if scanner.Err() != nil {
		return nil, nil, scanner.Err()
	}

	var updates [][]int

	for scanner.Scan() {
		pages := strings.Split(scanner.Text(), ",")

		pageNumbers := make([]int, len(pages))
		for i, p := range pages {
			page, err := strconv.Atoi(p)
			if err != nil {
				return nil, nil, fmt.Errorf("parse page: %w", err)
			}

			pageNumbers[i] = page
		}

		updates = append(updates, pageNumbers)
	}

	return rules, updates, scanner.Err()
}
