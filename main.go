package main

import (
	"fmt"
	"sudoku/game"
	"sudoku/ui"
)

func main() {
	fmt.Println("� Retro Sudoku - Simple Terminal UI Demo")
	fmt.Println("=================================================")
	fmt.Println()

	// Create a simple puzzle for demonstration
	generator := game.NewGenerator()
	puzzle := generator.GeneratePuzzle(game.Easy) // Use Easy for demo so it's not too hard

	fmt.Println("🧩 GAME FEATURES:")
	fmt.Println("• Original puzzle numbers are shown in BLACK/DARK")
	fmt.Println("• User-placed numbers are shown in BRIGHT BLUE")
	fmt.Println("• Blocked/impossible numbers are marked in RED")
	fmt.Println("• Cursor highlights change color based on active tool")
	fmt.Println()

	fmt.Println("🎯 TOOLS:")
	fmt.Println("• MOVE (Yellow): Navigate around the board")
	fmt.Println("• NUMBER (Blue): Place numbers you think are correct")
	fmt.Println("• BLOCK (Red): Mark numbers that cannot go in a cell")
	fmt.Println()

	fmt.Println("🕹️ RETRO CONTROLS:")
	fmt.Println("• WASD: D-Pad (Move cursor around the 9x9 grid)")
	fmt.Println("• Q/E: Shoulder buttons (Switch tools: Move ↔ Number ↔ Block)")
	fmt.Println("• F: Select button (Cycle through numbers 1-9)")
	fmt.Println("• SPACE: A button (Use the current tool)")
	fmt.Println("• C: B button (Clear current cell)")
	fmt.Println("• ESC: Start button (Quit the game)")
	fmt.Println()

	fmt.Println("🎨 Retro AESTHETIC:")
	fmt.Println("• Retro green color scheme")
	fmt.Println("• ASCII art borders and UI elements")
	fmt.Println("• Clear visual feedback for all actions")
	fmt.Println("• Coordinate system using A-I columns and 1-9 rows")
	fmt.Println()

	fmt.Println("Generated Easy Sudoku Puzzle:")
	puzzle.Display()
	fmt.Printf("Difficulty: Easy (%d empty cells, %.1f%% filled)\n",
		game.CountEmptyCells(puzzle), game.GetFilledRatio(puzzle)*100)
	fmt.Println()

	fmt.Println("Starting the Retro-style Terminal UI...")
	fmt.Println("Press ENTER to start, then follow the on-screen instructions!")
	var input string
	fmt.Scanln(&input)

	// Start the retro terminal UI
	terminalUI := ui.NewTerminalUI(puzzle)
	terminalUI.Run()

	fmt.Println()
	fmt.Println("🎮 Thanks for playing Retro Sudoku!")
	fmt.Println("   Created with Go + TinyGo for WebAssembly support")
	fmt.Println("   Perfect for retro gaming enthusiasts! 🕹️")
}
