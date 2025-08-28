package web

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sudoku/game"
	"syscall/js"
)

// WebUI handles the web-based user interface for Sudoku
type WebUI struct {
	sudoku         *game.Sudoku
	originalPuzzle *game.Sudoku // Store original puzzle to know which cells are editable
	selectedRow    int
	selectedCol    int
	selectedNumber int
	gameWon        bool
}

// NewWebUI creates a new web UI instance
func NewWebUI(puzzle *game.Sudoku) *WebUI {
	// Create a copy of the original puzzle
	original := &game.Sudoku{}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			original.Board[i][j] = puzzle.Board[i][j]
		}
	}

	return &WebUI{
		sudoku:         puzzle,
		originalPuzzle: original,
		selectedRow:    4, // Start in center
		selectedCol:    4,
		selectedNumber: 1,
		gameWon:        false,
	}
}

// Run starts the web UI and sets up event handlers
func (w *WebUI) Run() {
	// Register JavaScript functions
	js.Global().Set("selectCell", js.FuncOf(w.selectCell))
	js.Global().Set("setNumber", js.FuncOf(w.setNumber))
	js.Global().Set("clearCell", js.FuncOf(w.clearCell))
	js.Global().Set("newGame", js.FuncOf(w.newGame))
	js.Global().Set("setDifficulty", js.FuncOf(w.setDifficulty))
	js.Global().Set("getSudokuState", js.FuncOf(w.getSudokuState))

	// Initialize the display
	w.updateDisplay()

	// Keep the program running
	select {}
}

// selectCell handles cell selection from JavaScript
func (w *WebUI) selectCell(this js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		return nil
	}

	row := args[0].Int()
	col := args[1].Int()

	if row >= 0 && row < 9 && col >= 0 && col < 9 {
		w.selectedRow = row
		w.selectedCol = col
		w.updateDisplay()
	}

	return nil
}

// setNumber places a number in the selected cell
func (w *WebUI) setNumber(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return nil
	}

	number := args[0].Int()
	if number < 1 || number > 9 {
		return nil
	}

	// Check if the cell is editable (not part of original puzzle)
	if w.originalPuzzle.GetCell(w.selectedRow, w.selectedCol) != 0 {
		w.showMessage("Cannot modify original puzzle numbers!")
		return nil
	}

	// Check if the move is valid
	if !w.sudoku.IsValidMove(w.selectedRow, w.selectedCol, number) {
		w.showMessage("Invalid move! Number already exists in row, column, or box.")
		return nil
	}

	// Place the number
	w.sudoku.SetCell(w.selectedRow, w.selectedCol, number)
	w.updateDisplay()

	// Check if puzzle is solved
	if w.sudoku.IsComplete() {
		w.gameWon = true
		w.showMessage("Congratulations! You solved the puzzle!")
	}

	return nil
}

// clearCell clears the selected cell
func (w *WebUI) clearCell(this js.Value, args []js.Value) interface{} {
	// Check if the cell is editable (not part of original puzzle)
	if w.originalPuzzle.GetCell(w.selectedRow, w.selectedCol) != 0 {
		w.showMessage("Cannot modify original puzzle numbers!")
		return nil
	}

	w.sudoku.ClearCell(w.selectedRow, w.selectedCol)
	w.updateDisplay()
	return nil
}

// newGame generates a new puzzle
func (w *WebUI) newGame(this js.Value, args []js.Value) interface{} {
	difficulty := game.Medium
	if len(args) > 0 {
		diffStr := args[0].String()
		switch diffStr {
		case "easy":
			difficulty = game.Easy
		case "medium":
			difficulty = game.Medium
		case "hard":
			difficulty = game.Hard
		case "expert":
			difficulty = game.Expert
		}
	}

	generator := game.NewGenerator()
	newPuzzle := generator.GeneratePuzzle(difficulty)

	// Update both current and original puzzles
	w.sudoku = newPuzzle
	w.originalPuzzle = &game.Sudoku{}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			w.originalPuzzle.Board[i][j] = newPuzzle.Board[i][j]
		}
	}

	w.gameWon = false
	w.selectedRow = 4
	w.selectedCol = 4
	w.updateDisplay()

	return nil
}

// setDifficulty handles difficulty changes
func (w *WebUI) setDifficulty(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return nil
	}

	w.newGame(this, args)
	return nil
}

// getSudokuState returns the current game state as JSON
func (w *WebUI) getSudokuState(this js.Value, args []js.Value) interface{} {
	state := map[string]interface{}{
		"board":         w.sudoku.Board,
		"originalBoard": w.originalPuzzle.Board,
		"selectedRow":   w.selectedRow,
		"selectedCol":   w.selectedCol,
		"gameWon":       w.gameWon,
	}

	jsonBytes, err := json.Marshal(state)
	if err != nil {
		return nil
	}

	return string(jsonBytes)
}

// updateDisplay updates the HTML display
func (w *WebUI) updateDisplay() {
	// Update each cell in the HTML table
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			cellId := fmt.Sprintf("cell-%d-%d", row, col)
			cell := js.Global().Get("document").Call("getElementById", cellId)
			if cell.IsNull() {
				continue
			}

			value := w.sudoku.GetCell(row, col)
			isOriginal := w.originalPuzzle.GetCell(row, col) != 0
			isSelected := row == w.selectedRow && col == w.selectedCol

			// Set cell content
			if value == 0 {
				cell.Set("textContent", "")
			} else {
				cell.Set("textContent", strconv.Itoa(value))
			}

			// Set cell classes
			classList := cell.Get("classList")
			classList.Call("remove", "original", "user-input", "selected", "invalid")

			if isOriginal {
				classList.Call("add", "original")
			} else if value != 0 {
				classList.Call("add", "user-input")
			}

			if isSelected {
				classList.Call("add", "selected")
			}

			// Check for conflicts
			if value != 0 {
				// Temporarily clear the cell to check if the current value is valid
				originalValue := w.sudoku.GetCell(row, col)
				w.sudoku.ClearCell(row, col)
				isValid := w.sudoku.IsValidMove(row, col, originalValue)
				w.sudoku.SetCell(row, col, originalValue)

				if !isValid {
					classList.Call("add", "invalid")
				}
			}
		}
	}

	// Update selected coordinates display
	coordDisplay := js.Global().Get("document").Call("getElementById", "coordinates")
	if !coordDisplay.IsNull() {
		coord := fmt.Sprintf("%c%d", 'A'+w.selectedCol, w.selectedRow+1)
		coordDisplay.Set("textContent", coord)
	}

	// Update game status
	statusDisplay := js.Global().Get("document").Call("getElementById", "game-status")
	if !statusDisplay.IsNull() {
		if w.gameWon {
			statusDisplay.Set("textContent", "🎉 Puzzle Solved! 🎉")
			statusDisplay.Get("style").Set("color", "#00ff00")
		} else {
			emptyCount := game.CountEmptyCells(w.sudoku)
			statusDisplay.Set("textContent", fmt.Sprintf("Empty cells: %d", emptyCount))
			statusDisplay.Get("style").Set("color", "#ffffff")
		}
	}
}

// showMessage displays a temporary message to the user
func (w *WebUI) showMessage(message string) {
	messageDiv := js.Global().Get("document").Call("getElementById", "message")
	if !messageDiv.IsNull() {
		messageDiv.Set("textContent", message)
		messageDiv.Get("style").Set("display", "block")

		// Hide message after 3 seconds
		js.Global().Call("setTimeout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			messageDiv.Get("style").Set("display", "none")
			return nil
		}), 3000)
	}
}
