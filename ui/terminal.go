package ui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sudoku/game"

	"github.com/eiannone/keyboard"
)

// ANSI color codes for retro Game Boy-style colors
const (
	// Game Boy green palette
	ColorReset = "\033[0m"
	ColorBold  = "\033[1m"
	ColorDim   = "\033[2m"

	// Background colors (Game Boy style)
	BgGreen      = "\033[42m"
	BgDarkGreen  = "\033[48;5;22m"
	BgLightGreen = "\033[48;5;156m"

	// Text colors
	ColorBlack   = "\033[30m"
	ColorWhite   = "\033[37m"
	ColorBlue    = "\033[34m"
	ColorRed     = "\033[31m"
	ColorYellow  = "\033[33m"
	ColorGreen   = "\033[32m"
	ColorCyan    = "\033[36m"
	ColorMagenta = "\033[35m"

	// Bright colors
	ColorBrightBlue    = "\033[94m"
	ColorBrightRed     = "\033[91m"
	ColorBrightGreen   = "\033[92m"
	ColorBrightCyan    = "\033[96m"
	ColorBrightYellow  = "\033[93m"
	ColorBrightMagenta = "\033[95m"

	// Cursor positioning
	CursorHide  = "\033[?25l"
	CursorShow  = "\033[?25h"
	ClearScreen = "\033[2J"
	CursorHome  = "\033[H"
)

// Tool represents the current editing tool
type Tool int

const (
	ToolNumber Tool = iota // Blue tool - place numbers
	ToolBlock              // Red tool - mark impossible numbers
	ToolMove               // Yellow tool - just move cursor
)

// TerminalUI represents the Game Boy-style terminal interface
type TerminalUI struct {
	sudoku        *game.Sudoku
	originalBoard [9][9]int      // Keep track of original puzzle numbers
	userNumbers   [9][9]int      // User-placed numbers (blue)
	blockedNums   [9][9][10]bool // blocked[row][col][num] - numbers marked as impossible (red)

	cursorRow   int
	cursorCol   int
	currentTool Tool
	selectedNum int // Current number (1-9) for tools

	gameRunning bool
}

// NewTerminalUI creates a new terminal UI
func NewTerminalUI(sudoku *game.Sudoku) *TerminalUI {
	ui := &TerminalUI{
		sudoku:      sudoku,
		cursorRow:   4, // Start in center
		cursorCol:   4,
		currentTool: ToolMove,
		selectedNum: 1,
		gameRunning: true,
	}

	// Copy original board
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			ui.originalBoard[i][j] = sudoku.GetCell(i, j)
		}
	}

	return ui
}

// clearScreen clears the terminal screen
func (ui *TerminalUI) clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		fmt.Print(ClearScreen + CursorHome)
	}
}

// Run starts the terminal UI game loop
func (ui *TerminalUI) Run() {
	// Initialize keyboard
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Print(CursorHide)       // Hide cursor
	defer fmt.Print(CursorShow) // Show cursor when done

	for ui.gameRunning {
		ui.clearScreen()
		ui.draw()
		ui.handleInput()
	}
}

// draw renders the entire game interface
func (ui *TerminalUI) draw() {
	ui.drawTitle()
	ui.drawBoard()
	ui.drawToolbar()
	ui.drawControls()
	ui.drawStatus()
}

// drawTitle draws the Game Boy-style title
func (ui *TerminalUI) drawTitle() {
	fmt.Printf("%s%s", ColorBrightGreen, ColorBold)
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    🔢 RETRO SUDOKU 🔢                        ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Printf("%s\n", ColorReset)
}

