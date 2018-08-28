package main

import (
	"reflect"
	"testing"
)

func TestEval(t *testing.T) {
	cases := []struct {
		in   Pos
		want PosValue
	}{
		{Pos{}, 0},
		{Pos{0, [14]Size{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}, 1},
		{Pos{0, [14]Size{1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}, 2},
		{Pos{1, [14]Size{1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}, -2},
		{Pos{0, [14]Size{1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}}, 1},
		{Pos{0, [14]Size{1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}}, 0},
	}
	for _, c := range cases {
		got := c.in.Eval()
		if got != c.want {
			t.Errorf("(%v).Eval()=%v; want %v", c.in, got, c.want)
		}
	}
}

func TestOutOfMoves(t *testing.T) {
	cases := []struct {
		in   Pos
		want Pos
	}{
		{
			Pos{1, [14]Size{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0}},
			Pos{1, [14]Size{3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		},
		{
			Pos{0, [14]Size{1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0}},
			Pos{0, [14]Size{1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0}},
		},
	}
	for _, c := range cases {
		got := c.in.outofmoves()
		if got != c.want {
			t.Errorf("(%v).outofmoves()=%v; want %v", c.in, got, c.want)
		}
	}
}

func TestGameOver(t *testing.T) {
	cases := []struct {
		in   Pos
		want PosValue
	}{
		{Pos{0, [14]Size{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0}}, 0},
		{Pos{0, [14]Size{25, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0}}, infty},
		{Pos{1, [14]Size{1, 0, 0, 0, 0, 0, 0, 28, 1, 0, 0, 0, 1, 0}}, infty},
		{Pos{0, [14]Size{1, 0, 0, 0, 0, 0, 0, 28, 1, 0, 0, 0, 1, 0}}, -infty},
	}
	for _, c := range cases {
		got := c.in.GameOver()
		if got != c.want {
			t.Errorf("(%v).GameOver()=%v; want %v", c.in, got, c.want)
		}
	}
}

func TestMakeMove(t *testing.T) {
	cases := []struct {
		in   Pos
		move Move
		want Pos
	}{
		{
			Pos{0, [14]Size{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}},
			13,
			Pos{0, [14]Size{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		},
		{
			Pos{0, [14]Size{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0}},
			12,
			Pos{0, [14]Size{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}},
		},
		{
			Pos{0, [14]Size{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0}},
			12,
			Pos{1, [14]Size{1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}},
		},
		{
			Pos{0, [14]Size{0, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0}},
			11,
			Pos{1, [14]Size{8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1}},
		},
		{
			Pos{0, [14]Size{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}},
			SkipTurn,
			Pos{0, [14]Size{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}},
		},
		{
			Pos{1, [14]Size{0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0}},
			1,
			Pos{0, [14]Size{0, 0, 1, 1, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0}},
		},
		{
			Pos{1, [14]Size{0, 2, 0, 1, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0}},
			1,
			Pos{0, [14]Size{0, 0, 1, 2, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0}},
		},
	}
	for _, c := range cases {
		got := c.in.makemove(c.move)
		if got != c.want {
			t.Errorf("(%v).makemove(%v)=%v; want %v", c.in, c.move, got, c.want)
		}
	}
}

func TestGenmoves(t *testing.T) {
	cases := []struct {
		in   Pos
		want []Move
	}{
		{
			Pos{},
			[]Move{},
		},
		{
			Pos{0, [14]Size{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
			[]Move{},
		},
		{
			Pos{0, [14]Size{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}},
			[]Move{13},
		},
		{
			Pos{0, [14]Size{0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 5, 0, 2, 1}},
			[]Move{10, 12, 13},
		},
	}
	for _, c := range cases {
		got := c.in.genmoves()
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("(%v).genmoves()=%v; want %v", c.in, got, c.want)
		}
	}
}

func TestNegamax(t *testing.T) {
	cases := []struct {
		in   Pos
		want PosValue
	}{
		{Pos{0, [14]Size{0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1}}, 0},
	}
	depth := byte(5)
	pv := []Move{}
	for _, c := range cases {
		got, ok := c.in.negamax(depth, c.in.Player, -infty, infty, Table{}, &pv)
		if !ok {
			t.Errorf("negamax failed\n")
		}
		if got != c.want {
			t.Errorf("(%v).negamax()=%v; want %v", c.in, got, c.want)
		}
	}
}
