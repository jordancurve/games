// Solver for coin sliding puzzle.
//
// https://twitter.com/TamasGorbe/status/1033723716440674304
//
// Six coins are put on a table as shown on the left. Your task is to get the formation on the right in the least moves possible. A move means sliding a coin, without disturbing the rest, to a new place where it touches two others. Coins must stay on the table at all times.

/*
Hex grid:

    01    13    25    37
 00    12    24    36    48
    11    23    35    47
 10    22    34    46    58
    21    33    45    57
 20    32    44    56    68
    31    43    55    67
 30    42    54    66    78
    41    53    65    77
 40    52    64    76    88
    51    63    75    87

For example, 44 is adjacent to 33, 34, 45, 55, 54, and 43.
*/

package main

import (
	"fmt"
	"sort"
	"strings"
)

type Cell int
type CellSet map[Cell]bool
type CellList []Cell

type Pos CellSet

type Move struct {
	a, z Cell
}

type Entry struct {
	Moves MoveList
	Cset  CellSet
}

type MoveList []Move

func main() {
	start := CellList{34, 33, 44, 32, 43, 54}.Set()
	queue := []Entry{{MoveList{}, start}}
	seen := map[string]bool{}
	maxMoves := 0
	for len(queue) > 0 {
		var entry Entry
		entry, queue = queue[0], queue[1:]
		if len(entry.Moves)+1 > maxMoves {
			maxMoves = len(entry.Moves) + 1
			fmt.Printf("Checking for %d-move solutions...\n", maxMoves)
		}
		cset := entry.Cset
		key := cset.String()
		if seen[key] {
			continue
		}
		seen[key] = true
		for _, m := range cset.Moves() {
			newent := Entry{append(entry.Moves, m), cset.MakeMove(m)}
			if newent.Cset.IsWin() {
				fmt.Printf("Found solution: %s: %v\n", start, newent)
				return
			}
			queue = append(queue, newent)
		}
	}
}

func (cells CellList) String() string {
	cs := []string{}
	for _, c := range cells {
		cs = append(cs, fmt.Sprintf("%d", c))
	}
	return "{" + strings.Join(cs, ", ") + "}"
}

func (cset CellSet) String() string {
	return cset.List().String()
}

func (cset CellSet) MakeMove(move Move) CellSet {
	res := CellSet{}
	for c := range cset {
		if c != move.a {
			res[c] = true
		}
	}
	res[move.z] = true
	return res
}

func (cset CellSet) List() CellList {
	cells := CellList{}
	for c := range cset {
		if cset[c] {
			cells = append(cells, c)
		}
	}
	cells.Sort()
	return cells
}
func (cset CellSet) Moves() MoveList {
	moves := MoveList{}
	for c := range cset {
		if cset.Pinned(c) {
			continue
		}
		for d := Cell(0); d < Cell(100); d++ {
			if cset[d] {
				continue
			}
			delete(cset, c)
			if cset.Pinned(d) {
				cset[c] = true
				continue
			}
			if cset.AdjCount(d) >= 2 {
				moves = append(moves, Move{c, d})
			}
			cset[c] = true
		}
	}
	return moves
}

func (cset CellSet) AdjCount(c Cell) int {
	count := 0
	for _, n := range c.Neighbors() {
		if cset[n] {
			count++
		}
	}
	return count
}

func (c Cell) Neighbors() CellList {
	return filterLegalCells(CellList{c - 11, c - 10, c + 1, c + 11, c + 10, c - 1})
}

func filterLegalCells(cells CellList) CellList {
	cs := CellList{}
	for _, c := range cells {
		if c.legal() {
			cs = append(cs, c)
		}
	}
	return cs
}

func (c Cell) legal() bool {
	return c >= 0 && c <= 99
}

func (cset CellSet) IsWin() bool {
	neighborCount := map[Cell]int{}
	if len(cset) != 6 {
		return false
	}
	for c := range cset {
		for _, n := range c.Neighbors() {
			neighborCount[n]++
			if neighborCount[n] == 6 {
				return true
			}
		}
	}
	return false
}

func (cells CellList) IsWin() bool {
	neighborCount := map[Cell]int{}
	if len(cells) != 6 {
		return false
	}
	for _, c := range cells {
		for _, n := range c.Neighbors() {
			neighborCount[n]++
			if neighborCount[n] == 6 {
				return true
			}
		}
	}
	return false
}

func (cells CellList) Sort() {
	sort.Slice(cells, func(a, b int) bool { return cells[a] < cells[b] })
}

func (moves MoveList) Sort() {
	sort.Slice(moves, func(a, b int) bool {
		return moves[a].a < moves[b].a || (moves[a].a == moves[b].a && moves[a].z < moves[b].z)
	})
}

func (cset CellSet) Pinned(c Cell) bool {
	return (cset[c.UL()] && cset[c.UR()] && cset[c.DL()] && cset[c.DR()]) ||
		(cset[c.UR()] && cset[c.R()] && cset[c.DL()] && cset[c.Left()]) ||
		(cset[c.R()] && cset[c.DR()] && cset[c.Left()] && cset[c.UL()]) ||
		(cset[c.UR()] && cset[c.DR()] && cset[c.Left()]) ||
		(cset[c.R()] && cset[c.DL()] && cset[c.UL()])
}

func (cells CellList) Set() CellSet {
	cset := CellSet{}
	for _, c := range cells {
		cset[c] = true
	}
	return cset
}

func (c Cell) UL() Cell {
	return (c - 11).Filter()
}

func (c Cell) UR() Cell {
	return (c - 10).Filter()
}

func (c Cell) R() Cell {
	return (c + 1).Filter()
}

func (c Cell) DR() Cell {
	return (c + 11).Filter()
}

func (c Cell) DL() Cell {
	return (c + 10).Filter()
}

func (c Cell) Left() Cell {
	return (c - 1).Filter()
}

func (c Cell) Filter() Cell {
	if c < 0 || c > 99 {
		return Cell(-1)
	}
	return c
}

func (m Move) String() string {
	return fmt.Sprintf("%d->%d", m.a, m.z)
}

func (moves MoveList) String() string {
	ms := []string{}
	for _, m := range moves {
		ms = append(ms, m.String())
	}
	return strings.Join(ms, "; ")
}

func (e Entry) String() string {
	return fmt.Sprintf("%s: %s", e.Moves, e.Cset)
}
