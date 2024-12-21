package main

import (
	"bufio"
	"fmt"
	"iter"
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
	codes, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	total := 0

	for _, code := range codes {
		keypad1 := pressAll(0, code)

		keypad2 := make([]string, 0, len(keypad1))
		for _, keypad := range keypad1 {
			keypad2 = append(keypad2, pressAll(1, keypad)...)
		}

		keypad3 := make([]string, 0, len(keypad2))
		for _, keypad := range keypad2 {
			keypad3 = append(keypad3, pressAll(2, keypad)...)
		}

		num := 0
		for i := 0; i < 3; i++ {
			num = num*10 + int(code[i]-'0')
		}

		length := 1 << 62
		for _, sequence := range keypad3 {
			length = min(len(sequence), length)
		}

		complexity := num * length
		total += complexity
	}

	return total, nil
}

func part2() (int, error) {
	getMinSeqLengthMemoized = memoize3(getMinSeqLength)

	codes, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	total := 0

	for _, code := range codes {
		length := 0
		from := 'A'
		for _, to := range code {
			length += getMinSeqLengthMemoized(0, byte(from), byte(to))
			from = to
		}

		num := 0
		for i := 0; i < 3; i++ {
			num = num*10 + int(code[i]-'0')
		}

		complexity := num * length
		total += complexity
	}

	return total, nil
}

const robotCount = 25

func getMinSeqLength(keypadIndex int, from, to byte) int {
	keypad := numericKeypad
	if keypadIndex > 0 {
		keypad = directionalKeypad
	}

	pos := keypad[from]
	next := keypad[to]

	perms := press(keypadIndex, pos, next)

	if keypadIndex == robotCount {
		length := 1 << 62
		for _, perm := range perms {
			length = min(len(perm), length)
		}

		return length
	}

	bestLength := 1 << 62
	for _, perm := range perms {
		length := 0
		from := 'A'
		for _, to := range perm {
			length += getMinSeqLengthMemoized(keypadIndex+1, byte(from), byte(to))
			from = to
		}

		bestLength = min(bestLength, length)
	}

	return bestLength
}

var getMinSeqLengthMemoized func(keypadIndex int, from byte, to byte) int

var press = memoize3(func(keypadIndex int, pos, next Vec2D) []string {
	var presses []byte

	start := pos
	for pos.Col < next.Col {
		pos.Col++
		presses = append(presses, '>')
	}

	for pos.Row > next.Row {
		pos.Row--
		presses = append(presses, '^')
	}

	for pos.Col > next.Col {
		pos.Col--
		presses = append(presses, '<')
	}

	for pos.Row < next.Row {
		pos.Row++
		presses = append(presses, 'v')
	}

	uniquePermutations := map[string]struct{}{}
	var perms []string
	for perm := range permutations(presses) {
		if isValid(keypadIndex, start, perm) {
			p := string(append(perm, 'A'))
			if _, ok := uniquePermutations[p]; !ok {
				uniquePermutations[p] = struct{}{}
				perms = append(perms, p)
			}
		}
	}

	return perms
})

func pressAll(keypadIndex int, code string) []string {
	keypad := numericKeypad
	if keypadIndex > 0 {
		keypad = directionalKeypad
	}

	var parts [][]string

	pos := keypad['A']
	for _, button := range code {
		next := keypad[byte(button)]

		parts = append(parts, press(keypadIndex, pos, next))

		pos = next
	}

	prod := slices.Collect(cartesian(parts))

	res := make([]string, len(prod))
	for i, sequence := range prod {
		res[i] = strings.Join(sequence, "")
	}

	return res
}

