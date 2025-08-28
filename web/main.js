// WebAssembly and JavaScript bridge for Sudoku game

let wasmModule = null;
let selectedCell = { row: 4, col: 4 };

// Initialize the game when the page loads
document.addEventListener('DOMContentLoaded', function() {
    initializeSudokuBoard();
    setupKeyboardControls();
    loadWasm();
});

// Create the Sudoku board HTML structure
function initializeSudokuBoard() {
    const board = document.getElementById('sudoku-board');
    board.innerHTML = '';
    
    for (let row = 0; row < 9; row++) {
        const tr = document.createElement('tr');
        for (let col = 0; col < 9; col++) {
            const td = document.createElement('td');
            td.id = `cell-${row}-${col}`;
            td.onclick = () => selectCell(row, col);
            tr.appendChild(td);
        }
        board.appendChild(tr);
    }
}

// Load WebAssembly module
async function loadWasm() {
    const go = new Go();
    
    try {
        const result = await WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject);
        go.run(result.instance);
        wasmModule = result.instance;
        
        console.log("WebAssembly module loaded successfully");
        
        // Start with a medium difficulty game
        newGame('medium');
    } catch (error) {
        console.error("Failed to load WebAssembly module:", error);
        showFallbackMessage();
    }
}

// Fallback message if WASM fails to load
function showFallbackMessage() {
    const container = document.querySelector('.container');
    container.innerHTML = `
        <div style="text-align: center; padding: 50px;">
            <h2 style="color: var(--retro-error);">⚠️ WebAssembly Loading Failed</h2>
            <p>Unable to load the Sudoku game engine.</p>
            <p>Please ensure you're serving this page from a web server and that main.wasm is available.</p>
            <br>
            <p><strong>To build the WebAssembly module:</strong></p>
            <code style="background: var(--retro-dark); padding: 10px; display: block; margin: 10px;">
                GOOS=js GOARCH=wasm go build -o web/main.wasm cmd/wasm/main.go
            </code>
        </div>
    `;
}

// Set up keyboard event listeners
function setupKeyboardControls() {
    document.addEventListener('keydown', function(event) {
        switch (event.key) {
            case 'ArrowUp':
                event.preventDefault();
                moveSelection(-1, 0);
                break;
            case 'ArrowDown':
                event.preventDefault();
                moveSelection(1, 0);
                break;
            case 'ArrowLeft':
                event.preventDefault();
                moveSelection(0, -1);
                break;
            case 'ArrowRight':
                event.preventDefault();
                moveSelection(0, 1);
                break;
            case '1':
            case '2':
            case '3':
            case '4':
            case '5':
            case '6':
            case '7':
            case '8':
            case '9':
                event.preventDefault();
                setNumber(parseInt(event.key));
                break;
            case 'Delete':
            case 'Backspace':
                event.preventDefault();
                clearCell();
                break;
            case 'Escape':
                event.preventDefault();
                // Could add a pause/menu feature here
                break;
        }
    });
}

// Move selection with bounds checking
function moveSelection(deltaRow, deltaCol) {
    const newRow = Math.max(0, Math.min(8, selectedCell.row + deltaRow));
    const newCol = Math.max(0, Math.min(8, selectedCell.col + deltaCol));
    selectCell(newRow, newCol);
}

// These functions will be called by the Go WASM module
// They're defined globally so Go can access them

// Select a cell (called from Go and JavaScript)
function selectCell(row, col) {
    if (row < 0 || row > 8 || col < 0 || col > 8) return;
    
    selectedCell.row = row;
    selectedCell.col = col;
    
    // Remove previous selection
    document.querySelectorAll('.sudoku-board td').forEach(cell => {
        cell.classList.remove('selected');
    });
    
    // Add selection to new cell
    const cell = document.getElementById(`cell-${row}-${col}`);
    if (cell) {
        cell.classList.add('selected');
    }
    
    // Update coordinates display
    const coordDisplay = document.getElementById('coordinates');
    if (coordDisplay) {
        const coord = String.fromCharCode(65 + col) + (row + 1);
        coordDisplay.textContent = coord;
    }
    
    // Call Go function if available
    if (window.selectCell && typeof window.selectCell === 'function') {
        window.selectCell(row, col);
    }
}

