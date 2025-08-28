#!/bin/bash

# Simple web server script for Sudoku WebAssembly

echo "🌐 Starting web server for Sudoku WebAssembly..."

# Check if we're in the right directory
if [ ! -f "index.html" ]; then
    echo "❌ index.html not found. Make sure you're in the web/ directory."
    echo "💡 Run this script from the web/ directory, or cd to web/ first."
    exit 1
fi

# Check if main.wasm exists
if [ ! -f "main.wasm" ]; then
    echo "⚠️  main.wasm not found. Building WebAssembly module..."
    cd ..
    ./build-wasm.sh
    cd web
fi

# Try different web server options
PORT=8000

echo "🚀 Attempting to start web server on port $PORT..."

# Try Python 3 first
if command -v python3 &> /dev/null; then
    echo "🐍 Using Python 3..."
    echo "🌐 Open your browser and go to: http://localhost:$PORT"
    echo "🛑 Press Ctrl+C to stop the server"
    python3 -m http.server $PORT
elif command -v npx &> /dev/null; then
    echo "📦 Using Node.js http-server..."
    echo "🌐 Server will start automatically in your browser"
    echo "🛑 Press Ctrl+C to stop the server"
    npx http-server -p $PORT -o
elif command -v go &> /dev/null; then
    echo "🐹 Using Go simple file server..."
    echo "🌐 Open your browser and go to: http://localhost:$PORT"
    echo "🛑 Press Ctrl+C to stop the server"
    cat > temp_server.go << 'EOF'
package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    fs := http.FileServer(http.Dir("."))
    http.Handle("/", fs)
    
    fmt.Println("Server starting on :8000...")
    log.Fatal(http.ListenAndServe(":8000", nil))
}
EOF
    go run temp_server.go
    rm temp_server.go
else
    echo "❌ No suitable web server found!"
    echo "💡 Please install one of the following:"
    echo "   - Python: python -m http.server 8000"
    echo "   - Node.js: npx http-server"
    echo "   - Or use any other static file server"
    exit 1
fi
