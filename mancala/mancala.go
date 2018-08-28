// mancala game. Usage: mancala Player PocketCount1 PocketCount2 ... PocketCount13
// Press ctrl-c to interrupt search and display chosen move.
package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
)

/*
    13 12 11 10 9 8
  0                  7
	  1  2  3  4  5 6

Player 0's mancala is in pocket 0.
Player 1's mancala is in pocket 7.

*/

type Pos struct {
	Player     Player
	PocketSize [14]Size
}

type Flag byte

const (
	Exact = Flag(0)
	Lower = Flag(1)
	Upper = Flag(2)
)

type Entry struct {
	Flag           Flag
	Value          PosValue
	RemainingDepth byte
}

type Table map[Pos]Entry

type PosValue int

type Size byte

type Move byte

const SkipTurn = Move(255)

type Player byte

type Pocket byte

var infty = PosValue(1e6)

var sign = [2]PosValue{1, -1}

var Mancala = []Pocket{0, 7}

func atoi(s string) int {
	i, err := strconv.ParseInt(strings.Trim(s, "[] "), 10, 64)
	if err != nil {
		panic(err)
	}
	return int(i)
}

var interrupted bool

func main() {
	p := Pos{}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			interrupted = true
		}
	}()
	p.Player = Player(atoi(os.Args[1]))
	for i := 0; i < 14; i++ {
		p.PocketSize[i] = Size(atoi(os.Args[2+i]))
	}
	tbl := Table{}
	interrupted = false
	var depth byte
	fmt.Printf("%v = ...\n", p)
	pv := []Move{}
	for depth = byte(0); depth < byte(30); depth++ {
		v, ok := p.negamax(depth, p.Player, -infty, infty, tbl, &pv)
		if !ok {
			break
		}
		fmt.Printf("depth %d: val=%d pv=%v\n", depth, v, pv)
	}
	if len(pv) > 0 {
		fmt.Printf("-%v move %d-> %v\n", pv, pv[0], p.makemove(pv[0]))
	}
	q := p.makemove(pv[0])
	moves := q.genmoves()
	for _, m := range moves {
		fmt.Printf("  %v\n", q.makemove(m))
	}
}

func (p Pos) negamax(remainingDepth byte, player Player, alpha,
	beta PosValue, tbl Table, pv *[]Move) (PosValue, bool) {
	var line []Move
	//fmt.Printf("negamax(%d, %v, %d)\n", remainingDepth, p, player)
	if interrupted {
		fmt.Printf("interrupted\n")
		return 0, false
	}
	alphaOrig := alpha

	if entry, ok := tbl[p]; ok && entry.RemainingDepth >= remainingDepth {
		switch entry.Flag {
		case Exact:
			return entry.Value, true
		case Lower:
			alpha = max(alpha, entry.Value)
		case Upper:
			beta = min(beta, entry.Value)
		}
		if alpha >= beta {
			return entry.Value, true
		}
	}
	g := p.GameOver()
	if abs(g) > 0 {
		//fmt.Printf("game over: %v\n", g)
		return g, true
	}
	if remainingDepth == 0 {
		return PosValue(sign[p.Player]) * p.Eval(), true
	}
	var ms []Move
	if p.Player == player {
		ms = p.genmoves()
		if len(ms) == 0 {
			return p.outofmoves().Eval(), true
		}
	} else {
		ms = []Move{SkipTurn}
	}
	v := -infty
	for _, m := range ms {
		rd := remainingDepth
		if m != SkipTurn {
			rd = remainingDepth - 1
		}
		w, ok := p.makemove(m).negamax(rd, (1 - player), -beta, -alpha, tbl, &line)
		if !ok {
			return 0, false
		}
		v = max(v, -w)
		if v > alpha {
			alpha = v
			*pv = append([]Move{m}, line...)
		}
		if alpha >= beta {
			break
		}
	}
	entry := Entry{}
	entry.Value = v
	if v <= alphaOrig {
		entry.Flag = Upper
	} else if v >= beta {
		entry.Flag = Lower
	} else {
		entry.Flag = Exact
	}
	entry.RemainingDepth = remainingDepth
	tbl[p] = entry

	return v, true
}

func (p Pos) genmoves() []Move {
	moves := []Move{}
	m := Mancala[(1 - p.Player)]
	for i := Pocket(m + 1); i < Pocket(m+7); i++ {
		if p.PocketSize[i] > 0 {
			moves = append(moves, Move(i))
		}
	}
	return moves
}

func (p Pos) makemove(m Move) Pos {
	if m == SkipTurn {
		return p
	}
	n := p.PocketSize[m]
	p.PocketSize[m] = 0
	sq := Pocket(m + 1)
	for i := Size(1); i <= n; i++ {
		if sq%14 == Mancala[(1-p.Player)] {
			sq++
		}
		p.PocketSize[sq%14]++
		sq++
	}
	sq--
	sq %= 14
	man := Mancala[1-p.Player]
	if sq > man && sq < man+7 && p.PocketSize[sq] == 1 {
		p.PocketSize[Mancala[p.Player]] += p.PocketSize[14-sq]
		p.PocketSize[14-sq] = 0
	}
	if sq != Mancala[p.Player] {
		p.Player = 1 - p.Player
	}
	return p
}

func (p Pos) Eval() PosValue {
	//	return PosValue(sign[p.Player]) * (PosValue(p.PocketSize[0]) - PosValue(p.PocketSize[7]))

	a, b := Size(0), Size(0)
	for i := 0; i < 7; i++ {
		a += p.PocketSize[i]
		b += p.PocketSize[i+7]
	}
	return PosValue(sign[p.Player]) * (PosValue(a) - PosValue(b))
}

func (p Pos) GameOver() PosValue {
	if p.PocketSize[Mancala[p.Player]] > 24 {
		return infty
	}
	if p.PocketSize[Mancala[1-p.Player]] > 24 {
		return -infty
	}
	return 0
}

func abs(v PosValue) PosValue {
	if v < 0 {
		return -v
	}
	return v
}

func (p Pos) outofmoves() Pos {
	m := Mancala[p.Player]
	n := Mancala[1-p.Player]
	for i := m + 1; i < m+7; i++ {
		p.PocketSize[n] += p.PocketSize[i]
		p.PocketSize[i] = 0
	}
	return p
}

func max(x, y PosValue) PosValue {
	if x > y {
		return x
	}
	return y
}

func min(x, y PosValue) PosValue {
	if x < y {
		return x
	}
	return y
}
