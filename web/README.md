# 🎮 Sudoku WebAssembly

A retro-styled Sudoku game built with Go and WebAssembly, featuring a Game Boy-inspired green color scheme and smooth browser-based gameplay.

## 🚀 Quick Start

### Prerequisites
- Go 1.21+ installed
- A modern web browser with WebAssembly support

### Building and Running

1. **Build the WebAssembly module:**
   ```bash
   # On Windows
   ./build-wasm.bat
   
   # On Linux/Mac
   ./build-wasm.sh
   ```

2. **Start a local web server:**
   ```bash
   cd web
   python -m http.server 8000
   ```

3. **Open your browser and navigate to:**
   ```
   http://localhost:8000
   ```

## 🎯 Features

### Game Features
- **Multiple Difficulty Levels**: Easy, Medium, Hard, Expert
- **Smart Validation**: Invalid moves are highlighted in red
- **Visual Feedback**: Different colors for original vs user-entered numbers
- **Win Detection**: Automatic puzzle completion detection

### User Interface
- **Retro Aesthetic**: Game Boy-inspired green color scheme
- **Responsive Design**: Works on desktop and mobile devices
- **Multiple Input Methods**:
  - Click cells to select
  - Number buttons or keyboard (1-9)
  - Arrow keys for navigation
  - Touch/swipe support on mobile

### Controls
- **Mouse/Touch**: Click cells and number buttons
- **Keyboard**:
  - Arrow keys: Navigate the board
  - 1-9: Place numbers
  - Delete/Backspace: Clear cells
  - Escape: (Reserved for future features)

## 📁 File Structure

```
web/
├── index.html      # Main game interface
├── styles.css      # Retro Game Boy styling
├── main.js         # JavaScript bridge and UI logic
├── main.wasm       # WebAssembly module (generated)
└── wasm_exec.js    # Go WebAssembly runtime (copied from Go installation)

ui/web/
└── web.go          # Go WebAssembly interface

cmd/wasm/
└── main.go         # WebAssembly entry point
```

## 🛠️ Development

### Manual Build Commands

If the build scripts don't work, you can build manually:

```bash
# Set environment variables
export GOOS=js GOARCH=wasm  # Linux/Mac
set GOOS=js && set GOARCH=wasm  # Windows

# Build WebAssembly module
go build -o web/main.wasm cmd/wasm/main.go

# Copy WebAssembly runtime
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" web/  # Linux/Mac
copy "%GOROOT%\misc\wasm\wasm_exec.js" web\       # Windows
```

### Adding Features

The WebAssembly interface is in `ui/web/web.go`. Key functions:

- `selectCell()`: Handle cell selection
- `setNumber()`: Place numbers with validation
- `clearCell()`: Clear user-entered numbers
- `newGame()`: Generate new puzzles
- `updateDisplay()`: Refresh the UI

## 🎨 Customization

### Colors
The retro color scheme is defined in CSS variables at the top of `styles.css`:

```css
:root {
    --retro-dark: #0f380f;      /* Dark green */
    --retro-medium: #306230;    /* Medium green */
    --retro-light: #9bbc0f;     /* Light green */
    --retro-accent: #00ff41;    /* Bright green accent */
    --retro-error: #ff4444;     /* Error red */
}
```

### Board Size
To change the visual size, modify the cell dimensions in `styles.css`:

```css
.sudoku-board td {
    width: 45px;   /* Change width */
    height: 45px;  /* Change height */
    font-size: 1.4em;  /* Adjust font size */
}
```

## 🐛 Troubleshooting

### WebAssembly Module Failed to Load
- Ensure you're serving the files from a web server (not `file://`)
- Check that `main.wasm` was built successfully
- Verify `wasm_exec.js` is in the web directory
- Check browser developer console for specific errors

### Build Errors
- Ensure Go 1.21+ is installed
- Check that all Go modules are available (`go mod tidy`)
- Verify the game code compiles normally (`go build`)

### Performance Issues
- WebAssembly runs efficiently, but initial load may be slow
- Large difficulty changes regenerate the entire puzzle
- Browser developer tools can help profile performance

## 🌐 Browser Compatibility

Tested and working on:
- Chrome 90+
- Firefox 89+
- Safari 14+
- Edge 90+

WebAssembly is supported in all modern browsers. For older browsers, consider showing a compatibility message.

## 📱 Mobile Support

The game includes touch controls:
- Tap cells to select
- Swipe to navigate
- Responsive design adapts to screen size
- Number pad is touch-friendly

## 🚀 Deployment

For production deployment:

1. Build the WebAssembly module
2. Upload all files in `web/` directory to your web server
3. Configure proper MIME types:
   - `.wasm` files: `application/wasm`
   - `.js` files: `application/javascript`

## 📄 License

This project is part of the Sudoku game implementation. See the main project README for license information.
