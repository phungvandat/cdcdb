package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

var (
	log             = logrus.New()
	startPosition   uint64
	rlStartPosition = sync.RWMutex{}
	repConn         *pgx.ReplicationConn
	dbConn          *pgx.Conn
	once            = sync.Once{}
	mapTable        = map[uint32]*tableInfo{}
)

func init() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.FullTimestamp = true
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	log.SetFormatter(customFormatter)
}

func main() {
	var ctx, cancel = context.WithCancel(context.Background())
	initConn(ctx)
	defer func() {
		cancel()
		closeRepConn()
	}()
	setTableInfo()
	createSlot()
	startReplication()
	go receiveMessages(ctx)
	go keepAliveConn(ctx)

	errChn := make(chan error)
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		errChn <- fmt.Errorf("%s", <-ch)
	}()
	<-errChn
}

func PrintJSON(val interface{}) {
	b, err := json.MarshalIndent(val, "", "\t")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))
}
