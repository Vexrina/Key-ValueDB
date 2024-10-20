package main

import (
	"BD/pkg/database"
	"BD/pkg/parser"
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	db := database.NewDataBaseImpl()
	parser := parser.ParserImpl{
		Databases: *db,
	}
	for {
		fmt.Print("our db $ ")
		cmd, _ := reader.ReadString('\n')

		result, err := parser.Parse(cmd)

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result)
	}
}
