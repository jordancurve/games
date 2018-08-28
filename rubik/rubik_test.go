package main

import (
	"testing"
)

func TestSolvedFace(t *testing.T) {
	tests := []struct {
		state State
		face  int
		want  bool
	}{
		{
			"111122223333444455556666",
			0,
			true,
		},
		{
			"111x22223333444455556666",
			0,
			false,
		},
		{
			"x11122223333444455556666",
			1,
			true,
		},
	}
	for _, test := range tests {
		got := test.state.solvedFace(test.face)
		if got != test.want {
			t.Errorf("(%q).solvedFace(%d)=%v; want %v", test.state, test.face, got, test.want)
		}
	}
}

func TestSolved(t *testing.T) {
	tests := []struct {
		state State
		want  bool
	}{
		{
			"111122223333444455556666",
			true,
		},
		{
			"11112222333344445555666x",
			false,
		},
	}
	for _, test := range tests {
		got := test.state.solved()
		if got != test.want {
			t.Errorf("(%q).solved()=%v; want %v", test.state, got, test.want)
		}
	}
}

func TestMakeMove(t *testing.T) {
	tests := []struct {
		state State
		moves string
		want  State
	}{
		{
			"OBYWYRBWGGGGRRYOOORWWBBY",
			"U",
			"OBWRYOBOGGGGYRWOYRRWWBBY",
		},
		{
			"OBYWYRBWGGGGRRYOOORWWBBY",
			"UUUU",
			"OBYWYRBWGGGGRRYOOORWWBBY",
		},
		{
			"OBYWYRBWGGGGRRYOOORWWBBY",
			"RRRR",
			"OBYWYRBWGGGGRRYOOORWWBBY",
		},
		{
			"OBYWYRBWGGGGRRYOOORWWBBY",
			"BBBB",
			"OBYWYRBWGGGGRRYOOORWWBBY",
		},
	}
	for _, test := range tests {
		got := test.state.makeMove(test.moves)
		if got != test.want {
			t.Errorf("(%q).makeMove(%s)=%v; want %v", test.state, test.moves, got, test.want)
		}
	}
}
