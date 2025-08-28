package game

import (
	"math/rand"
	"time"
)

// Difficulty represents the difficulty level of a Sudoku puzzle
type Difficulty int

const (
	Easy   Difficulty = iota // ~36-46 clues (55-65% filled)
	Medium                   // ~32-35 clues (35-40% filled)
	Hard                     // ~28-31 clues (30-35% filled)
	Expert                   // ~22-27 clues (25-30% filled)
)

// Generator handles Sudoku puzzle generation
type Generator struct {
	rand *rand.Rand
}

// NewGenerator creates a new Sudoku generator
func NewGenerator() *Generator {
	return &Generator{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// NewGeneratorWithSeed creates a new Sudoku generator with a specific seed
func NewGeneratorWithSeed(seed int64) *Generator {
	return &Generator{
		rand: rand.New(rand.NewSource(seed)),
	}
}

// GenerateComplete generates a complete, valid Sudoku solution
func (g *Generator) GenerateComplete() *Sudoku {
	sudoku := NewSudoku()
	g.fillBoard(sudoku)
	return sudoku
}

// GeneratePuzzle generates a puzzle with the specified difficulty
func (g *Generator) GeneratePuzzle(difficulty Difficulty) *Sudoku {
	// First generate a complete solution
	complete := g.GenerateComplete()

	// Then remove cells based on difficulty
	return g.createPuzzle(complete, difficulty)
}

// GeneratePuzzleWithEmptyRatio generates a puzzle with a specific empty cell ratio
// emptyRatio should be between 0.0 and 1.0 (e.g., 0.1 for 10% empty, 90% filled)
func (g *Generator) GeneratePuzzleWithEmptyRatio(emptyRatio float64) *Sudoku {
	if emptyRatio < 0.0 {
		emptyRatio = 0.0
	}
	if emptyRatio > 1.0 {
		emptyRatio = 1.0
	}

	// Generate complete solution
	complete := g.GenerateComplete()

	// Calculate number of cells to remove
	totalCells := 81
	cellsToRemove := int(float64(totalCells) * emptyRatio)

	return g.removeCells(complete, cellsToRemove)
}

// fillBoard fills an empty board with a valid Sudoku solution using backtracking
func (g *Generator) fillBoard(sudoku *Sudoku) bool {
	// Find empty cell
	row, col := g.findEmptyCell(sudoku)
	if row == -1 {
		return true // Board is complete
	}

	// Try numbers 1-9 in random order
	numbers := g.getRandomNumbers()
	for _, num := range numbers {
		if sudoku.IsValidMove(row, col, num) {
			sudoku.SetCell(row, col, num)

			if g.fillBoard(sudoku) {
				return true
			}

			// Backtrack
			sudoku.ClearCell(row, col)
		}
	}

	return false
}

// createPuzzle creates a puzzle from a complete solution based on difficulty
func (g *Generator) createPuzzle(complete *Sudoku, difficulty Difficulty) *Sudoku {
	var cellsToRemove int

	switch difficulty {
	case Easy:
		cellsToRemove = 35 + g.rand.Intn(11) // 35-45 cells removed (36-46 clues)
	case Medium:
		cellsToRemove = 46 + g.rand.Intn(10) // 46-55 cells removed (26-35 clues)
	case Hard:
		cellsToRemove = 50 + g.rand.Intn(10) // 50-59 cells removed (22-31 clues)
	case Expert:
		cellsToRemove = 54 + g.rand.Intn(10) // 54-63 cells removed (18-27 clues)
	default:
		cellsToRemove = 40 // Default to medium-ish
	}

	return g.removeCells(complete, cellsToRemove)
}

// removeCells removes a specified number of cells from a complete puzzle
func (g *Generator) removeCells(complete *Sudoku, cellsToRemove int) *Sudoku {
	// Create a copy of the complete puzzle
	puzzle := &Sudoku{}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			puzzle.Board[i][j] = complete.Board[i][j]
		}
	}

	// Create list of all cell positions
	positions := make([][2]int, 0, 81)
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			positions = append(positions, [2]int{i, j})
		}
	}

	// Shuffle positions
	g.shufflePositions(positions)

	// Remove cells
	removed := 0
	for _, pos := range positions {
		if removed >= cellsToRemove {
			break
		}

		row, col := pos[0], pos[1]
		if !puzzle.IsEmpty(row, col) {
			// Store original value
			originalValue := puzzle.GetCell(row, col)
			puzzle.ClearCell(row, col)

			// Check if puzzle still has unique solution (simplified check)
			// For now, we'll just remove the cell - in a full implementation,
			// you'd want to verify the puzzle still has a unique solution
			removed++

			// In a more sophisticated implementation, you would:
			// 1. Temporarily remove the cell
			// 2. Check if the puzzle still has a unique solution
			// 3. If not, restore the cell and try next position
			// 4. If yes, keep it removed

			// For now, we'll use a simpler approach and just remove cells
			_ = originalValue // Suppress unused variable warning
		}
	}

	return puzzle
}

// findEmptyCell finds the first empty cell in the board
func (g *Generator) findEmptyCell(sudoku *Sudoku) (int, int) {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if sudoku.IsEmpty(i, j) {
				return i, j
			}
		}
	}
	return -1, -1 // No empty cell found
}

// getRandomNumbers returns numbers 1-9 in random order
func (g *Generator) getRandomNumbers() []int {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	// Fisher-Yates shuffle
	for i := len(numbers) - 1; i > 0; i-- {
		j := g.rand.Intn(i + 1)
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}

	return numbers
}

// shufflePositions shuffles an array of positions
func (g *Generator) shufflePositions(positions [][2]int) {
	for i := len(positions) - 1; i > 0; i-- {
		j := g.rand.Intn(i + 1)
		positions[i], positions[j] = positions[j], positions[i]
	}
}

// GetDifficultyName returns the string name of a difficulty level
func GetDifficultyName(difficulty Difficulty) string {
	switch difficulty {
	case Easy:
		return "Easy"
	case Medium:
		return "Medium"
	case Hard:
		return "Hard"
	case Expert:
		return "Expert"
	default:
		return "Unknown"
	}
}

// GetDifficultyFromEmptyRatio converts an empty ratio to an approximate difficulty
func GetDifficultyFromEmptyRatio(emptyRatio float64) Difficulty {
	if emptyRatio <= 0.45 {
		return Easy
	} else if emptyRatio <= 0.60 {
		return Medium
	} else if emptyRatio <= 0.75 {
		return Hard
	} else {
		return Expert
	}
}

// CountEmptyCells counts the number of empty cells in a puzzle
func CountEmptyCells(sudoku *Sudoku) int {
	count := 0
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if sudoku.IsEmpty(i, j) {
				count++
			}
		}
	}
	return count
}

// GetEmptyRatio calculates the ratio of empty cells (0.0 to 1.0)
func GetEmptyRatio(sudoku *Sudoku) float64 {
	emptyCells := CountEmptyCells(sudoku)
	return float64(emptyCells) / 81.0
}

// GetFilledRatio calculates the ratio of filled cells (0.0 to 1.0)
func GetFilledRatio(sudoku *Sudoku) float64 {
	return 1.0 - GetEmptyRatio(sudoku)
}
