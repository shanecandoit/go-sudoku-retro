package game

import (
	"fmt"
	"strings"
)

// Sudoku represents a 9x9 Sudoku board
type Sudoku struct {
	Board [9][9]int // 0 represents empty cell, 1-9 represent filled cells
}

// NewSudoku creates a new empty Sudoku board
func NewSudoku() *Sudoku {
	return &Sudoku{}
}

// NewSudokuWithPuzzle creates a Sudoku board with a predefined puzzle
func NewSudokuWithPuzzle(puzzle [9][9]int) *Sudoku {
	return &Sudoku{Board: puzzle}
}

// SetCell sets a value at the given position
func (s *Sudoku) SetCell(row, col, value int) bool {
	if !s.isValidPosition(row, col) || !s.isValidValue(value) {
		return false
	}
	s.Board[row][col] = value
	return true
}

// GetCell returns the value at the given position
func (s *Sudoku) GetCell(row, col int) int {
	if !s.isValidPosition(row, col) {
		return -1 // Invalid position
	}
	return s.Board[row][col]
}

// ClearCell clears the cell at the given position
func (s *Sudoku) ClearCell(row, col int) bool {
	if !s.isValidPosition(row, col) {
		return false
	}
	s.Board[row][col] = 0
	return true
}

// IsEmpty checks if a cell is empty
func (s *Sudoku) IsEmpty(row, col int) bool {
	return s.GetCell(row, col) == 0
}

// IsValidMove checks if placing a value at a position is valid according to Sudoku rules
func (s *Sudoku) IsValidMove(row, col, value int) bool {
	if !s.isValidPosition(row, col) || !s.isValidValue(value) {
		return false
	}

	// Save current value
	original := s.Board[row][col]

	// Temporarily place the value
	s.Board[row][col] = value

	// Check if it's valid
	valid := s.isValidPlacement(row, col, value)

	// Restore original value
	s.Board[row][col] = original

	return valid
}

// IsComplete checks if the Sudoku puzzle is completely and correctly filled
func (s *Sudoku) IsComplete() bool {
	// Check if all cells are filled
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if s.Board[i][j] == 0 {
				return false
			}
		}
	}

	// Check if all placements are valid
	return s.IsValid()
}

// IsValid checks if the current board state is valid (no conflicts)
func (s *Sudoku) IsValid() bool {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if s.Board[i][j] != 0 {
				if !s.isValidPlacement(i, j, s.Board[i][j]) {
					return false
				}
			}
		}
	}
	return true
}

// Display prints the Sudoku board to console in a nice format
func (s *Sudoku) Display() {
	fmt.Println("┌───────┬───────┬───────┐")

	for i := 0; i < 9; i++ {
		if i == 3 || i == 6 {
			fmt.Println("├───────┼───────┼───────┤")
		}

		fmt.Print("│ ")
		for j := 0; j < 9; j++ {
			if j == 3 || j == 6 {
				fmt.Print("│ ")
			}

			if s.Board[i][j] == 0 {
				fmt.Print(". ")
			} else {
				fmt.Printf("%d ", s.Board[i][j])
			}
		}
		fmt.Println("│")
	}
	fmt.Println("└───────┴───────┴───────┘")
}

// String returns a string representation of the board
func (s *Sudoku) String() string {
	var builder strings.Builder

	builder.WriteString("┌─────────┬─────────┬─────────┐\n")

	for i := 0; i < 9; i++ {
		if i == 3 || i == 6 {
			builder.WriteString("├─────────┼─────────┼─────────┤\n")
		}

		builder.WriteString("│ ")
		for j := 0; j < 9; j++ {
			if j == 3 || j == 6 {
				builder.WriteString("│ ")
			}

			if s.Board[i][j] == 0 {
				builder.WriteString(". ")
			} else {
				builder.WriteString(fmt.Sprintf("%d ", s.Board[i][j]))
			}
		}
		builder.WriteString("│\n")
	}

	builder.WriteString("└─────────┴─────────┴─────────┘")

	return builder.String()
}

