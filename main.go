package main

import (
	db "BD/pkg/database"
	server "BD/pkg/http"
	"BD/pkg/parser"
	"BD/pkg/xlog"
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	xlog.SetupLog()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databases := make(map[string]db.DataBaseImpl)
	parse := parser.ParserImpl{
		Databases: databases,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go server.Run(
		&parse,
		port,
	)
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
	} else {
		<-ctx.Done()
		xlog.Info("gracefully shutting down...")
	}
}
