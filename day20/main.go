package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
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
	m, start, end, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	visitedNoCheats := simpleBFS(m, start, end)

	fastestTime := visitedNoCheats[end].Picoseconds + 1

	cheats, _ := bfs(m, start, end, visitedNoCheats)

	counts := map[int]int{}
	for _, seconds := range cheats {
		counts[fastestTime-seconds]++
	}

	total := 0
	for amount, count := range counts {
		if amount >= 100 {
			total += count
		}
	}

	return total, nil
}

func part2() (int, error) {
	m, start, end, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	visitedNoCheats := simpleBFS(m, start, end)

	total := 0

	visited := map[Vec2D]struct{}{}
	cheats := map[LongCheat]struct{}{}

	pos := start
	for picoseconds := 0; pos != end; picoseconds++ {
		visited[pos] = struct{}{}

		for dRow := -20; dRow <= 20; dRow++ {
			for dCol := -20; dCol <= 20; dCol++ {
				next := pos.Add(Dir2D{Row: dRow, Col: dCol})
				if !m.InBounds(next) {
					continue
				}

				if m[next.Row][next.Col] == '#' {
					continue
				}

				dist := distManhattan(pos, next)

				if dist == 0 || dist > 20 {
					continue
				}

				cheat := LongCheat{Start: pos, End: next}
				if _, ok := cheats[cheat]; ok {
					continue
				}

				timeSaved := visitedNoCheats[next].Picoseconds - picoseconds - dist + 1
				cheats[cheat] = struct{}{}

				if timeSaved >= 100 {
					total++
				}
			}
		}

		for _, dir := range []Dir2D{DirectionUp, DirectionDown, DirectionLeft, DirectionRight} {
			next := pos.Add(dir)

			if !m.InBounds(next) {
				continue
			}

			if m[next.Row][next.Col] == '#' {
				continue
			}

			if _, ok := visited[next]; ok {
				continue
			}

			pos = next
			break
		}
	}

	return total, nil
}

type LongCheat struct {
	Start Vec2D
	End   Vec2D
}

type Cheat struct {
	Start     Vec2D
	Direction Dir2D
}

func (c Cheat) IsUsed() bool {
	return c.Direction != Dir2D{}
}

type Tile struct {
	Cheat Cheat
	Pos   Vec2D
}

func simpleBFS(m Map, start, end Vec2D) map[Vec2D]PathNode {
	visited := map[Vec2D]PathNode{start: {Tile: Tile{Pos: start}}}

	q := Path{{Tile: Tile{Pos: start}}}
	for len(q) > 0 {
		node := heap.Pop(&q).(PathNode)
		if node.Tile.Pos == end {
			return collectPath(start, end, visited)
		}

		for _, dir := range []Dir2D{DirectionUp, DirectionDown, DirectionLeft, DirectionRight} {
			next := PathNode{
				Picoseconds: node.Picoseconds + 1,
				Tile: Tile{
					Cheat: node.Tile.Cheat,
					Pos:   node.Tile.Pos.Add(dir),
				},
			}

			if !m.InBounds(next.Tile.Pos) {
				continue
			}

			if tile := m[next.Tile.Pos.Row][next.Tile.Pos.Col]; tile == '#' {
				continue
			}

			if _, ok := visited[next.Tile.Pos]; ok {
				continue
			}

			heap.Push(&q, next)
			visited[next.Tile.Pos] = node
		}
	}

	panic("unreachable")
}

func collectPath(start, end Vec2D, visited map[Vec2D]PathNode) map[Vec2D]PathNode {
	res := map[Vec2D]PathNode{
		start: visited[start],
	}

	for end != start {
		from := visited[end]
		res[end] = from

		end = from.Tile.Pos
	}

	return res
}