// drawBoard renders the Sudoku board with Game Boy styling
func (ui *TerminalUI) drawBoard() {
	fmt.Printf("%s%s", ColorBrightGreen, ColorBold)
	fmt.Println("    A  B  C   D  E  F   G  H  I")
	fmt.Printf("%s", ColorReset)

	for row := 0; row < 9; row++ {
		// Draw horizontal separators
		if row == 3 || row == 6 {
			fmt.Printf("%s%s", ColorBrightGreen, ColorBold)
			fmt.Println("  ╠═════════╬═════════╬═════════╣")
			fmt.Printf("%s", ColorReset)
		} else if row == 0 {
			fmt.Printf("%s%s", ColorBrightGreen, ColorBold)
			fmt.Println("  ╔═════════╦═════════╦═════════╗")
			fmt.Printf("%s", ColorReset)
		}

		// Draw row number
		fmt.Printf("%s%s%d%s ", ColorBrightGreen, ColorBold, row+1, ColorReset)

		// Draw board row
		for col := 0; col < 9; col++ {
			// Draw vertical separators - only thick borders between 3x3 boxes
			if col == 0 {
				fmt.Printf("%s%s║%s", ColorBrightGreen, ColorBold, ColorReset)
			} else if col == 3 || col == 6 {
				fmt.Printf("%s%s║%s", ColorBrightGreen, ColorBold, ColorReset)
			}
			// No separator for other columns (removes the thin │ lines)

			ui.drawCell(row, col)
		}
		fmt.Printf("%s%s║%s\n", ColorBrightGreen, ColorBold, ColorReset)
	}

	// Bottom border
	fmt.Printf("%s%s", ColorBrightGreen, ColorBold)
	fmt.Println("  ╚═════════╩═════════╩═════════╝")
	fmt.Printf("%s\n", ColorReset)
}

// drawCell renders a single cell with appropriate styling
func (ui *TerminalUI) drawCell(row, col int) {
	isCursor := (row == ui.cursorRow && col == ui.cursorCol)
	originalNum := ui.originalBoard[row][col]
	userNum := ui.userNumbers[row][col]

	// Determine what to display
	var displayChar string
	var cellColor string

	if originalNum != 0 {
		// Original puzzle number (black/dark)
		displayChar = fmt.Sprintf("%d", originalNum)
		cellColor = ColorBlack + ColorBold
	} else if userNum != 0 {
		// User-placed number (blue)
		displayChar = fmt.Sprintf("%d", userNum)
		cellColor = ColorBrightBlue + ColorBold
	} else {
		// Empty cell
		displayChar = " "
		cellColor = ColorReset
	}

	// Add cursor highlight
	if isCursor {
		switch ui.currentTool {
		case ToolNumber:
			fmt.Printf("%s%s %s %s", BgDarkGreen, ColorBrightBlue, displayChar, ColorReset)
		case ToolBlock:
			fmt.Printf("%s%s %s %s", BgDarkGreen, ColorBrightRed, displayChar, ColorReset)
		case ToolMove:
			fmt.Printf("%s%s %s %s", BgLightGreen, ColorBlack, displayChar, ColorReset)
		}
	} else {
		fmt.Printf(" %s%s%s ", cellColor, displayChar, ColorReset)
	}
}

// drawToolbar renders the Game Boy-style toolbar with vertical tool layout
func (ui *TerminalUI) drawToolbar() {
	fmt.Printf("%s%s╔═══════════════════ TOOLS ═══════════════════╗%s\n", ColorBrightGreen, ColorBold, ColorReset)

	// MOVE tool
	if ui.currentTool == ToolMove {
		fmt.Printf("║ %s%s▶▶ MOVE CURSOR ◀◀%s                           ║\n", ColorBold, ColorYellow, ColorReset)
	} else {
		fmt.Printf("║      Move Cursor                            ║\n")
	}

	// BLUE PENCIL tool (NUMBER)
	if ui.currentTool == ToolNumber {
		fmt.Printf("║ %s%s▶▶ BLUE PENCIL (MARK MAYBE) ◀◀%s         ║\n", ColorBold, ColorBrightBlue, ColorReset)
	} else {
		fmt.Printf("║   🔢 Blue Pencil (mark maybe)               ║\n")
	}

	// RED PENCIL tool (BLOCK)
	if ui.currentTool == ToolBlock {
		fmt.Printf("║ %s%s▶▶ RED PENCIL (MARK NOT POSSIBLE) ◀◀%s   ║\n", ColorBold, ColorBrightRed, ColorReset)
	} else {
		fmt.Printf("║   🚫 Red Pencil (mark not possible)         ║\n")
	}

	// CLEAR option
	fmt.Printf("║   🧽 Clear Marks at Cursor (C)              ║\n")
	fmt.Printf("╬═════════════════════════════════════════════╬\n")

	// Always show numbers section
	fmt.Print("║ Numbers: ")
	for i := 1; i <= 9; i++ {
		if i == ui.selectedNum {
			fmt.Printf("%s%s[%d]%s ", ColorBold, ColorBrightGreen, i, ColorReset)
		} else {
			fmt.Printf(" %s%d%s ", ColorDim, i, ColorReset)
		}
	}
	fmt.Println("       ║")

	fmt.Printf("%s%s╚═════════════════════════════════════════════╝%s\n", ColorBrightGreen, ColorBold, ColorReset)
}

