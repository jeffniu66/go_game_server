package main

import (
	"go_game_server/robot/client_robot"
)

//GOOS=windows GOARCH=386 go build -ldflags '-w -s' -o robot.exe client.go

func main() {
	client_robot.InitHandler()
	client_robot.InitReceive()
	client_robot.Start()
}
