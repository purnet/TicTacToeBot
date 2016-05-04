package models

import "encoding/json"

type ClientRpcRequest struct {
	Method string      `json:"method"`
	Params interface{} `json:"params,omitempty"`
	Id     int         `json:"id,omitempty"`
}

type ClientRpcResponse struct {
	Result interface{} `json:"result"`
	Error  string      `json:"error,omitempty"`
	Id     int         `json:"id,omitempty"`
}

type ServerRpcRequest struct {
	Method string           `json:"method"`
	Params *json.RawMessage `json:"params,omitempty"`
	Id     int              `json:"id,omitempty"`
}

type ServerRpcResponse struct {
	Result *json.RawMessage `json:"result"`
	Error  string           `json:"error,omitempty"`
	Id     int              `json:"id,omitempty"`
}

type RegistrationParams struct {
	Token               string `json:"token"`
	BotName             string `json:"botname"`
	BotVersion          string `json:"botversion"`
	Game                string `json:"game"`
	RpcEndPoint         string `json:"rpcendpoint"`
	ProgrammingLanguage string `json:"programminglanguage"`
	Website             string `json:"website,omitempty"`
	Description         string `json:"description,omitempty"`
}

type RegistrationResponse struct {
	Message string `json:"message"`
}

type StatusPingResponse struct {
	Ping string `json:"ping"`
}

type ErrorParams struct {
	GameId    int    `json:"gameid"`
	Message   string `json:"message"`
	ErrorCode int    `json:"errorcode"`
}

type StatusResponseParams struct {
	Status string `json:"status"`
}

// Custom Models for TicTacToe only from here on
type NextMoveParams struct {
	GameId    int      `json:"gameid"`
	Mark      string   `json:"mark"`
	GameState []string `json:"gamestate"`
}

type NextMoveResponseParams struct {
	Position int `json:"position"`
}

type Complete struct {
	GameId    int      `json:"gameid"`
	Mark      string   `json:"mark"`
	Winner    bool     `json:"winner"`
	GameState []string `json:"gamestate"`
}