// drawControls shows the control scheme
func (ui *TerminalUI) drawControls() {
	fmt.Printf("%s%s╔══════════════ CONTROLS ═══════════════════════╗%s\n",
		ColorBrightGreen, ColorBold, ColorReset)
	fmt.Printf("║ %sWASD%s/arrows: D-Pad (Move cursor)              ║\n", ColorBrightCyan, ColorReset)
	fmt.Printf("║ %sQ/E%s: Shoulder buttons (Switch tools)          ║\n", ColorBrightCyan, ColorReset)
	fmt.Printf("║ %sF%s: Select button (Cycle numbers 1-9)          ║\n", ColorBrightCyan, ColorReset)
	fmt.Printf("║ %sSPACE%s: A button (Use current tool)            ║\n", ColorBrightCyan, ColorReset)
	fmt.Printf("║ %sC%s: B button (Clear marks at cursor)           ║\n", ColorBrightCyan, ColorReset)
	fmt.Printf("║ %sESC%s: Quit game                                ║\n", ColorBrightCyan, ColorReset)
	fmt.Printf("%s%s╚═══════════════════════════════════════════════╝%s\n", ColorBrightGreen, ColorBold, ColorReset)
}

// drawStatus shows current game status
func (ui *TerminalUI) drawStatus() {
	fmt.Printf("Cursor: %s%s%c%d%s  Tool: %s", ColorBrightYellow, ColorBold,
		'A'+ui.cursorCol, ui.cursorRow+1, ColorReset, ui.getToolName())

	if ui.currentTool != ToolMove {
		fmt.Printf("  Selected: %s%s%d%s", ColorBrightMagenta, ColorBold, ui.selectedNum, ColorReset)
	}

	// Show blocked numbers for current cell if any
	if ui.hasBlockedNumbers(ui.cursorRow, ui.cursorCol) {
		fmt.Printf("  Blocked: %s", ui.getBlockedNumbersString(ui.cursorRow, ui.cursorCol))
	}

	fmt.Println()
}

// getToolName returns the current tool name with color
func (ui *TerminalUI) getToolName() string {
	switch ui.currentTool {
	case ToolMove:
		return fmt.Sprintf("%s%sMOVE%s", ColorYellow, ColorBold, ColorReset)
	case ToolNumber:
		return fmt.Sprintf("%s%sNUMBER%s", ColorBrightBlue, ColorBold, ColorReset)
	case ToolBlock:
		return fmt.Sprintf("%s%sBLOCK%s", ColorBrightRed, ColorBold, ColorReset)
	default:
		return "UNKNOWN"
	}
}

// hasBlockedNumbers checks if a cell has any blocked numbers
func (ui *TerminalUI) hasBlockedNumbers(row, col int) bool {
	for i := 1; i <= 9; i++ {
		if ui.blockedNums[row][col][i] {
			return true
		}
	}
	return false
}

// getBlockedNumbersString returns a string of blocked numbers for a cell
func (ui *TerminalUI) getBlockedNumbersString(row, col int) string {
	var blocked []string
	for i := 1; i <= 9; i++ {
		if ui.blockedNums[row][col][i] {
			blocked = append(blocked, fmt.Sprintf("%s%s%d%s", ColorBrightRed, ColorBold, i, ColorReset))
		}
	}
	if len(blocked) == 0 {
		return "none"
	}
	return fmt.Sprintf("[%s]", fmt.Sprintf("%v", blocked))
}

