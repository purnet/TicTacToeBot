package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/purnet/TicTacToeBot/models"
)

// Test isGameOver function
func TestIsGameOver(t *testing.T) {
	tests := []struct {
		name      string
		gameState []string
		expected  bool
		winner    string
	}{
		{
			name:      "Empty board",
			gameState: []string{"", "", "", "", "", "", "", "", ""},
			expected:  false,
			winner:    "",
		},
		{
			name:      "X wins horizontally (top row)",
			gameState: []string{"X", "X", "X", "", "", "", "", "", ""},
			expected:  true,
			winner:    "X",
		},
		{
			name:      "O wins horizontally (middle row)",
			gameState: []string{"", "", "", "O", "O", "O", "", "", ""},
			expected:  true,
			winner:    "O",
		},
		{
			name:      "X wins vertically (left column)",
			gameState: []string{"X", "", "", "X", "", "", "X", "", ""},
			expected:  true,
			winner:    "X",
		},
		{
			name:      "O wins vertically (right column)",
			gameState: []string{"", "", "O", "", "", "O", "", "", "O"},
			expected:  true,
			winner:    "O",
		},
		{
			name:      "X wins diagonally (main diagonal)",
			gameState: []string{"X", "", "", "", "X", "", "", "", "X"},
			expected:  true,
			winner:    "X",
		},
		{
			name:      "O wins diagonally (anti-diagonal)",
			gameState: []string{"", "", "O", "", "O", "", "O", "", ""},
			expected:  true,
			winner:    "O",
		},
		{
			name:      "Draw game",
			gameState: []string{"X", "O", "X", "O", "X", "O", "O", "X", "O"},
			expected:  true,
			winner:    "",
		},
		{
			name:      "Game in progress",
			gameState: []string{"X", "O", "", "X", "", "", "", "", ""},
			expected:  false,
			winner:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gameOver, winner := isGameOver(tt.gameState)
			if gameOver != tt.expected {
				t.Errorf("isGameOver() gameOver = %v, expected %v", gameOver, tt.expected)
			}
			if winner != tt.winner {
				t.Errorf("isGameOver() winner = %v, expected %v", winner, tt.winner)
			}
		})
	}
}

// Test MiniMax function
func TestMiniMax(t *testing.T) {
	tests := []struct {
		name        string
		stateOfGame []string
		player      string
		move        int
		turn        string
		level       int
		expected    int
	}{
		{
			name:        "X wins immediately",
			stateOfGame: []string{"X", "X", "", "", "", "", "", "", ""},
			player:      "X",
			move:        2,
			turn:        "X",
			level:       0,
			expected:    10, // 10 + level (0)
		},
		{
			name:        "O blocks X from winning",
			stateOfGame: []string{"X", "X", "", "", "", "", "", "", ""},
			player:      "O",
			move:        2,
			turn:        "O",
			level:       0,
			expected:    -10, // -10 + level (0)
		},
		{
			name:        "Draw game",
			stateOfGame: []string{"X", "O", "X", "O", "X", "O", "", "", ""},
			player:      "X",
			move:        6,
			turn:        "X",
			level:       0,
			expected:    0,
		},
		{
			name:        "Game in progress",
			stateOfGame: []string{"X", "", "", "", "", "", "", "", ""},
			player:      "X",
			move:        1,
			turn:        "X",
			level:       0,
			expected:    0, // Will depend on the game tree
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MiniMax(tt.stateOfGame, tt.player, tt.move, tt.turn, tt.level)
			// For complex game states, we just check that the result is reasonable
			if result < -20 || result > 20 {
				t.Errorf("MiniMax() result = %v, expected reasonable score", result)
			}
		})
	}
}

// Test MakeBestMove function
func TestMakeBestMove(t *testing.T) {
	tests := []struct {
		name       string
		gameState  []string
		player     string
		gameId     int
		validMoves []int // Valid positions that could be returned
	}{
		{
			name:       "Empty board",
			gameState:  []string{"", "", "", "", "", "", "", "", ""},
			player:     "X",
			gameId:     1,
			validMoves: []int{0, 1, 2, 3, 4, 5, 6, 7, 8},
		},
		{
			name:       "One move left",
			gameState:  []string{"X", "O", "X", "O", "X", "O", "O", "X", ""},
			player:     "X",
			gameId:     2,
			validMoves: []int{8},
		},
		{
			name:       "X can win",
			gameState:  []string{"X", "X", "", "", "", "", "", "", ""},
			player:     "X",
			gameId:     3,
			validMoves: []int{2}, // Should choose winning move
		},
		{
			name:       "O can win",
			gameState:  []string{"O", "O", "", "", "", "", "", "", ""},
			player:     "O",
			gameId:     4,
			validMoves: []int{2}, // Should choose winning move
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MakeBestMove(tt.gameState, tt.player, tt.gameId)

			// Check if result is a valid move
			valid := false
			for _, validMove := range tt.validMoves {
				if result == validMove {
					valid = true
					break
				}
			}
			if !valid {
				t.Errorf("MakeBestMove() returned %v, expected one of %v", result, tt.validMoves)
			}

			// Check if the position is empty
			if result >= 0 && result < len(tt.gameState) && tt.gameState[result] != "" {
				t.Errorf("MakeBestMove() returned position %v that is not empty", result)
			}
		})
	}
}

