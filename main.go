package main

import (
	"BD/pkg/database"
	"BD/pkg/http"
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

	go http.Run(make(map[string]database.DataBaseImpl), "8080")
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
