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
	wires, gates, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	z := getWires(gates, 'z')

	res := calculate(wires, gates, z)

	return res, nil
}

func part2() (string, error) {
	_, gates, err := readInput()
	if err != nil {
		return "", fmt.Errorf("readInput: %w", err)
	}

	followedBy := map[string][]string{}
	for wire, gate := range gates {
		followedBy[gate.LHS] = append(followedBy[gate.LHS], wire)
		followedBy[gate.RHS] = append(followedBy[gate.RHS], wire)
	}

	z := getWires(gates, 'z')

	incorrectGates := make([]string, 0, 8)

	for wire, gate := range gates {
		if wire[0] == 'z' {
			if wire != z[len(z)-1] {
				if _, ok := gate.Op.(OperationXor); !ok {
					incorrectGates = append(incorrectGates, wire)
				}
			}
		} else {
			switch gate.Op.(type) {
			case OperationXor:
				if gate.LHS[0] != 'x' && gate.LHS[0] != 'y' && gate.RHS[0] != 'x' && gate.RHS[0] != 'y' {
					incorrectGates = append(incorrectGates, wire)
				} else {
					if !hasNext[OperationXor](gates, followedBy, wire) {
						incorrectGates = append(incorrectGates, wire)
					}
				}
			case OperationAnd:
				if gate.LHS != "x00" && gate.LHS != "y00" && gate.RHS != "x00" && gate.RHS != "y00" {
					if !hasNext[OperationOr](gates, followedBy, wire) {
						incorrectGates = append(incorrectGates, wire)
					}
				}
			}
		}
	}

	slices.Sort(incorrectGates)

	return strings.Join(incorrectGates, ","), nil
}

func hasNext[T Operation](gates map[string]Gate, followedBy map[string][]string, wire string) bool {
	for _, next := range followedBy[wire] {
		if _, ok := gates[next].Op.(T); ok {
			return true
		}
	}

	return false
}

func calculate(wires map[string]int, gates map[string]Gate, z []string) int {
	simulateAll(z, wires, gates)
	return toNumber(wires, z)
}

func toNumber(wires map[string]int, output []string) int {
	result := 0

	for i, wire := range output {
		result += wires[wire] << i
	}

	return result
}

func simulateAll(output []string, wires map[string]int, gates map[string]Gate) {
	for _, wire := range output {
		simulate(wires, gates, wire)
	}
}

func getWires[T any](wires map[string]T, start byte) []string {
	var res []string
	for wire := range wires {
		if wire[0] == start {
			res = append(res, wire)
		}
	}

	slices.Sort(res)
	return res
}

func simulate(wires map[string]int, gates map[string]Gate, wire string) int {
	if v, ok := wires[wire]; ok {
		return v
	}

	gate := gates[wire]
	v := gate.Op.Eval(simulate(wires, gates, gate.LHS), simulate(wires, gates, gate.RHS))
	wires[wire] = v

	return v
}

type Operation interface {
	Eval(int, int) int
}

type OperationAnd struct{}

func (op OperationAnd) Eval(a, b int) int {
	return a & b
}

type OperationOr struct{}

func (op OperationOr) Eval(a, b int) int {
	return a | b
}

type OperationXor struct{}

func (op OperationXor) Eval(a, b int) int {
	return a ^ b
}

type Gate struct {
	LHS string
	Op  Operation
	RHS string
}

func readInput() (map[string]int, map[string]Gate, error) {
	f, err := os.Open("input/day-24.txt")
	if err != nil {
		return nil, nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	wires := map[string]int{}
	for scanner.Scan() {
		var wire string
		var value int

		line := scanner.Text()
		if line == "" {
			break
		}

		if _, err := fmt.Sscanf(line, "%3s: %d", &wire, &value); err != nil {
			return nil, nil, fmt.Errorf("parse wire: %w", err)
		}

		wires[wire] = value
	}

	gates := map[string]Gate{}
	for scanner.Scan() {
		var out string
		var op string
		var gate Gate

		if _, err := fmt.Sscanf(scanner.Text(), "%s %s %s -> %s", &gate.LHS, &op, &gate.RHS, &out); err != nil {
			return nil, nil, fmt.Errorf("parse gate: %w", err)
		}

		switch op {
		case "AND":
			gate.Op = OperationAnd{}
		case "OR":
			gate.Op = OperationOr{}
		case "XOR":
			gate.Op = OperationXor{}
		}

		gates[out] = gate
	}

	return wires, gates, scanner.Err()
}
