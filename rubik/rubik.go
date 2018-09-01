// This program is a solver for the 2x2x2 Rubik's cube.
//
// Example usage: go run rubik.go WROOOYRRGGGGWBRYWYBOYBBW
package main

import (
	"fmt"
	"os"
	"strings"
)

type Move []int
type State string

/*
Moves:

     NULL                     UP                     BACK                   RIGHT

     0   1                   0  1                    0  1                    0  9
     2   3                   7  5                    2  3                    2 11

4 5  8   9  12 13      4 16  10  8  2  13     4  5   8  9  12 13       4  5  8  17  14 12
6 7  10 11  14 15      6 17  11  9  3  15     21 20  6  7  10 11       6  7  10 19  15 13

     16 17                   14 12                  18 16                    16 21
     18 19                   18 19                  19 17                    18 23

     20 21                   20 21                  15 14                    20  1
     22 23                   22 23                  22 23                    22  3
*/

var U = []int{0, 1, 7, 5, 4, 16, 6, 17, 10, 8, 11, 9, 2, 13, 3, 15, 14, 12, 18, 19, 20, 21, 22, 23}
var B = []int{0, 1, 2, 3, 4, 5, 21, 20, 8, 9, 6, 7, 12, 13, 10, 11, 18, 16, 19, 17, 15, 14, 22, 23}
var R = []int{0, 9, 2, 11, 4, 5, 6, 7, 8, 17, 10, 19, 14, 12, 15, 13, 16, 21, 18, 23, 20, 1, 22, 3}

type moveName string

var nameToMove = map[moveName]Move{
	"U": U,
	"B": B,
	"R": R,
}

var seen map[State]bool

type MoveSeq struct {
	moves string
	state State
}

func main() {
	s := State(strings.Join(os.Args[1:], ""))

	sol, finalState, ok := solve(s)
	if !ok {
		fmt.Printf("no solution\n")
	} else {
		fmt.Printf("%s --> [%d moves: %q] --> %s\n", s, len(sol), sol, finalState)
	}
}

func solve(s State) (string, State, bool) {
	if s.solved() {
		return "", s, true
	}
	seen := map[State]bool{s: true}
	queue := []MoveSeq{{"", s}}
	for len(queue) > 0 {
		seq := queue[0]
		queue = queue[1:]
		for _, moveName := range []string{"U", "R", "B"} {
			st := seq.state
			nextState := st.makeMoves(moveName)
			nextMoves := seq.moves + moveName
			if nextState.solved() {
				return nextMoves, nextState, true
			}
			if !seen[nextState] {
				seen[nextState] = true
				queue = append(queue, MoveSeq{nextMoves, nextState})
			}
		}
	}
	return "", s, false
}

func (s State) solved() bool {
	for face := 0; face < 6; face++ {
		if !s.solvedFace(face) {
			return false
		}
	}
	return true
}

func (s State) solvedFace(face int) bool {
	for square := 1; square < 4; square++ {
		if s[face*4] != s[face*4+square] {
			return false
		}
	}
	return true
}

func (s State) makeMoves(moves string) State {
	for _, m := range moves {
		s = s.makeMove(nameToMove[moveName(m)])
	}
	return s
}

func (s State) makeMove(move Move) State {
	nextState := make([]byte, len(s))
	for i := 0; i < len(move); i++ {
		nextState[i] = byte(s[move[i]])
	}
	return State(nextState)
}

func (s State) String() string {
	faces := []string{}
	for face := 0; face < 6; face++ {
		faces = append(faces, string(s[4*face:4*face+4]))
	}
	return strings.Join(faces, " ")
}
