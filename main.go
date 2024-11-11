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
	go http.Run(make(map[string]database.DataBaseImpl), "8080")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	db := database.NewDataBaseImpl()
	parse := parser.ParserImpl{
		Databases: *db,
	}

	for {
		fmt.Print("our db $ ")
		cmd, _ := reader.ReadString('\n')

		result, err := parse.Parse(cmd)

		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(result)
	}

}
