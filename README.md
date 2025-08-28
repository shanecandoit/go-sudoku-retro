# Sudoku Game in Go

A modern Sudoku puzzle game implemented in Go, designed to run both as a native application and in web browsers using WebAssembly (WASM) with TinyGo.

## Features

- 🎯 Complete Sudoku puzzle generation and solving
- 🌐 Web browser support via WebAssembly
- 🎮 Interactive gameplay with validation
- 🔍 Hint system and auto-solve functionality
- 📱 Responsive design for desktop and mobile
- ⚡ Fast compilation with TinyGo
- 🎨 Clean, intuitive user interface

## Architecture

This project is structured to support multiple deployment targets:

- **Native Go**: Traditional command-line or GUI application
- **WebAssembly**: Browser-based game using TinyGo compilation

### For WebAssembly Development

- [TinyGo](https://tinygo.org/getting-started/install/)
- A modern web browser with WASM support
- Basic HTTP server (for local development)

## Building for WebAssembly

```bash
# Build WASM binary
tinygo build -o web/sudoku.wasm -target wasm main.go

# Serve the web application
cd web
python -m http.server 8080
# or
npx serve .
```

Then open `http://localhost:8080` in your browser.

## Building for Production

```bash
# Native binary
go build -o sudoku main.go

# Optimized WASM build
tinygo build -o web/sudoku.wasm -target wasm -opt 2 main.go
```

## Project Structure

```text
sudoku/
├── main.go              # Main entry point
├── game/               # Core game logic
│   ├── sudoku.go       # Sudoku puzzle logic
│   ├── generator.go    # Puzzle generation
│   ├── solver.go       # Solving algorithms
│   └── validator.go    # Input validation
├── ui/                 # User interface layers
│   ├── terminal/       # CLI interface
│   ├── web/           # Web/WASM interface
│   └── common/        # Shared UI utilities
├── web/               # Web assets
│   ├── index.html     # Main HTML page
│   ├── style.css      # Styles
│   ├── app.js         # JavaScript integration
│   └── sudoku.wasm    # Compiled WASM binary
├── assets/            # Game assets
└── tests/             # Test files
```

## WebAssembly Integration

The game leverages TinyGo's excellent WebAssembly support to provide a seamless browser experience:

### Key Benefits of TinyGo for WASM

- **Small binary size**: TinyGo produces much smaller WASM files than standard Go
- **Fast compilation**: Quick build times for development
- **Good browser compatibility**: Works across modern browsers
- **Easy JavaScript integration**: Simple interop with web APIs

Happy Sudoku solving! 🧩
