package main

import (
	_ "go_game_server/server/handler"
	"testing"
)

func TestAi(t *testing.T) {
	AiLogin(100, 10)
}

func TestAiMatch(t *testing.T) {
	AiLoginAndMatch(100, 10, "a", "ws://localhost:80/ws")
}
