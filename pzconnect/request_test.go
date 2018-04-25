package pzconnect

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"testing"

	"github.com/goinggo/tracelog"
	"github.com/ugorji/go/codec"
)

func TestHandleRequest(t *testing.T) {
	// log strat and stop
	tracelog.Start(tracelog.LevelError)
	defer func() { tracelog.Stop() }()

	var getTests = []struct {
		callback   ReceiveHandlerFunc
		connection IConnection
		payload    []byte
		expected   bool // expected success/failure
	}{
		{
			nil,
			&MockConnection{},
			nil, //corrupt data
			false,
		},
		{
			nil,
			&MockConnection{},
			getMsgPackData(123, []byte("123")),
			false,
		},
		{
			func(gm GameMessage) {
				// log.Println("test callback is fired")
			},
			&MockConnection{},
			getMsgPackData(123, []byte("123")),
			true,
		},
	}

	setUp()
	defer func() {
		tearDown()
	}()

	for _, tt := range getTests {
		receiveCallback = tt.callback

		actual := handleRequest(tt.connection, tt.payload)
		if tt.expected != actual {
			t.Errorf("handleRequest(\"%v\", \"%v\"): expected %v, actual %v", tt.connection, tt.payload, tt.expected, actual)
		}
	}
}

func getMsgPackData(senderID uint64, payload []byte) []byte {
	req := pzMessage{
		API:           1,
		EventType:     1,
		SenderID:      senderID,
		DestinationID: "ws://localhost:8080?client_id=" + strconv.FormatUint(senderID, 10),
		AppID:         1,
		PayloadLength: uint32(len(payload)),
	}

	var resp []byte
	enc := codec.NewEncoderBytes(&resp, msgpackHandle)
	err := enc.Encode(req)
	if err != nil {
		fmt.Printf("Encode err: %#v\n", err)
		return resp
	}

	// build buffer
	buf := buildBuf(resp, payload)

	return buf
}

func buildBuf(messagePack []byte, payload []byte) []byte {
	// build the buf
	// write msgpack length
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, uint16(len(messagePack)))
	if err != nil {
		fmt.Println("binary.Write failed for msgpacklength:", err)
	}

	// write msgpack
	err = binary.Write(buf, binary.LittleEndian, messagePack)
	if err != nil {
		fmt.Println("binary.Write failed for messagePack:", err)
	}

	// write payload
	err = binary.Write(buf, binary.LittleEndian, payload)
	if err != nil {
		fmt.Println("binary.Write failed for payload:", err)
	}

	return buf.Bytes()
}