// GetRowConflicts returns the conflicting values in the same row
func (s *Sudoku) GetRowConflicts(row, col int) []int {
	if !s.isValidPosition(row, col) {
		return nil
	}

	value := s.Board[row][col]
	if value == 0 {
		return nil
	}

	var conflicts []int
	for j := 0; j < 9; j++ {
		if j != col && s.Board[row][j] == value {
			conflicts = append(conflicts, j)
		}
	}
	return conflicts
}

// GetColConflicts returns the conflicting values in the same column
func (s *Sudoku) GetColConflicts(row, col int) []int {
	if !s.isValidPosition(row, col) {
		return nil
	}

	value := s.Board[row][col]
	if value == 0 {
		return nil
	}

	var conflicts []int
	for i := 0; i < 9; i++ {
		if i != row && s.Board[i][col] == value {
			conflicts = append(conflicts, i)
		}
	}
	return conflicts
}

// GetBoxConflicts returns the conflicting values in the same 3x3 box
func (s *Sudoku) GetBoxConflicts(row, col int) [][2]int {
	if !s.isValidPosition(row, col) {
		return nil
	}

	value := s.Board[row][col]
	if value == 0 {
		return nil
	}

	boxRow := (row / 3) * 3
	boxCol := (col / 3) * 3

	var conflicts [][2]int
	for i := boxRow; i < boxRow+3; i++ {
		for j := boxCol; j < boxCol+3; j++ {
			if (i != row || j != col) && s.Board[i][j] == value {
				conflicts = append(conflicts, [2]int{i, j})
			}
		}
	}
	return conflicts
}

// Private helper methods

func (s *Sudoku) isValidPosition(row, col int) bool {
	return row >= 0 && row < 9 && col >= 0 && col < 9
}

func (s *Sudoku) isValidValue(value int) bool {
	return value >= 0 && value <= 9
}

func (s *Sudoku) isValidPlacement(row, col, value int) bool {
	if value == 0 {
		return true // Empty cell is always valid
	}

	// Check row
	for j := 0; j < 9; j++ {
		if j != col && s.Board[row][j] == value {
			return false
		}
	}

	// Check column
	for i := 0; i < 9; i++ {
		if i != row && s.Board[i][col] == value {
			return false
		}
	}

	// Check 3x3 box
	boxRow := (row / 3) * 3
	boxCol := (col / 3) * 3
	for i := boxRow; i < boxRow+3; i++ {
		for j := boxCol; j < boxCol+3; j++ {
			if (i != row || j != col) && s.Board[i][j] == value {
				return false
			}
		}
	}

	return true
}

// GetSamplePuzzle returns a sample Sudoku puzzle for testing
func GetSamplePuzzle() [9][9]int {
	return [9][9]int{
		{5, 3, 0, 0, 7, 0, 0, 0, 0},
		{6, 0, 0, 1, 9, 5, 0, 0, 0},
		{0, 9, 8, 0, 0, 0, 0, 6, 0},
		{8, 0, 0, 0, 6, 0, 0, 0, 3},
		{4, 0, 0, 8, 0, 3, 0, 0, 1},
		{7, 0, 0, 0, 2, 0, 0, 0, 6},
		{0, 6, 0, 0, 0, 0, 2, 8, 0},
		{0, 0, 0, 4, 1, 9, 0, 0, 5},
		{0, 0, 0, 0, 8, 0, 0, 7, 9},
	}
}

// GetSampleSolution returns the solution to the sample puzzle
func GetSampleSolution() [9][9]int {
	return [9][9]int{
		{5, 3, 4, 6, 7, 8, 9, 1, 2},
		{6, 7, 2, 1, 9, 5, 3, 4, 8},
		{1, 9, 8, 3, 4, 2, 5, 6, 7},
		{8, 5, 9, 7, 6, 1, 4, 2, 3},
		{4, 2, 6, 8, 5, 3, 7, 9, 1},
		{7, 1, 3, 9, 2, 4, 8, 5, 6},
		{9, 6, 1, 5, 3, 7, 2, 8, 4},
		{2, 8, 7, 4, 1, 9, 6, 3, 5},
		{3, 4, 5, 2, 8, 6, 1, 7, 9},
	}
}
