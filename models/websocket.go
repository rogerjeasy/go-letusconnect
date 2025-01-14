package models

import "github.com/gofiber/websocket/v2"

type WebSocketMessage struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
	From    string      `json:"from"`
	To      string      `json:"to"`
	Time    int64       `json:"time"`
}

type WebSocketConnection struct {
	UserID string
	Conn   *websocket.Conn
}
