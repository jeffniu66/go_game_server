package main

import "fmt"

func mainb() {
	for {
		fmt.Println("please enter data eg. AI 1000 10 ws://localhost:80/ws")
		var (
			tag      string
			roomNum  int
			matchNum int
			addr     string
		)
		fmt.Scan(&tag, &roomNum, &matchNum, &addr)
		if tag == "" || roomNum <= 0 || matchNum <= 0 || addr == "" {
			fmt.Println("please enter data eg. 1000 10 ws://localhost:80/ws")
			continue
		}
		AiLoginAndMatch(roomNum, matchNum, tag, addr)
		break
	}
}
