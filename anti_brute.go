package main

import (
	"fmt"

	"github.com/tima-fey/anti_brute/internal/localDB"
)

func main() {
	localDB := localDB.DbInit()
	fmt.Println(localDB)
	out := make(chan bool, 3)
	localDB.Address.Add("test", out)
	answer := <-out
	fmt.Println(answer)
	fmt.Println(localDB)
}
