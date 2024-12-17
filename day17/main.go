package main

import (
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

func part1() (string, error) {
	vm, err := readInput()
	if err != nil {
		return "", fmt.Errorf("readInput: %w", err)
	}

	vm.ExecuteProgram()

	output := make([]string, len(vm.output))
	for i, num := range vm.output {
		output[i] = strconv.Itoa(num)
	}

	return strings.Join(output, ","), nil
}

func part2() (int, error) {
	vm, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	return dfs(vm, 0, 1), nil
}

func dfs(vm *VM, i int, length int) int {
	for j := 0; j < 8; j++ {
		vm.output = vm.output[:0]

		vm.RegisterA = i + j
		vm.ExecuteProgram()

		if slices.Equal(vm.Program[len(vm.Program)-length:], vm.output) {
			if length == len(vm.Program) {
				return i + j
			}

			if res := dfs(vm, (i+j)*8, length+1); res != -1 {
				return res
			}
		}
	}

	return -1
}

type VM struct {
	RegisterA int
	RegisterB int
	RegisterC int

	Program            []int
	instructionPointer int

	instructions []func(int)

	output []int
}

func (vm *VM) ExecuteProgram() {
	for vm.instructionPointer = 0; vm.instructionPointer < len(vm.Program); vm.instructionPointer += 2 {
		opcode := vm.Program[vm.instructionPointer]
		operand := vm.Program[vm.instructionPointer+1]

		instruction := vm.instructions[opcode]
		instruction(operand)
	}
}

func (vm *VM) dv(operand int) int {
	combo := vm.combo(operand)

	num := vm.RegisterA
	denom := 1 << combo

	return num / denom
}

func (vm *VM) adv(operand int) {
	vm.RegisterA = vm.dv(operand)
}

func (vm *VM) bxl(operand int) {
	vm.RegisterB ^= operand
}

func (vm *VM) bst(operand int) {
	combo := vm.combo(operand)
	vm.RegisterB = combo % 8
}

func (vm *VM) jnz(operand int) {
	if vm.RegisterA == 0 {
		return
	}

	vm.instructionPointer = operand - 2
}

func (vm *VM) bxc(int) {
	vm.RegisterB ^= vm.RegisterC
}

func (vm *VM) out(operand int) {
	combo := vm.combo(operand)

	val := combo % 8
	vm.output = append(vm.output, val)
}

func (vm *VM) bdv(operand int) {
	vm.RegisterB = vm.dv(operand)
}

func (vm *VM) cdv(operand int) {
	vm.RegisterC = vm.dv(operand)
}

func (vm *VM) combo(operand int) int {
	if operand <= 3 {
		return operand
	}

	switch operand {
	case 4:
		return vm.RegisterA
	case 5:
		return vm.RegisterB
	case 6:
		return vm.RegisterC
	}

	panic(fmt.Errorf("unknown operand %d", operand))
}

func readInput() (*VM, error) {
	f, err := os.Open("input/day-17.txt")
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var vm VM
	vm.instructions = []func(int){
		vm.adv, vm.bxl, vm.bst, vm.jnz, vm.bxc, vm.out, vm.bdv, vm.cdv,
	}

	var program string
	if _, err := fmt.Fscanf(
		f,
		"Register A: %d\nRegister B: %d\nRegister C: %d\n\nProgram: %s",
		&vm.RegisterA,
		&vm.RegisterB,
		&vm.RegisterC,
		&program,
	); err != nil {
		return nil, fmt.Errorf("read program: %w", err)
	}

	data := strings.Split(program, ",")
	vm.Program = make([]int, len(data))

	for i, item := range data {
		num, err := strconv.Atoi(item)
		if err != nil {
			return nil, fmt.Errorf("atoi: %w", err)
		}

		vm.Program[i] = num
	}

	return &vm, nil
}
