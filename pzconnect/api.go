package pzconnect

import (
	"bytes"
	"encoding/binary"

	"github.com/goinggo/tracelog"
	"github.com/ugorji/go/codec"
)

// GetOnlineClientIds returns all connected clients
func GetOnlineClientIds() []uint64 {
	tracelog.Info("pzconnect", "GetOnlineClientIds", "GetOnlineClientIds")
	var onlineIDs []uint64
	for c := range clientsHub.GetClients() {
		onlineIDs = append(onlineIDs, c)
	}

	return onlineIDs
}

// Broadcast send message to all connected clients
func Broadcast(payload []byte) []uint64 {
	tracelog.Info("pzconnect", "broadcast", "broadcast message")
	failedIDs := SendMessage(GetOnlineClientIds(), payload)

	tracelog.Trace("pzconnect", "Broadcast", "failedIDs %#x", failedIDs)
	return failedIDs
}

// SendMessage is a send event handler, which send message to reciepents
func SendMessage(receiverIDs []uint64, payload []byte) []uint64 {
	tracelog.Info("pzconnect", "SendMessage", "Sending message")

	// create buf and reuse for all receivers
	buf := createBufFromPayload(payload)
	if buf == nil {
		tracelog.Info("pzconnect", "sendMessage", "Buf creation failed")
		return receiverIDs
	}

	var failedIDs []uint64
	for _, receiverID := range receiverIDs {
		success := sendMessage(receiverID, buf)
		if !success {
			failedIDs = append(failedIDs, receiverID)
		}
	}

	tracelog.Trace("pzconnect", "SendMessage", "failedIDs %#x", failedIDs)
	return failedIDs
}

func sendMessage(receiverID uint64, buf []byte) bool {
	tracelog.Trace("pzconnect", "sendMessage", "Target Client Id: %#x", receiverID)
	conn := clientsHub.Get(receiverID)
	if conn == nil {
		tracelog.Trace("pzconnect", "sendMessage", "receiver: %#x is not available", receiverID)
		return false
	}

	conn.SendMessage(buf)
	return true
}

func createBufFromPayload(payload []byte) []byte {
	payloadLength := uint32(len(payload))
	message := pzMessage{
		PayloadLength: payloadLength,
		API:           apiVersion,
	}
	messagePack, err := msgpackEncode(message)
	if err != nil {
		tracelog.Errorf(err, "pzconnect", "sendMessage", "Encode err")
		return nil
	}

	buf, err := prepareBuf(messagePack, payload)
	if err != nil {
		tracelog.Errorf(err, "pzconnect", "sendMessage", "binary.Write failed")
		return nil
	}

	return buf
}

func prepareBuf(messagePack []byte, payload []byte) ([]byte, error) {
	// write msgpack length
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, uint16(len(messagePack)))
	if err != nil {
		tracelog.Errorf(err, "pzconnect", "createBuf", "binary.Write failed on len")
		return nil, err
	}

	// write msgpack
	err = binary.Write(buf, binary.LittleEndian, messagePack)
	if err != nil {
		tracelog.Errorf(err, "pzconnect", "createBuf", "binary.Write failed on msgpack")
		return nil, err
	}

	// write payload
	if payload != nil {
		err = binary.Write(buf, binary.LittleEndian, payload)
		if err != nil {
			tracelog.Errorf(err, "pzconnect", "createBuf", "binary.Write failed on payload")
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func msgpackEncode(message pzMessage) ([]byte, error) {
	var messagePack []byte
	enc := codec.NewEncoderBytes(&messagePack, msgpackHandle)
	err := enc.Encode(message)
	if err != nil {
		tracelog.Errorf(err, "pzconnect", "sendMessage", "Encode err")
		return nil, err
	}

	return messagePack, nil
}