// Set a number in the selected cell
function setNumber(number) {
    if (number < 1 || number > 9) return;
    
    // Call Go function if available
    if (window.setNumber && typeof window.setNumber === 'function') {
        window.setNumber(number);
    } else {
        // Fallback for testing without WASM
        console.log(`Setting number ${number} at ${selectedCell.row}, ${selectedCell.col}`);
    }
}

// Clear the selected cell
function clearCell() {
    // Call Go function if available
    if (window.clearCell && typeof window.clearCell === 'function') {
        window.clearCell();
    } else {
        // Fallback for testing without WASM
        console.log(`Clearing cell at ${selectedCell.row}, ${selectedCell.col}`);
    }
}

// Start a new game
function newGame(difficulty = 'medium') {
    // Call Go function if available
    if (window.newGame && typeof window.newGame === 'function') {
        window.newGame(difficulty);
    } else {
        // Fallback for testing without WASM
        console.log(`Starting new ${difficulty} game`);
        showMessage(`New ${difficulty} game would start here!`);
    }
}

// Set difficulty and start new game
function setDifficulty(difficulty) {
    newGame(difficulty);
}

// Show a temporary message
function showMessage(text) {
    const messageDiv = document.getElementById('message');
    if (messageDiv) {
        messageDiv.textContent = text;
        messageDiv.style.display = 'block';
        
        setTimeout(() => {
            messageDiv.style.display = 'none';
        }, 3000);
    }
}

// Get game state (for debugging or save/load features)
function getSudokuState() {
    if (window.getSudokuState && typeof window.getSudokuState === 'function') {
        return window.getSudokuState();
    }
    return null;
}

// Update difficulty selector when page loads
document.addEventListener('DOMContentLoaded', function() {
    const difficultySelect = document.getElementById('difficulty');
    if (difficultySelect) {
        difficultySelect.value = 'medium';
    }
});

// Add visual feedback for button presses
document.addEventListener('DOMContentLoaded', function() {
    const numberButtons = document.querySelectorAll('.number-btn');
    numberButtons.forEach(button => {
        button.addEventListener('click', function() {
            this.style.transform = 'scale(0.95)';
            setTimeout(() => {
                this.style.transform = 'scale(1)';
            }, 100);
        });
    });
});

// Handle page visibility changes (pause game when tab is hidden)
document.addEventListener('visibilitychange', function() {
    if (document.hidden) {
        // Game is hidden, could pause here if needed
        console.log('Game paused');
    } else {
        // Game is visible again
        console.log('Game resumed');
    }
});

// Add touch support for mobile devices
let touchStartX = 0;
let touchStartY = 0;

document.addEventListener('touchstart', function(e) {
    touchStartX = e.touches[0].clientX;
    touchStartY = e.touches[0].clientY;
}, { passive: true });

document.addEventListener('touchend', function(e) {
    if (!touchStartX || !touchStartY) return;
    
    const touchEndX = e.changedTouches[0].clientX;
    const touchEndY = e.changedTouches[0].clientY;
    
    const deltaX = touchEndX - touchStartX;
    const deltaY = touchEndY - touchStartY;
    
    const minSwipeDistance = 50;
    
    if (Math.abs(deltaX) > minSwipeDistance || Math.abs(deltaY) > minSwipeDistance) {
        if (Math.abs(deltaX) > Math.abs(deltaY)) {
            // Horizontal swipe
            if (deltaX > 0) {
                moveSelection(0, 1); // Right
            } else {
                moveSelection(0, -1); // Left
            }
        } else {
            // Vertical swipe
            if (deltaY > 0) {
                moveSelection(1, 0); // Down
            } else {
                moveSelection(-1, 0); // Up
            }
        }
    }
    
    touchStartX = 0;
    touchStartY = 0;
}, { passive: true });
