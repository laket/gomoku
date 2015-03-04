package gomoku

import "errors"
import "fmt"
import "math/rand"
import "encoding/csv"
import "strconv"
import "io"
import "os"

const (
	ROWS    = 15
	COLUMNS = 15
)

// StoneType is one of these
const (
	NONE = iota
	BLACK
	WHITE
)

// Board[r][c] : r, c
type StoneType uint8

// BattleType is one of these
const (
	EVEN = iota
	BLACKWIN
	WHITEWIN
)

type BattleResult uint8

type Board [ROWS][COLUMNS]StoneType

func (board *Board) Print() {
	columns := len(board[0])

	fmt.Printf("   ")
	for c := 0; c < columns; c++ {
		fmt.Printf("%c ", 'a'+c)
	}
	fmt.Println()

	for r := 0; r < len(board); r++ {
		fmt.Printf("%02d ", r)

		for c := 0; c < len(board[r]); c++ {
			var char byte

			switch board[r][c] {
			case NONE:
				char = '-'
			case BLACK:
				char = 'b'
			case WHITE:
				char = 'w'
			default:
				char = '.'
			}

			fmt.Printf("%c ", char)
		}
		fmt.Println()
	}
}

// RandomGenerate generates a random gomoku board.
func RandomGenerate() *Board {
	var board Board

	for r := 0; r < len(board); r++ {
		for c := 0; c < len(board[c]); c++ {
			board[r][c] = StoneType(rand.Intn(3))
		}
	}

	return &board
}

// GenerateFinishBoard generates a board which is finished
// Games finish when black or white wins or random trial takes constant times.
func GenerateFinishBoard() *Board {
	var board Board

	const (
		NUM_TRIAL = 200
	)

	type Point struct {
		r int
		c int
	}
	positionCandidate := make([]Point, ROWS*COLUMNS, ROWS*COLUMNS)

	for r := 0; r < len(board); r++ {
		for c := 0; c < len(board[c]); c++ {
			positionCandidate[r*COLUMNS+c] = Point{r, c}
		}
	}

	numCandidate := len(positionCandidate)

	for i := 0; i < NUM_TRIAL; i++ {
		// postionCandidate[numCandidate:] contains used positions.
		nextIndex := rand.Intn(numCandidate)
		nextPos := positionCandidate[nextIndex]

		// trash process
		positionCandidate[nextIndex] = positionCandidate[numCandidate-1]

		board[nextPos.r][nextPos.c] = StoneType(rand.Intn(3))
	}

	return &board
}

// checkDiagonal seeks a sequence of 5 stones.
// stone is StoneType to be sought
// (startR, startC) is start position to seek
// (dr, dc) is seeking direction
//
// checkDiagonal returns positions contained in the sequence.
// (-1, -1) means "not found".
func (board *Board) seekSequence(stone StoneType, startR int, startC int, dr int, dc int) (int, int) {
	numSeq := 0
	nextR := startR
	nextC := startC

	for 0 <= nextC && nextC < COLUMNS && 0 <= nextR && nextR < ROWS {
		if board[nextR][nextC] == stone {
			numSeq++

			if numSeq == 5 {
				return nextR, nextC
			}
		} else {
			numSeq = 0
		}

		nextR += dr
		nextC += dc
	}

	return -1, -1
}

// CheckResult returns the battle result of board
// If both Black and White satisfy win condition,
// one side, which depends on implement, wins.
func (board *Board) CheckResult() BattleResult {
	checkWithStoneType := func(stone StoneType) (int, int) {
		// check row line

		for r := 0; r < ROWS; r++ {
			posr, posc := board.seekSequence(stone, r, 0, 0, 1)

			if posr != -1 && posc != -1 {
				return posr, posc
			}
		}
		for c := 0; c < COLUMNS; c++ {
			posr, posc := board.seekSequence(stone, 0, c, 1, 0)

			if posr != -1 && posc != -1 {
				return posr, posc
			}
		}

		// check diagonal line
		// start from upper side
		for c := 0; c < COLUMNS; c++ {
			resR, resC := board.seekSequence(stone, 0, c, 1, 1)
			if resR != -1 {
				return resR, resC
			}
			resR, resC = board.seekSequence(stone, 0, c, 1, -1)
			if resR != -1 {
				return resR, resC
			}
		}
		for r := 1; r < ROWS; r++ {
			resR, resC := board.seekSequence(stone, r, 0, 1, 1)
			if resR != -1 {
				return resR, resC
			}
			resR, resC = board.seekSequence(stone, r, COLUMNS-1, 1, -1)
			if resR != -1 {
				return resR, resC
			}
		}

		return -1, -1
	}

	r, c := checkWithStoneType(BLACK)

	if r != -1 {
		fmt.Printf("Black Win : (%d, %c)\n", r, 'a'+c)
		return BLACKWIN
	}

	r, c = checkWithStoneType(WHITE)

	if r != -1 {
		fmt.Printf("White Win : (%d, %c)\n", r, 'a'+c)
		return WHITEWIN
	}

	return EVEN
}

var (
	ErrInvalidBoardSize  = errors.New("input number sequence size is invalid for board")
	ErrInvalidBoardValue = errors.New("input number sequence contains an invalid value")
)

// NewBoard makes a new Board object from number sequence.
// seq is row major order sequence, which consists of gomoku board.
// values of seq is one of StoneType constant.
func NewBoard(seq []uint8) (*Board, error) {
	if len(seq) != ROWS*COLUMNS {
		return nil, ErrInvalidBoardSize
	}

	var board Board

	for i, v := range seq {
		if v != NONE && v != BLACK && v != WHITE {
			return nil, ErrInvalidBoardValue
		}

		board[i/COLUMNS][i%COLUMNS] = StoneType(v)
	}

	return &board, nil
}

// constructs boards from a csv file.
// Now this function is slow
func NewBoardsFromCSV(r io.Reader) ([]*Board, error) {
	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true

	boards := make([]*Board, 0, 32)

	for {
		record, err := reader.Read()

		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}

		seq := make([]uint8, len(record), len(record))
		for i, v := range record {
			tmp, err := strconv.Atoi(v)

			if err != nil {
				fmt.Fprintf(os.Stderr, err.Error())
			}

			seq[i] = uint8(tmp)

		}

		newBoard, err := NewBoard(seq)

		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
		}

		boards = append(boards, newBoard)
	}

	return boards, nil
}

// ToCSV writes csv file.
func BoardsToCSV(boards []*Board, w io.Writer) []string {
	// board[ROWS][COLS]
	writer := csv.NewWriter(w)

	record := make([]string, 0, ROWS*COLUMNS)

	for _, board := range boards {
		record = record[:0]

		for r := 0; r < ROWS; r++ {
			for c := 0; c < COLUMNS; c++ {
				record = append(record, strconv.Itoa(int(board[r][c])))

			}
		}

		writer.Write(record)
	}

	writer.Flush()
	return record
}
