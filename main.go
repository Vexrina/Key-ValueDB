package main

import (
	db "BD/pkg/database"
	server "BD/pkg/http"
	"BD/pkg/parser"
	"bufio"
	"fmt"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	peers := getPeers()

	databases := make(map[string]db.DataBaseImpl)
	parse := parser.ParserImpl{
		Databases: databases,
	}

	go server.Run(
		&parse,
		port,
		peers,
	)
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	if os.Getenv("CLI") == "y" {
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
}

func getPeers() []string {
	var peers []string
	idx := 1
	for {
		peer := os.Getenv("PEER" + fmt.Sprint(idx))
		if peer == "" {
			break
		}
		peers = append(peers, peer)
		idx += 1
	}
	return peers
}
