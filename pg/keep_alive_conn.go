package main

import (
	"context"
	"time"

	"github.com/jackc/pgx"
)

// keep alive connection
func keepAliveConn(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.NewTicker(time.Second).C:
			sendStandbyStatus()
		}
	}
}

func sendStandbyStatus() {
	num := readStartPosition()
	status, err := pgx.NewStandbyStatus(num)
	if err != nil {
		log.WithError(err).Errorln("NewStandbyStatus")
		return
	}
	status.ReplyRequested = 0
	err = repConn.SendStandbyStatus(status)
	if err != nil {
		log.WithError(err).Errorln("SendStandbyStatus")
	}
}