func bfs(m Map, start, end Vec2D, visitedNoCheats map[Vec2D]PathNode) (map[Cheat]int, int) {
	cheats := map[Cheat]int{}

	fastestTime := visitedNoCheats[end].Picoseconds + 1

	visited := map[Tile]struct{}{}

	q := Path{{Tile: Tile{Pos: start}}}
	for len(q) > 0 {
		node := heap.Pop(&q).(PathNode)
		if node.Tile.Pos == end {
			if !node.Tile.Cheat.IsUsed() {
				return cheats, node.Picoseconds
			}

			cheats[node.Tile.Cheat] = node.Picoseconds
			fmt.Println(cheats)
			continue
		}

		for _, dir := range []Dir2D{DirectionUp, DirectionDown, DirectionLeft, DirectionRight} {
			next := PathNode{
				Picoseconds: node.Picoseconds + 1,
				Tile: Tile{
					Cheat: node.Tile.Cheat,
					Pos:   node.Tile.Pos.Add(dir),
				},
			}

			if !m.InBounds(next.Tile.Pos) {
				continue
			}

			if tile := m[next.Tile.Pos.Row][next.Tile.Pos.Col]; tile == '#' {
				if node.Tile.Cheat.IsUsed() {
					continue
				}

				cheatStart := next.Tile.Pos

				next.Tile.Pos = cheatStart.Add(dir)
				noCheatsNode, ok := visitedNoCheats[next.Tile.Pos]
				if ok {
					cheatTime := fastestTime - noCheatsNode.Picoseconds + next.Picoseconds
					cheats[Cheat{
						Start:     cheatStart,
						Direction: dir,
					}] = cheatTime
				}

				continue
			}

			if _, ok := visitedNoCheats[next.Tile.Pos]; !ok {
				continue
			}

			if _, ok := visited[next.Tile]; ok {
				continue
			}

			heap.Push(&q, next)
			visited[next.Tile] = struct{}{}
		}
	}

	panic("unreachable")
}

type Path []PathNode

func (p Path) Len() int {
	return len(p)
}

func (p Path) Less(i, j int) bool {
	return p[i].Picoseconds < p[j].Picoseconds
}

func (p Path) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p *Path) Push(x any) {
	*p = append(*p, x.(PathNode))
}

func (p *Path) Pop() any {
	n := len(*p)
	x := (*p)[n-1]
	*p = (*p)[:n-1]
	return x
}

type PathNode struct {
	Picoseconds int
	Tile        Tile
}

type Map []string

func (m Map) InBounds(v Vec2D) bool {
	return v.Row >= 0 && v.Row < len(m) && v.Col >= 0 && v.Col < len(m[v.Row])
}

type Vec2D struct {
	Row int
	Col int
}

func (v Vec2D) Add(d Dir2D) Vec2D {
	return Vec2D{Row: v.Row + d.Row, Col: v.Col + d.Col}
}

func distManhattan(a, b Vec2D) int {
	dRow := a.Row - b.Row
	if dRow < 0 {
		dRow = -dRow
	}

	dCol := a.Col - b.Col
	if dCol < 0 {
		dCol = -dCol
	}

	return dRow + dCol
}

type Dir2D Vec2D

var (
	DirectionUp    = Dir2D{Row: -1, Col: 0}
	DirectionDown  = Dir2D{Row: 1, Col: 0}
	DirectionLeft  = Dir2D{Row: 0, Col: -1}
	DirectionRight = Dir2D{Row: 0, Col: 1}
)

func readInput() (Map, Vec2D, Vec2D, error) {
	f, err := os.Open("input/day-20.txt")
	if err != nil {
		return nil, Vec2D{}, Vec2D{}, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var m Map
	var start, end Vec2D

	scanner := bufio.NewScanner(f)
	for row := 0; scanner.Scan(); row++ {
		line := scanner.Text()
		m = append(m, line)

		for col, tile := range line {
			switch tile {
			case 'S':
				start = Vec2D{Row: row, Col: col}
			case 'E':
				end = Vec2D{Row: row, Col: col}
			}
		}
	}

	return m, start, end, scanner.Err()
}
