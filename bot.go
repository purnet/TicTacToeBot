package main

import (
	"fmt"
	"sync"

	"log"
	"net/http"

	"io/ioutil"

	"bytes"
	"encoding/json"

	"os"

	"github.com/purnet/TicTacToeBot/models"
)

type GameBot interface {
	Register(game string, botName string, rpcendpoint string, botversion string, website string, description string) bool
	StatusPing(id int) []byte
	Error(rpcReq models.ServerRpcRequest) []byte
	ServeHTTP(rw http.ResponseWriter, req *http.Request)
	BaseUrl() string
	Token() string
	SetBaseUrl(baseUrl string)
	SetToken(token string)
}

type TicTacToeBot struct {
	baseUrl string
	token   string
}

func (b *TicTacToeBot) StatusPing(id int) []byte {
	s := models.StatusPingResponse{"OK"}
	rpc := CreateRPCResponse(s, "", id)
	return rpc
}

func (b *TicTacToeBot) Register(game string, botName string, rpcendpoint string, botversion string, website string, description string) bool {
	params := models.RegistrationParams{b.Token(), botName, botversion, game, rpcendpoint, "Go", website, description}
	JsonRpcBody := CreateRPCRequest("RegistrationService.Register", params, 1)
	fmt.Println(string(JsonRpcBody))
	respBody, _, _ := b.RpcRequest(JsonRpcBody)

	var resp models.ServerRpcResponse
	err := json.Unmarshal(respBody, &resp)
	if err != nil {
		fmt.Println(err)
	}
	rr := models.RegistrationResponse{}
	byteResult, e := json.Marshal(resp.Result)
	if e != nil {
		fmt.Println(e)
	}

	json.Unmarshal(byteResult, &rr)
	fmt.Println("Response message: ", rr.Message)

	return true

}

func (b *TicTacToeBot) SetBaseUrl(baseUrl string) {
	b.baseUrl = baseUrl
}

func (b *TicTacToeBot) SetToken(token string) {
	b.token = token
}

func (b *TicTacToeBot) BaseUrl() string {
	return b.baseUrl
}

func (b *TicTacToeBot) Token() string {
	return b.token
}

func (b *TicTacToeBot) Error(rpcReq models.ServerRpcRequest) []byte {
	params := models.ErrorParams{}
	byteResult, e := json.Marshal(rpcReq.Params)
	if e != nil {
		fmt.Println(e)
	}
	json.Unmarshal(byteResult, &params)
	fmt.Printf("Game: %v encounted Error: %v: %s\n", params.GameId, params.ErrorCode, params.Message)

	s := models.StatusResponseParams{"OK"}
	rpc := CreateRPCResponse(s, "", rpcReq.Id)
	return rpc
}

func CreateRPCRequest(method string, params interface{}, id int) []byte {
	req := models.ClientRpcRequest{method, params, id}
	JsonReqBody, _ := json.Marshal(req)
	return JsonReqBody
}

func CreateRPCResponse(result interface{}, error string, id int) []byte {
	resp := models.ClientRpcResponse{result, error, id}
	JsonRespBody, _ := json.Marshal(resp)
	return JsonRespBody
}

