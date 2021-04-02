package main

import (
	"fmt"
)

const (
	protoVersion     = 1
	publicationNames = "wal_listener"
)

func startReplication() {
	var plgs = []string{
		fmt.Sprintf("proto_version '%v'", protoVersion),
		fmt.Sprintf("publication_names '%v'", publicationNames),
	}

	// SELECT * FROM pg_logical_slot_get_binary_changes('slotName', NULL, NULL, 'proto_version', '1', 'publication_names', 'wal_listener');
	err := repConn.StartReplication(slotName, startPosition, -1, plgs...)
	if err != nil {
		panic(err)
	}

	for oid := range mapTable {
		fmt.Printf("Listening %v...\n", mapTable[oid].Name)
	}
}
