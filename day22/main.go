package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

const steps = 2_000

func part1() (int, error) {
	buyers, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	total := 0
	for _, secret := range buyers {
		for range steps {
			secret = next(secret)
		}

		total += secret
	}

	return total, nil
}

const sequenceLength = 4

func part2() (int, error) {
	buyers, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	changes := make([][]int, len(buyers))
	prices := make([][]int, len(buyers))
	sequences := make([]map[Sequence]int, len(buyers))
	indexToSeq := make([][]IndexedSequence, len(buyers))

	for i, secret := range buyers {
		changes[i] = make([]int, steps)
		prices[i] = make([]int, steps)
		sequences[i] = map[Sequence]int{}

		price := secret % 10

		for j := range steps {
			secret = next(secret)

			nextPrice := secret % 10
			prices[i][j] = nextPrice
			changes[i][j] = nextPrice - price

			if j > 2 {
				seq := Sequence(changes[i][j-3 : j+1])
				if _, ok := sequences[i][seq]; !ok {
					sequences[i][seq] = j
					indexToSeq[i] = append(indexToSeq[i], IndexedSequence{Index: j, Sequence: seq})
				}
			}

			price = nextPrice
		}
	}

	mostBananas := 0
	usedSequences := map[Sequence]struct{}{}
	for _, indexedSequences := range indexToSeq {
		for _, seq := range indexedSequences {
			if _, ok := usedSequences[seq.Sequence]; ok {
				continue
			}
			usedSequences[seq.Sequence] = struct{}{}

			bananas := 0
			for buyer, m := range sequences {
				idx, ok := m[seq.Sequence]
				if !ok {
					continue
				}

				bananas += prices[buyer][idx]
			}

			mostBananas = max(mostBananas, bananas)
		}
	}

	return mostBananas, nil
}

type IndexedSequence struct {
	Index    int
	Sequence Sequence
}

type Sequence [sequenceLength]int

func next(num int) int {
	res := prune(mix(num*64, num))
	res = prune(mix(res/32, res))
	res = prune(mix(res*2048, res))

	return res
}

func mix(a, b int) int {
	return a ^ b
}

func prune(a int) int {
	return a % 16777216
}

func readInput() ([]int, error) {
	f, err := os.Open("input/day-22.txt")
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var buyers []int

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		str := scanner.Text()

		num, err := strconv.Atoi(str)
		if err != nil {
			return nil, fmt.Errorf("parse number: %w", err)
		}

		buyers = append(buyers, num)
	}

	return buyers, scanner.Err()
}
