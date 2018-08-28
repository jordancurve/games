package main

import (
	"reflect"
	"testing"
)

func TestIsWin(t *testing.T) {
	cases := []struct {
		cells CellList
		want  bool
	}{
		{[]Cell{}, false},
		{[]Cell{00, 01, 12, 22, 21, 10, 11}, false},
		{[]Cell{00, 01, 12, 22, 21, 10}, true},
		{[]Cell{55, 56, 67, 76, 77, 65}, true},
		{[]Cell{55, 56, 67, 76, 77, 66}, false},
	}
	for _, c := range cases {
		got := c.cells.IsWin()
		if got != c.want {
			t.Errorf("(%v).IsWin()=%v; want %v\n", c.cells, got, c.want)
		}
	}
}

func TestNeighbors(t *testing.T) {
	cases := []struct {
		cell Cell
		want CellList
	}{
		{0, CellList{1, 11, 10}},
		{54, CellList{43, 44, 55, 65, 64, 53}},
	}
	for _, c := range cases {
		got := c.cell.Neighbors()
		got.Sort()
		c.want.Sort()
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("(%v).Neighbors()=%v; want %v", c.cell, got, c.want)
		}
	}
}

func TestPinned(t *testing.T) {
	cases := []struct {
		cells CellList
		cell  Cell
		want  bool
	}{
		{CellList{33, 34, 54, 55}, 44, true},
		{CellList{32, 34, 54, 55}, 44, false},
		{CellList{33, 34, 45, 55}, 44, false},
		{CellList{33, 34, 45, 55, 54}, 44, true},
		{CellList{32, 44, 53}, 43, true},
		{CellList{33, 42, 54}, 43, true},
		{CellList{33, 42, 54}, 43, true},
	}
	for _, c := range cases {
		cset := c.cells.Set()
		got := cset.Pinned(c.cell)
		if got != c.want {
			t.Errorf("(%v).Pinned(%v)=%v; want %v", c.cells, c.cell, got, c.want)
		}
	}
}

func TestMoves(t *testing.T) {
	cases := []struct {
		cells CellList
		want  MoveList
	}{
		{CellList{}, []Move{}},
		{CellList{22, 33}, []Move{}},
		{CellList{22, 33, 32}, []Move{{22, 43}, {32, 23}, {33, 21}}},
		{CellList{22, 23, 33, 43, 44}, []Move{
			{22, 34}, {22, 54}, {22, 32},
			{23, 34}, {23, 54}, {23, 32},
			{44, 32}, {44, 12}, {44, 34},
			{43, 32}, {43, 12}, {43, 34},
		}},
		{CellList{22, 23, 43, 44, 35}, []Move{
			{35, 34}, {35, 54}, {35, 32}, {35, 12},
			{43, 33}, {43, 12}, {43, 24}, {43, 45},
			{44, 33}, {44, 32}, {44, 12}, {44, 24}, {44, 34},
			{22, 33}, {22, 54}, {22, 24}, {22, 45},
			{23, 33}, {23, 32}, {23, 54}, {23, 34}, {23, 45},
		}},
	}
	for _, c := range cases {
		cset := c.cells.Set()
		got := cset.Moves()
		got.Sort()
		c.want.Sort()
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("(%v).Moves()=%v; want %v", c.cells, got, c.want)
		}
	}
}
