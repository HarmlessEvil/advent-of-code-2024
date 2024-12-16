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

	score, _ := dijkstra(m, Reindeer{Pos: start, Dir: Dir2D{Row: 0, Col: 1}}, end)

	return score, nil
}

func part2() (int, error) {
	m, start, end, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	_, tiles := dijkstra(m, Reindeer{Pos: start, Dir: Dir2D{Row: 0, Col: 1}}, end)

	return len(tiles), nil
}

type PathNode struct {
	From  []Reindeer
	Score int
}

func dijkstra(m Map, start Reindeer, end Vec2D) (int, map[Vec2D]struct{}) {
	visited := map[Reindeer]PathNode{}

	tiles := Paths{{Reindeer: start}}
	for len(tiles) > 0 {
		tile := heap.Pop(&tiles).(Tile)
		if node, ok := visited[tile.Reindeer]; ok && node.Score < tile.Score {
			continue
		}

		from := tile.Reindeer

		if from.Pos == end {
			tilesOnShortestPaths := collectTiles(visited, tile.Reindeer)
			return tile.Score, tilesOnShortestPaths
		}

		dir := from.Dir
		for _, next := range []Tile{
			{Score: tile.Score + 1, Reindeer: Reindeer{Dir: dir}},
			{Score: tile.Score + 1_001, Reindeer: Reindeer{Dir: dir.RotateClockwise()}},
			{Score: tile.Score + 1_001, Reindeer: Reindeer{Dir: dir.RotateCounterclockwise()}},
		} {
			next.Reindeer.Pos = from.Pos.Add(next.Reindeer.Dir)
			if m.At(next.Reindeer.Pos) == '#' {
				continue
			}

			node, ok := visited[next.Reindeer]
			if ok && next.Score > node.Score {
				continue
			}

			if !ok || next.Score == node.Score {
				node.From = append(node.From, from)
			} else {
				node.From = []Reindeer{from}
			}

			if !ok || next.Score < node.Score {
				heap.Push(&tiles, next)
			}

			node.Score = next.Score
			visited[next.Reindeer] = node
		}
	}

	panic("unreachable")
}

func collectTiles(visited map[Reindeer]PathNode, end Reindeer) map[Vec2D]struct{} {
	res := map[Vec2D]struct{}{end.Pos: {}}

	q := []Reindeer{end}
	for len(q) > 0 {
		to := q[0]
		q = q[1:]

		node, ok := visited[to]
		if !ok {
			continue
		}

		for _, from := range node.From {
			res[from.Pos] = struct{}{}
			q = append(q, from)
		}
	}

	return res
}

type Paths []Tile

func (p Paths) Len() int {
	return len(p)
}

func (p Paths) Less(i, j int) bool {
	return p[i].Score < p[j].Score
}

func (p Paths) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p *Paths) Push(x any) {
	*p = append(*p, x.(Tile))
}

func (p *Paths) Pop() any {
	old := *p
	n := len(old)
	x := old[n-1]
	*p = old[0 : n-1]
	return x
}

type Tile struct {
	Reindeer Reindeer
	Score    int
}

type Reindeer struct {
	Pos Vec2D
	Dir Dir2D
}

type Vec2D struct {
	Row int
	Col int
}

func (v Vec2D) Add(d Dir2D) Vec2D {
	return Vec2D{Row: v.Row + d.Row, Col: v.Col + d.Col}
}

type Dir2D Vec2D

func (d Dir2D) RotateClockwise() Dir2D {
	return Dir2D{Row: d.Col, Col: -d.Row}
}

func (d Dir2D) RotateCounterclockwise() Dir2D {
	return Dir2D{Row: -d.Col, Col: d.Row}
}

type Map []string

func (m Map) At(pos Vec2D) byte {
	return m[pos.Row][pos.Col]
}

func readInput() (Map, Vec2D, Vec2D, error) {
	f, err := os.Open("input/day-16.txt")
	if err != nil {
		return nil, Vec2D{}, Vec2D{}, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var start, end Vec2D
	var m Map

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
