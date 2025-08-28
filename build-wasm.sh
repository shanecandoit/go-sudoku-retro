#!/bin/bash

# Build script for Sudoku WebAssembly

echo "🚀 Building Sudoku WebAssembly..."

# Set environment variables for WebAssembly build
export GOOS=js
export GOARCH=wasm

# Build the WebAssembly module
echo "📦 Compiling Go to WebAssembly..."
go build -o web/main.wasm cmd/wasm/main.go

if [ $? -eq 0 ]; then
    echo "✅ WebAssembly build successful!"
    
    # Copy wasm_exec.js from Go installation
    echo "📋 Copying wasm_exec.js..."
    cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" web/
    
    if [ $? -eq 0 ]; then
        echo "✅ wasm_exec.js copied successfully!"
        echo ""
        echo "🎮 Sudoku WebAssembly is ready!"
        echo "📁 Files created in web/ directory:"
        echo "   - main.wasm (WebAssembly module)"
        echo "   - wasm_exec.js (Go WebAssembly runtime)"
        echo "   - index.html (Game interface)"
        echo "   - styles.css (Retro styling)"
        echo "   - main.js (JavaScript bridge)"
        echo ""
        echo "🌐 To run the game:"
        echo "   1. Start a local web server in the web/ directory"
        echo "   2. Open index.html in your browser"
        echo ""
        echo "💡 Example web server commands:"
        echo "   Python 3: python -m http.server 8000"
        echo "   Node.js:  npx http-server"
        echo "   Go:       go run -m http.FileServer ."
    else
        echo "❌ Failed to copy wasm_exec.js"
        echo "💡 You may need to copy it manually from:"
        echo "   $(go env GOROOT)/misc/wasm/wasm_exec.js"
    fi
else
    echo "❌ WebAssembly build failed!"
    echo "💡 Make sure you have Go installed and your code compiles correctly"
fi
