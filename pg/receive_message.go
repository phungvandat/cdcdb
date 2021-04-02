package main

import (
	"bytes"
	"context"
	"encoding/binary"
)

func receiveMessages(ctx context.Context) {
	for {
		mess, err := repConn.WaitForReplicationMessage(ctx)
		if err != nil {
			if err == context.Canceled {
				log.Warnln("Job canceled")
				return
			}
			log.WithError(err).Errorln("WaitForReplicationMessage")
			continue
		}
		if mess.WalMessage != nil {
			msg := decodeMessage(mess.WalMessage.WalData)
			if msg != nil {
				PrintJSON(msg)
			}

			var walStart = mess.WalMessage.WalStart
			if walStart > startPosition {
				setStartPosition(walStart)
				sendStandbyStatus()
			}
		}
	}
}

type decoder struct {
	order binary.ByteOrder
	buf   *bytes.Buffer
}

func (d *decoder) int32() uint32 {
	return d.order.Uint32(d.buf.Next(4))
}

func (d *decoder) newTuple() bool {
	return d.buf.Next(1)[0] == 'N'
}

func (d *decoder) tupleData() [][]byte {
	var (
		size = int(d.int16())
		data = make([][]byte, size)
	)
	for i := 0; i < size; i++ {
		s := d.buf.Next(1)
		if len(s) == 0 {
			continue
		}
		switch s[0] {
		case 't':
			data[i] = d.buf.Next(int(d.int32()))
		}
	}
	return data
}

func (d *decoder) int16() int16 {
	var val int16
	r := bytes.NewReader(d.buf.Next(2))
	_ = binary.Read(r, d.order, &val)
	return val
}

func decodeMessage(data []byte) *walDecodeMessage {
	if len(data) <= 1 {
		return nil
	}
	var (
		firstByte = data[0]
		d         = decoder{order: binary.BigEndian, buf: bytes.NewBuffer(data[1:])}
	)
	switch firstByte {
	case 'I': // insert
		return d.Insert()
	case 'U': // update
		return nil
	case 'D': // delete
		return nil
	}
	return nil
}

func (d *decoder) Insert() *walDecodeMessage {
	var (
		id        = d.int32()
		_         = d.newTuple()
		data      = d.tupleData()
		table, ok = mapTable[id]
		wdm       = &walDecodeMessage{
			Event:   "INSERT",
			Columns: make(map[string]interface{}),
		}
	)

	if !ok {
		log.WithField("table_oid", id).Warningf("not exists")
		return nil
	}

	wdm.Table = table.Name
	if len(table.OrderColumns) != len(data) {
		log.WithField("table_oid", id).Warningf("columns table changed")
		return nil
	}

	for idx := range data {
		var d = data[idx]
		wdm.Columns[table.OrderColumns[idx].Name] = string(d)
	}

	return wdm
}

type walDecodeMessage struct {
	Event   string                 `json:"event"`
	Table   string                 `json:"table"`
	Columns map[string]interface{} `json:"columns"`
}

func setStartPosition(walStart uint64) {
	rlStartPosition.Lock()
	defer rlStartPosition.Unlock()
	startPosition = walStart
}

func readStartPosition() uint64 {
	rlStartPosition.RLock()
	defer rlStartPosition.RUnlock()
	return startPosition
}