// handleInput processes user input with immediate keypress detection
func (ui *TerminalUI) handleInput() {
	fmt.Printf("\n%s%sPress any key (WASD/QE/F/SPACE/C/ESC)...%s", ColorBrightCyan, ColorBold, ColorReset)

	char, key, err := keyboard.GetSingleKey()
	if err != nil {
		return
	}

	// Handle special keys
	switch key {
	case keyboard.KeySpace:
		ui.useTool()
	case keyboard.KeyEsc:
		ui.gameRunning = false
	case keyboard.KeyArrowUp:
		ui.moveCursor(-1, 0)
	case keyboard.KeyArrowDown:
		ui.moveCursor(1, 0)
	case keyboard.KeyArrowLeft:
		ui.moveCursor(0, -1)
	case keyboard.KeyArrowRight:
		ui.moveCursor(0, 1)
	default:
		// Handle character keys
		switch char {
		// D-Pad movement (WASD)
		case 'w', 'W':
			ui.moveCursor(-1, 0)
		case 's', 'S':
			ui.moveCursor(1, 0)
		case 'a', 'A':
			ui.moveCursor(0, -1)
		case 'd', 'D':
			ui.moveCursor(0, 1)

		// Shoulder buttons - tool switching
		case 'q', 'Q':
			ui.cycleTool()
		case 'e', 'E':
			ui.cycleToolReverse()

		// Select button - number cycling
		case 'f', 'F':
			ui.cycleNumber()

		// B button - clear cell
		case 'c', 'C':
			ui.clearCurrentCell()

		// Direct number selection (1-9)
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			ui.selectedNum = int(char - '0')
		}
	}
} // moveCursor moves the cursor in the specified direction
func (ui *TerminalUI) moveCursor(deltaRow, deltaCol int) {
	newRow := ui.cursorRow + deltaRow
	newCol := ui.cursorCol + deltaCol

	if newRow >= 0 && newRow < 9 && newCol >= 0 && newCol < 9 {
		ui.cursorRow = newRow
		ui.cursorCol = newCol
	}
}

// cycleTool cycles through the available tools (forward)
func (ui *TerminalUI) cycleTool() {
	switch ui.currentTool {
	case ToolMove:
		ui.currentTool = ToolNumber
	case ToolNumber:
		ui.currentTool = ToolBlock
	case ToolBlock:
		ui.currentTool = ToolMove
	}
}

// cycleToolReverse cycles through the available tools (backward)
func (ui *TerminalUI) cycleToolReverse() {
	switch ui.currentTool {
	case ToolMove:
		ui.currentTool = ToolBlock
	case ToolNumber:
		ui.currentTool = ToolMove
	case ToolBlock:
		ui.currentTool = ToolNumber
	}
}

// cycleNumber cycles through numbers 1-9
func (ui *TerminalUI) cycleNumber() {
	ui.selectedNum++
	if ui.selectedNum > 9 {
		ui.selectedNum = 1
	}
}

// useTool applies the current tool action
func (ui *TerminalUI) useTool() {
	row, col := ui.cursorRow, ui.cursorCol

	// Can't modify original puzzle numbers
	if ui.originalBoard[row][col] != 0 {
		return
	}

	switch ui.currentTool {
	case ToolNumber:
		// Place number (blue)
		ui.userNumbers[row][col] = ui.selectedNum
		// Clear any blocked status for this number
		ui.blockedNums[row][col][ui.selectedNum] = false

	case ToolBlock:
		// Toggle blocked number (red)
		ui.blockedNums[row][col][ui.selectedNum] = !ui.blockedNums[row][col][ui.selectedNum]
		// If we're blocking the current user number, clear it
		if ui.blockedNums[row][col][ui.selectedNum] && ui.userNumbers[row][col] == ui.selectedNum {
			ui.userNumbers[row][col] = 0
		}
	}
}

// clearCurrentCell clears the current cell
func (ui *TerminalUI) clearCurrentCell() {
	row, col := ui.cursorRow, ui.cursorCol

	// Can't clear original puzzle numbers
	if ui.originalBoard[row][col] != 0 {
		return
	}

	// Clear user number and all blocked numbers
	ui.userNumbers[row][col] = 0
	for i := 1; i <= 9; i++ {
		ui.blockedNums[row][col][i] = false
	}
}
