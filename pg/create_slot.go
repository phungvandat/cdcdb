package main

import (
	"fmt"
	"os"
)

var slotName = os.Getenv("PG_SLOT_NAME")

const (
	dfSlotName     = "example_slot"
	dfOutputPlugin = "pgoutput"
)

func createSlot() {
	if slotName == "" {
		slotName = dfSlotName
	}

	// SELECT pg_create_logical_replication_slot('slotName', 'dfOutputPlugin');
	err := repConn.CreateReplicationSlot(slotName, dfOutputPlugin)
	if err != nil {
		if err.Error() == fmt.Sprintf("ERROR: replication slot \"%v\" already exists (SQLSTATE 42710)", slotName) {
			return
		}
		panic(err)
	}
}
