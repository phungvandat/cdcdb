package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	db, closeDB := initClient()
	defer closeDB()

	e, closeE := NewEngine(db)
	defer closeE()
	go e.CDC(ctx)

	errChn := make(chan error)
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		errChn <- nil
	}()
	<-errChn
	cancel()
}

func PrintJSON(val interface{}) {
	b, err := json.MarshalIndent(val, "", "\t")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))
}
