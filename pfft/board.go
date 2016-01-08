package main

const boardsize = 15

// kelimelik board
var (
	boardLetterMult = [boardsize][boardsize]int{
		{1, 1, 1, 1, 1, 2, 1, 1, 1, 2, 1, 1, 1, 1, 1},
		{1, 3, 1, 1, 1, 1, 2, 1, 2, 1, 1, 1, 1, 3, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 3, 1, 1, 1, 1, 1, 3, 1, 1, 1, 1},
		{2, 1, 1, 1, 1, 2, 1, 1, 1, 2, 1, 1, 1, 1, 2},
		{1, 2, 1, 1, 1, 1, 2, 1, 2, 1, 1, 1, 1, 2, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 2, 1, 1, 1, 1, 2, 1, 2, 1, 1, 1, 1, 2, 1},
		{2, 1, 1, 1, 1, 2, 1, 1, 1, 2, 1, 1, 1, 1, 2},
		{1, 1, 1, 1, 3, 1, 1, 1, 1, 1, 3, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 3, 1, 1, 1, 1, 2, 1, 2, 1, 1, 1, 1, 3, 1},
		{1, 1, 1, 1, 1, 2, 1, 1, 1, 2, 1, 1, 1, 1, 1},
	}
	boardWordMult = [boardsize][boardsize]int{
		{1, 1, 3, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{3, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 3},
		{1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 2, 1, 1, 1, 1, 2, 1, 1, 1, 1, 2, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1},
		{3, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 3},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 3, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 1},
	}
)