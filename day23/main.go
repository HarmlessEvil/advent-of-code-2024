package main

import (
	"bufio"
	"fmt"
	"iter"
	"maps"
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
	m, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	triangles := map[[3]string]struct{}{}
	for from, edges := range m {
		for to := range edges {
			for computer := range m {
				if _, ok := m[to][computer]; !ok {
					continue
				}

				if _, ok := m[computer][from]; !ok {
					continue
				}

				triangle := [3]string{from, to, computer}
				slices.Sort(triangle[:])

				triangles[triangle] = struct{}{}
			}
		}
	}

	count := 0
	for triangle := range triangles {
		if triangle[0][0] == 't' || triangle[1][0] == 't' || triangle[2][0] == 't' {
			count++
		}
	}

	return count, nil
}

func part2() (string, error) {
	m, err := readInput()
	if err != nil {
		return "", fmt.Errorf("readInput: %w", err)
	}

	var largestClique map[string]struct{}
	for clique := range BronKerbosch(m) {
		if len(clique) > len(largestClique) {
			largestClique = clique
		}
	}

	return strings.Join(slices.Sorted(maps.Keys(largestClique)), ","), nil
}

// BronKerbosch is an enumeration algorithm for finding all maximal cliques in an undirected graph.
//
// Source: https://en.wikipedia.org/wiki/Bron–Kerbosch_algorithm
func BronKerbosch(m map[string]map[string]struct{}) iter.Seq[map[string]struct{}] {
	R := map[string]struct{}{}
	X := map[string]struct{}{}

	P := map[string]struct{}{}
	for v := range m {
		P[v] = struct{}{}
	}

	return func(yield func(map[string]struct{}) bool) {
		var recurse func(R, P, X map[string]struct{}) bool
		recurse = func(R, P, X map[string]struct{}) bool {
			if len(P) == 0 && len(X) == 0 {
				if !yield(R) {
					return false
				}
			}

			for v := range P {
				nextR := maps.Clone(R)
				nextR[v] = struct{}{}

				nextP := maps.Clone(P)
				for p := range P {
					if _, ok := m[v][p]; !ok {
						delete(nextP, p)
					}
				}

				nextX := maps.Clone(X)
				for x := range X {
					if _, ok := m[v][x]; !ok {
						delete(nextX, x)
					}
				}

				if !recurse(nextR, nextP, nextX) {
					return false
				}

				delete(P, v)
				X[v] = struct{}{}
			}

			return true
		}

		recurse(R, P, X)
	}
}

func readInput() (map[string]map[string]struct{}, error) {
	f, err := os.Open("input/day-23.txt")
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	m := map[string]map[string]struct{}{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		edge := scanner.Bytes()

		from := string(edge[:2])
		to := string(edge[3:])

		if _, ok := m[from]; !ok {
			m[from] = map[string]struct{}{}
		}
		m[from][to] = struct{}{}

		if _, ok := m[to]; !ok {
			m[to] = map[string]struct{}{}
		}
		m[to][from] = struct{}{}
	}

	return m, scanner.Err()
}
