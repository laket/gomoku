package main

import "time"
import "math/rand"

import "github.com/laket72/gomoku/gomoku"

func main() {
	rand.Seed(time.Now().Unix())
	board := gomoku.GenerateFinishBoard()
	board.Print()
	board.CheckResult()
}
