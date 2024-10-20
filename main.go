package main

import (
	"BD/pkg/database"
	"fmt"
	"time"
)

func main() {
	table := database.NewTableImpl()
	value := database.Value{
		"Hello, World!",
		time.Now(),
	}
	table.Add(1, value)

	fmt.Print(table.Size())
}