func cartesian[T any](m [][]T) iter.Seq[[]T] {
	nextIndex := func(ixs []int) {
		for j := 0; j < len(ixs); j++ {
			idx := len(ixs) - j - 1

			ixs[idx]++

			if idx == 0 || ixs[idx] < len(m[idx]) {
				return
			}

			ixs[idx] = 0
		}
	}

	return func(yield func([]T) bool) {
		indexes := make([]int, len(m))
		for ; indexes[0] < len(m[0]); nextIndex(indexes) {
			buf := make([]T, len(indexes))
			for j, k := range indexes {
				buf[j] = m[j][k]
			}

			if !yield(buf) {
				return
			}
		}
	}
}

func isValid(keypadIndex int, start Vec2D, perm []byte) bool {
	var f func(Vec2D) bool
	if keypadIndex == 0 {
		f = isValidNumeric
	} else {
		f = isValidDirectional
	}

	pos := start
	for _, key := range perm {
		pos = pos.Add(keyToDirection[key])
		if !f(pos) {
			return false
		}
	}

	return true
}

func permutations[T any](s []T) iter.Seq[[]T] {
	s = slices.Clone(s)

	return func(yield func([]T) bool) {
		var permute func(s []T, n int) bool
		permute = func(s []T, n int) bool {
			if n <= 1 {
				return yield(slices.Clone(s))
			}

			if !permute(s, n-1) {
				return false
			}

			for i := 0; i < n-1; i++ {
				if n%2 == 0 {
					s[i], s[n-1] = s[n-1], s[i]
				} else {
					s[0], s[n-1] = s[n-1], s[0]
				}

				if !permute(s, n-1) {
					return false
				}
			}

			return true
		}

		permute(s, len(s))
	}
}

type Vec2D struct {
	Row int
	Col int
}

func (v Vec2D) Add(d Dir2D) Vec2D {
	return Vec2D{Row: v.Row + d.Row, Col: v.Col + d.Col}
}

type Dir2D Vec2D

var numericKeypad = map[byte]Vec2D{
	'7': {Row: 0, Col: 0},
	'8': {Row: 0, Col: 1},
	'9': {Row: 0, Col: 2},
	'4': {Row: 1, Col: 0},
	'5': {Row: 1, Col: 1},
	'6': {Row: 1, Col: 2},
	'1': {Row: 2, Col: 0},
	'2': {Row: 2, Col: 1},
	'3': {Row: 2, Col: 2},
	'0': {Row: 3, Col: 1},
	'A': {Row: 3, Col: 2},
}

func isValidNumeric(pos Vec2D) bool {
	if pos == (Vec2D{Row: 3, Col: 0}) {
		return false
	}

	return pos.Row >= 0 && pos.Row < 4 && pos.Col >= 0 && pos.Col < 3
}

var directionalKeypad = map[byte]Vec2D{
	'^': {Row: 0, Col: 1},
	'A': {Row: 0, Col: 2},
	'<': {Row: 1, Col: 0},
	'v': {Row: 1, Col: 1},
	'>': {Row: 1, Col: 2},
}

func isValidDirectional(pos Vec2D) bool {
	if pos == (Vec2D{}) {
		return false
	}

	return pos.Row >= 0 && pos.Row < 2 && pos.Col >= 0 && pos.Col < 3
}

var keyToDirection = map[byte]Dir2D{
	'^': {Row: -1, Col: 0},
	'<': {Row: 0, Col: -1},
	'v': {Row: 1, Col: 0},
	'>': {Row: 0, Col: 1},
}

func memoize3[P1, P2, P3 comparable, R any](f func(P1, P2, P3) R) func(P1, P2, P3) R {
	type Key struct {
		P1 P1
		P2 P2
		P3 P3
	}

	memo := map[Key]R{}

	return func(p1 P1, p2 P2, p3 P3) R {
		key := Key{P1: p1, P2: p2, P3: p3}

		if r, ok := memo[key]; ok {
			return r
		}

		r := f(p1, p2, p3)
		memo[key] = r
		return r
	}
}

func readInput() ([]string, error) {
	f, err := os.Open("input/day-21.txt")
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var codes []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		codes = append(codes, scanner.Text())
	}

	return codes, scanner.Err()
}
