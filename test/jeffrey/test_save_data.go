package main

import (
	"fmt"
	"go_game_server/server/include"
)

func main() {
	item := &include.Item{ItemId: 101, Num: 1, UserId: 625272}
	fmt.Println("item: ", item)
	// db.SaveData(item)
}