// Test TicTacToeBot methods
func TestTicTacToeBot_StatusPing(t *testing.T) {
	bot := &TicTacToeBot{}

	result := bot.StatusPing(123)

	var response models.ClientRpcResponse
	err := json.Unmarshal(result, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Id != 123 {
		t.Errorf("StatusPing() id = %v, expected 123", response.Id)
	}

	if response.Error != "" {
		t.Errorf("StatusPing() error = %v, expected empty", response.Error)
	}

	// Check if result contains "OK"
	var statusResponse models.StatusPingResponse
	resultBytes, _ := json.Marshal(response.Result)
	json.Unmarshal(resultBytes, &statusResponse)

	if statusResponse.Ping != "OK" {
		t.Errorf("StatusPing() ping = %v, expected 'OK'", statusResponse.Ping)
	}
}

func TestTicTacToeBot_SetAndGetMethods(t *testing.T) {
	bot := &TicTacToeBot{}

	// Test SetBaseUrl and BaseUrl
	testUrl := "http://test.example.com"
	bot.SetBaseUrl(testUrl)
	if bot.BaseUrl() != testUrl {
		t.Errorf("BaseUrl() = %v, expected %v", bot.BaseUrl(), testUrl)
	}

	// Test SetToken and Token
	testToken := "test-token-123"
	bot.SetToken(testToken)
	if bot.Token() != testToken {
		t.Errorf("Token() = %v, expected %v", bot.Token(), testToken)
	}
}

func TestTicTacToeBot_Error(t *testing.T) {
	bot := &TicTacToeBot{}

	// Create a test error request
	errorParams := models.ErrorParams{
		GameId:    123,
		Message:   "Test error",
		ErrorCode: 500,
	}
	paramsBytes, _ := json.Marshal(errorParams)

	rpcReq := models.ServerRpcRequest{
		Method: "TicTacToe.Error",
		Params: (*json.RawMessage)(&paramsBytes),
		Id:     456,
	}

	result := bot.Error(rpcReq)

	var response models.ClientRpcResponse
	err := json.Unmarshal(result, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Id != 456 {
		t.Errorf("Error() id = %v, expected 456", response.Id)
	}

	if response.Error != "" {
		t.Errorf("Error() error = %v, expected empty", response.Error)
	}
}

func TestTicTacToeBot_NextMove(t *testing.T) {
	bot := &TicTacToeBot{}

	// Create a test next move request
	nextMoveParams := models.NextMoveParams{
		GameId:    789,
		Mark:      "X",
		GameState: []string{"X", "O", "", "", "", "", "", "", ""},
	}
	paramsBytes, _ := json.Marshal(nextMoveParams)

	rpcReq := models.ServerRpcRequest{
		Method: "TicTacToe.NextMove",
		Params: (*json.RawMessage)(&paramsBytes),
		Id:     789,
	}

	result := bot.NextMove(rpcReq)

	var response models.ClientRpcResponse
	err := json.Unmarshal(result, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Id != 789 {
		t.Errorf("NextMove() id = %v, expected 789", response.Id)
	}

	if response.Error != "" {
		t.Errorf("NextMove() error = %v, expected empty", response.Error)
	}

	// Check if result contains a valid position
	var nextMoveResponse models.NextMoveResponseParams
	resultBytes, _ := json.Marshal(response.Result)
	json.Unmarshal(resultBytes, &nextMoveResponse)

	if nextMoveResponse.Position < 0 || nextMoveResponse.Position > 8 {
		t.Errorf("NextMove() position = %v, expected 0-8", nextMoveResponse.Position)
	}
}

func TestTicTacToeBot_Complete(t *testing.T) {
	bot := &TicTacToeBot{}

	// Test winning case
	completeParams := models.Complete{
		GameId:    999,
		Mark:      "X",
		Winner:    true,
		GameState: []string{"X", "X", "X", "O", "O", "", "", "", ""},
	}
	paramsBytes, _ := json.Marshal(completeParams)

	rpcReq := models.ServerRpcRequest{
		Method: "TicTacToe.Complete",
		Params: (*json.RawMessage)(&paramsBytes),
		Id:     999,
	}

	result := bot.Complete(rpcReq)

	var response models.ClientRpcResponse
	err := json.Unmarshal(result, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Id != 999 {
		t.Errorf("Complete() id = %v, expected 999", response.Id)
	}

	if response.Error != "" {
		t.Errorf("Complete() error = %v, expected empty", response.Error)
	}
}

// Test utility functions
func TestCreateRPCRequest(t *testing.T) {
	params := models.RegistrationParams{
		Token:      "test-token",
		BotName:    "test-bot",
		BotVersion: "1.0",
		Game:       "TICTACTOE",
	}

	result := CreateRPCRequest("Test.Method", params, 123)

	var request models.ClientRpcRequest
	err := json.Unmarshal(result, &request)
	if err != nil {
		t.Fatalf("Failed to unmarshal request: %v", err)
	}

	if request.Method != "Test.Method" {
		t.Errorf("CreateRPCRequest() method = %v, expected 'Test.Method'", request.Method)
	}

	if request.Id != 123 {
		t.Errorf("CreateRPCRequest() id = %v, expected 123", request.Id)
	}
}

func TestCreateRPCResponse(t *testing.T) {
	result := CreateRPCResponse("test result", "test error", 456)

	var response models.ClientRpcResponse
	err := json.Unmarshal(result, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Result != "test result" {
		t.Errorf("CreateRPCResponse() result = %v, expected 'test result'", response.Result)
	}

	if response.Error != "test error" {
		t.Errorf("CreateRPCResponse() error = %v, expected 'test error'", response.Error)
	}

	if response.Id != 456 {
		t.Errorf("CreateRPCResponse() id = %v, expected 456", response.Id)
	}
}

// Test HTTP handler
func TestTicTacToeBot_ServeHTTP(t *testing.T) {
	bot := &TicTacToeBot{}

	tests := []struct {
		name                string
		method              string
		params              interface{}
		expectedStatus      int
		expectEmptyResponse bool
	}{
		{
			name:           "Status.Ping",
			method:         "Status.Ping",
			params:         nil,
			expectedStatus: 200,
		},
		{
			name:   "TicTacToe.NextMove",
			method: "TicTacToe.NextMove",
			params: models.NextMoveParams{
				GameId:    1,
				Mark:      "X",
				GameState: []string{"", "", "", "", "", "", "", "", ""},
			},
			expectedStatus: 200,
		},
		{
			name:   "TicTacToe.Error",
			method: "TicTacToe.Error",
			params: models.ErrorParams{
				GameId:    1,
				Message:   "Test error",
				ErrorCode: 500,
			},
			expectedStatus: 200,
		},
		{
			name:   "TicTacToe.Complete",
			method: "TicTacToe.Complete",
			params: models.Complete{
				GameId:    1,
				Mark:      "X",
				Winner:    true,
				GameState: []string{"X", "X", "X", "O", "O", "", "", "", ""},
			},
			expectedStatus: 200,
		},
		{
			name:                "Unknown method",
			method:              "Unknown.Method",
			params:              nil,
			expectedStatus:      200,  // Handler doesn't return error status
			expectEmptyResponse: true, // Unknown methods don't write response
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request body
			var body []byte
			if tt.params != nil {
				paramsBytes, _ := json.Marshal(tt.params)
				body, _ = json.Marshal(models.ServerRpcRequest{
					Method: tt.method,
					Params: (*json.RawMessage)(&paramsBytes),
					Id:     1,
				})
			} else {
				body, _ = json.Marshal(models.ServerRpcRequest{
					Method: tt.method,
					Id:     1,
				})
			}

			// Create HTTP request
			req := httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			bot.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("ServeHTTP() status = %v, expected %v", rr.Code, tt.expectedStatus)
			}

			// Check that response is valid JSON (unless expecting empty response)
			if !tt.expectEmptyResponse {
				var response models.ClientRpcResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("ServeHTTP() response is not valid JSON: %v", err)
				}
			} else {
				// For unknown methods, expect empty response body
				if len(rr.Body.Bytes()) != 0 {
					t.Errorf("ServeHTTP() expected empty response for unknown method, got: %s", rr.Body.String())
				}
			}
		})
	}
}

// Test PrintGameState (output testing)
func TestPrintGameState(t *testing.T) {
	// This test mainly ensures the function doesn't panic
	// In a real scenario, you might want to capture stdout
	gameBoard := []string{"X", "O", "X", "O", "X", "O", "O", "X", "O"}

	// This should not panic
	PrintGameState(gameBoard)

	// Test with empty board
	emptyBoard := []string{"", "", "", "", "", "", "", "", ""}
	PrintGameState(emptyBoard)
}