func (b TicTacToeBot) RpcRequest(body []byte) ([]byte, string, int) {
	client := &http.Client{}
	fmt.Printf("sdf%s", b.BaseUrl())
	req, err := http.NewRequest("POST", b.BaseUrl(), bytes.NewBuffer(body))
	req.Header.Add("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)

	return respBody, resp.Status, resp.StatusCode
}

func isGameOver(gs []string) (bool, string) {
	switch {
	case gs[0] != "" && gs[0] == gs[3] && gs[3] == gs[6]:
		return true, gs[0]
	case gs[0] != "" && gs[0] == gs[4] && gs[4] == gs[8]:
		return true, gs[0]
	case gs[1] != "" && gs[1] == gs[4] && gs[4] == gs[7]:
		return true, gs[1]
	case gs[2] != "" && gs[2] == gs[5] && gs[5] == gs[8]:
		return true, gs[2]
	case gs[2] != "" && gs[2] == gs[4] && gs[4] == gs[6]:
		return true, gs[2]
	case gs[0] != "" && gs[0] == gs[1] && gs[1] == gs[2]:
		return true, gs[0]
	case gs[3] != "" && gs[3] == gs[4] && gs[4] == gs[5]:
		return true, gs[3]
	case gs[6] != "" && gs[6] == gs[7] && gs[7] == gs[8]:
		return true, gs[6]
	default:
		for _, s := range gs {
			if s == "" {
				return false, ""
			}
		}
		return true, ""
	}
}

func MiniMax(stateOfGame []string, player string, move int, turn string, level int) int {
	sog := make([]string, 9, 9)
	copy(sog, stateOfGame)
	sog[move] = turn
	gameOver, piece := isGameOver(sog)
	if gameOver {
		switch piece {
		case "":
			return 0
		case player:
			return 10 + level
		default:
			return -10 + level
		}
	} else {
		var newTurn string
		if turn == "X" {
			newTurn = "O"
		} else {
			newTurn = "X"
		}
		moves := make(map[int]int)
		for i, s := range sog {
			if s == "" {
				moves[i] = MiniMax(sog, player, i, newTurn, level-1)
			}
		}
		var bestScore int
		scoreSet := false
		for _, score := range moves {
			if !scoreSet || (score > bestScore && newTurn == player) || (newTurn != player && score < bestScore) {
				bestScore = score
				scoreSet = true
			}
		}
		return bestScore
	}
}

type ChannelResult struct {
	Position int
	Score    int
}

func MakeBestMove(gameState []string, player string, gameId int) (pos int) {
	availableMoves := make(map[int]int)
	var wg sync.WaitGroup
	ch := make(chan ChannelResult)
	for i, s := range gameState {
		if s == "" {
			wg.Add(1)

			go func(c chan ChannelResult, pos int, state []string) {
				defer wg.Done()
				score := MiniMax(state, player, pos, player, 0)
				c <- ChannelResult{pos, score}
			}(ch, i, gameState)
		}
	}

	go func() {
		wg.Wait()
		fmt.Println("Closing Channel")
		close(ch)
	}()

	for c := range ch {
		fmt.Printf("Channel read: %d - %d\n", c.Position, c.Score)
		availableMoves[c.Position] = c.Score
	}

	var bestScore int
	scoreSet := false
	for move, score := range availableMoves {
		fmt.Printf("Game:%v Position: %v has score of %v \n", gameId, move, score)
		if !scoreSet || score > bestScore {
			bestScore = score
			pos = move
			scoreSet = true
		}
	}
	return pos
}

func (b TicTacToeBot) NextMove(rpcReq models.ServerRpcRequest) []byte {
	params := models.NextMoveParams{}
	byteResult, e := json.Marshal(rpcReq.Params)
	if e != nil {
		fmt.Println(e)
	}

	json.Unmarshal(byteResult, &params)
	fmt.Printf("Game: %v You are playing %s \n", params.GameId, params.Mark)
	PrintGameState(params.GameState)
	myMove := MakeBestMove(params.GameState, params.Mark, params.GameId)
	fmt.Printf("Game: %v your chosen move is position %v \n", params.GameId, myMove)
	pos := models.NextMoveResponseParams{myMove}
	rpc := CreateRPCResponse(pos, "", rpcReq.Id)
	return rpc
}

func (b TicTacToeBot) Complete(rpcReq models.ServerRpcRequest) []byte {
	params := models.Complete{}
	byteResult, e := json.Marshal(rpcReq.Params)
	if e != nil {
		fmt.Println(e)
	}

	json.Unmarshal(byteResult, &params)
	var tellMe string
	if params.Winner {
		tellMe = "Congatulations you WON!!"
	} else {
		tellMe = "Better Luck next time fool.."
	}
	fmt.Printf("%s GameId: %v where you were playing %s \n", tellMe, params.GameId, params.Mark)
	PrintGameState(params.GameState)
	s := models.StatusResponseParams{"OK"}
	rpc := CreateRPCResponse(s, "", rpcReq.Id)
	return rpc
}

func PrintGameState(gameBoard []string) {
	for num, s := range gameBoard {
		var val string
		if s == "" {
			val = " "
		} else {
			val = s
		}
		if (num+1)%3 == 0 {
			fmt.Printf("%s \n", val)
		} else {
			fmt.Printf("%s | ", val)
		}
	}
}

func (b *TicTacToeBot) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var rpcRequest models.ServerRpcRequest
	err := decoder.Decode(&rpcRequest)
	if err != nil {
		panic(err)
	}
	var body []byte
	switch rpcRequest.Method {
	case "Status.Ping":
		body = b.StatusPing(rpcRequest.Id)
		rw.Write(body)
	case "TicTacToe.NextMove":
		body = b.NextMove(rpcRequest)
		rw.Write(body)
	case "TicTacToe.Error":
		body = b.Error(rpcRequest)
		rw.Write(body)
	case "TicTacToe.Complete":
		body = b.Complete(rpcRequest)
		rw.Write(body)
	default:
		fmt.Printf("Request method %s is of unknown type\n", rpcRequest)
	}

	return
}

func main() {

	var b GameBot
	b = &TicTacToeBot{}
	b.SetBaseUrl(os.Getenv("MERKNERA_URL"))
	b.SetToken(os.Getenv("TOKEN"))

	if b.Register("TICTACTOE", os.Getenv("BOTNAME"), os.Getenv("MY_URL"), "2.1", "", "") {
		fmt.Println("Registration Complete... Tic Tac Toe Has begun")
	}

	http.Handle("/", b)

	err := http.ListenAndServe(":3003", nil)
	if err != nil {
		log.Fatal(err)
	}
}
