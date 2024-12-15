package main

import (
	"bufio"
	"bytes"
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
	m, robot, movements, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	for _, movement := range movements {
		robot = moveRobot(m, robot, movement)
	}

	total := 0
	for row, line := range m {
		for col, cell := range line {
			if cell == 'O' {
				total += row*100 + col
			}
		}
	}

	return total, nil
}

func part2() (int, error) {
	m, _, movements, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	wideMap, robot := mapToWideMap(m)

	for _, movement := range movements {
		robot = moveRobotOnWideMap(wideMap, robot, movement)
	}

	total := 0
	for row, line := range wideMap {
		for col, cell := range line {
			if cell == '[' {
				total += row*100 + col
			}
		}
	}

	return total, nil
}

type WideMap Map

func (m WideMap) At(v Vec2D) *byte {
	return &m[v.Row][v.Col]
}

func mapToWideMap(m Map) (WideMap, Vec2D) {
	var robot Vec2D

	res := make(WideMap, len(m))
	for row, line := range m {
		res[row] = make([]byte, len(line)*2)

		for col, cell := range line {
			switch cell {
			case '#':
				res[row][col*2] = '#'
				res[row][col*2+1] = '#'
			case 'O':
				res[row][col*2] = '['
				res[row][col*2+1] = ']'
			case '.':
				res[row][col*2] = '.'
				res[row][col*2+1] = '.'
			case '@':
				res[row][col*2] = '@'
				res[row][col*2+1] = '.'

				robot = Vec2D{Row: row, Col: col * 2}
			}
		}
	}

	return res, robot
}

func moveRobot(m Map, robot Vec2D, movement Dir2D) Vec2D {
	next := robot.Add(movement)
	for cell := m.At(next); *cell == 'O'; cell = m.At(next) {
		next = next.Add(movement)
	}

	if cell := m.At(next); *cell == '#' {
		return robot
	}

	nextBox := robot.Add(movement)
	if next != nextBox {
		*m.At(next) = 'O'
	}

	*m.At(nextBox) = '@'
	*m.At(robot) = '.'

	return nextBox
}

func moveRobotOnWideMap(m WideMap, robot Vec2D, movement Dir2D) Vec2D {
	if movement.IsVertical() {
		return moveRobotVerticallyOnWideMap(m, robot, movement)
	}

	return moveRobotHorizontallyOnWideMap(m, robot, movement)
}

func moveRobotHorizontallyOnWideMap(m WideMap, robot Vec2D, movement Dir2D) Vec2D {
	next := robot.Add(movement)
	for cell := m.At(next); *cell == '[' || *cell == ']'; cell = m.At(next) {
		next = next.Add(movement).Add(movement)
	}

	if cell := m.At(next); *cell == '#' {
		return robot
	}

	nextBox := robot.Add(movement)
	if next != nextBox {
		line := m[robot.Row]

		var dst, src []byte
		if movement == DirectionLeft {
			dst = line[next.Col:nextBox.Col]
			src = line[next.Col+1 : nextBox.Col+1]
		} else {
			dst = line[nextBox.Col+1 : next.Col+1]
			src = line[nextBox.Col:next.Col]
		}

		copy(dst, src)
	}

	*m.At(nextBox) = '@'
	*m.At(robot) = '.'

	return nextBox
}

func moveRobotVerticallyOnWideMap(m WideMap, robot Vec2D, movement Dir2D) Vec2D {
	var boxes []Vec2D

	visited := map[Vec2D]struct{}{}

	front := []Vec2D{robot}
	for len(front) > 0 {
		pos := front[0]
		front = front[1:]

		currentCell := *m.At(pos)

		next := pos.Add(movement)
		cell := *m.At(next)

		if cell == '#' || currentCell != '@' && *m.At(next.Add(DirectionRight)) == '#' {
			return robot
		}

		if currentCell == '@' && cell == '.' {
			break
		}

		if currentCell == '@' || cell == currentCell {
			if cell == ']' {
				next = next.Add(DirectionLeft)
			}

			if _, ok := visited[next]; !ok {
				visited[next] = struct{}{}

				front = append(front, next)
				boxes = append(boxes, next)
			}
		} else {
			left := next.Add(DirectionLeft)
			if _, ok := visited[left]; !ok && *m.At(left) == '[' {
				visited[left] = struct{}{}

				front = append(front, left)
				boxes = append(boxes, left)
			}

			right := next.Add(DirectionRight)
			if _, ok := visited[right]; !ok && *m.At(right) == '[' {
				visited[right] = struct{}{}

				front = append(front, right)
				boxes = append(boxes, right)
			}
		}
	}

	slices.Reverse(boxes)
	for _, box := range boxes {
		*m.At(box) = '.'
		*m.At(box.Add(DirectionRight)) = '.'

		next := box.Add(movement)
		*m.At(next) = '['
		*m.At(next.Add(DirectionRight)) = ']'
	}

	next := robot.Add(movement)
	*m.At(next) = '@'
	*m.At(robot) = '.'

	return next
}

type Map [][]byte

func (m Map) At(v Vec2D) *byte {
	return &m[v.Row][v.Col]
}

type Vec2D struct {
	Row int
	Col int
}

func (v Vec2D) Add(d Dir2D) Vec2D {
	return Vec2D{Row: v.Row + d.Row, Col: v.Col + d.Col}
}

type Dir2D Vec2D

var (
	DirectionUp    = Dir2D{Row: -1, Col: 0}
	DirectionDown  = Dir2D{Row: 1, Col: 0}
	DirectionLeft  = Dir2D{Row: 0, Col: -1}
	DirectionRight = Dir2D{Row: 0, Col: 1}
)

func (d Dir2D) IsVertical() bool {
	return d == DirectionUp || d == DirectionDown
}

func readInput() (Map, Vec2D, []Dir2D, error) {
	f, err := os.Open("input/day-15.txt")
	if err != nil {
		return nil, Vec2D{}, nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var m Map
	var robot Vec2D

	scanner := bufio.NewScanner(f)
	for row := 0; scanner.Scan(); row++ {
		line := scanner.Bytes()
		if len(line) == 0 {
			break
		}

		if col := bytes.IndexByte(line, '@'); col >= 0 {
			robot = Vec2D{Row: row, Col: col}
		}

		m = append(m, slices.Clone(line))
	}

	var movements []Dir2D
	for scanner.Scan() {
		line := scanner.Bytes()
		slices.Grow(movements, len(line))

		for _, b := range line {
			switch b {
			case '^':
				movements = append(movements, DirectionUp)
			case 'v':
				movements = append(movements, DirectionDown)
			case '<':
				movements = append(movements, DirectionLeft)
			case '>':
				movements = append(movements, DirectionRight)
			}
		}
	}

	return m, robot, movements, scanner.Err()
}
