# TicTacToeBot

A Go-based Tic-Tac-Toe bot that implements the minimax algorithm for optimal gameplay.

## Project Structure

- `bot.go` - Main bot implementation with game logic and HTTP handlers
- `bot_test.go` - Comprehensive unit tests
- `models/jsonrpc.go` - JSON-RPC data structures
- `go.mod` - Go module definition

## Features

- **Game Logic**: Complete Tic-Tac-Toe game state evaluation
- **AI Algorithm**: Minimax algorithm with parallel processing for optimal moves
- **HTTP API**: JSON-RPC based HTTP server for bot communication
- **Comprehensive Testing**: Full unit test coverage

## Running the Tests

To run the unit tests, you need Go installed on your system:

1. **Install Go** (if not already installed):
   - Download from https://golang.org/dl/
   - Follow installation instructions for your OS

2. **Run all tests**:
   ```bash
   go test -v
   ```

3. **Run specific test functions**:
   ```bash
   go test -v -run TestIsGameOver
   go test -v -run TestMiniMax
   go test -v -run TestMakeBestMove
   ```

4. **Run tests with coverage**:
   ```bash
   go test -v -cover
   ```

## Test Coverage

The test suite covers:

- **Game Logic Functions**:
  - `isGameOver()` - Tests all winning conditions and draw scenarios
  - `MiniMax()` - Tests the minimax algorithm with various game states
  - `MakeBestMove()` - Tests optimal move selection

- **Bot Methods**:
  - `StatusPing()` - Tests ping response
  - `Register()` - Tests bot registration (mocked)
  - `Error()` - Tests error handling
  - `NextMove()` - Tests move calculation
  - `Complete()` - Tests game completion handling

- **HTTP Handler**:
  - `ServeHTTP()` - Tests all RPC method handlers

- **Utility Functions**:
  - `CreateRPCRequest()` - Tests request creation
  - `CreateRPCResponse()` - Tests response creation
  - `PrintGameState()` - Tests game board printing

- **Setters/Getters**:
  - `SetBaseUrl()` / `BaseUrl()`
  - `SetToken()` / `Token()`

## Running the Bot

To run the bot server:

```bash
export MERKNERA_URL="your_merknera_url"
export TOKEN="your_token"
export BOTNAME="your_bot_name"
export MY_URL="your_bot_url"

go run .
```

The bot will start an HTTP server on port 3003 and register itself with the game server.

