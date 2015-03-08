package main

import "time"
import "os"
import "flag"
import "fmt"
import "math/rand"

import "github.com/laket72/gomoku/gomoku"

func main() {
	rand.Seed(time.Now().Unix())

	var (
		dataPath  string // 出力するデータのパス
		labelPath string // 出力するラベルのパス
		numData   int    // データ数
	)

	flag.StringVar(&dataPath, "o", "data.csv", "output data path")
	flag.StringVar(&labelPath, "l", "label.csv", "label data path")
	flag.IntVar(&numData, "n", 1000, "the number of data")
	flag.Parse()

	boards := make([]*gomoku.Board, 0, numData)
	results := make([]int, 0, numData)

	for i := 0; i < numData; i++ {
		board := gomoku.GenerateFinishBoard()
		boards = append(boards, board)
		results = append(results, int(board.CheckResult()))
	}

	f, _ := os.OpenFile(dataPath, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()

	gomoku.BoardsToCSV(boards, f)

	labelFile, _ := os.OpenFile(labelPath, os.O_WRONLY|os.O_CREATE, 0666)
	defer labelFile.Close()

	for _, v := range results {
		fmt.Fprintf(labelFile, "%d\n", v)
	}

	//boards[0].Print()
	//board.Print()
	//board.CheckResult()
}
