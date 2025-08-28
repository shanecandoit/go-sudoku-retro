package main

import (
	"sudoku/game"
	"sudoku/ui/web"
)

func main() {
	// Create a default medium difficulty puzzle
	generator := game.NewGenerator()
	puzzle := generator.GeneratePuzzle(game.Medium)

	// Create and start the web UI
	webUI := web.NewWebUI(puzzle)
	webUI.Run()
}
