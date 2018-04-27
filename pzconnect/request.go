package pzconnect

import (
	"encoding/binary"

	"github.com/goinggo/tracelog"
	"github.com/ugorji/go/codec"
)

// Front controller which unmarshalls and processes the request
func handleRequest(conn IConnection, message []byte) bool {
	defaultMsgpackHeaderLen := 2

	// parse the buf
	msglength := len(message)
	tracelog.Trace("pzconnect", "handleRequest", "message len: %d, message content: %s", msglength, string(message))
	if msglength <= defaultMsgpackHeaderLen {
		tracelog.Info("pzconnect", "handleRequest", "invalid msg length, so possible corruption in pzlib, ignore msg")
		return false
	}
	msgpackLen := binary.LittleEndian.Uint16(message[0 : defaultMsgpackHeaderLen+1])
	if msgpackLen == 0 {
		tracelog.Info("pzconnect", "handleRequest", "invalid msgpack length, so possible corruption in pzlib, ignore msg")
		return false
	}
	msgpack := message[defaultMsgpackHeaderLen : defaultMsgpackHeaderLen+int(msgpackLen)+1] // end limit is exlusive
	payload := message[int(msgpackLen)+defaultMsgpackHeaderLen:]                            // start limit is inclusive
	tracelog.Trace("pzconnect", "handleRequest", "payload len: %d, payload: %s", len(payload), payload)

	// decode from msgpack
	var req pzMessage
	var dec = codec.NewDecoderBytes(msgpack, msgpackHandle)
	var err = dec.Decode(&req)
	if err != nil {
		tracelog.Error(err, "pzconnect", "handleRequest")
		return false
	}
	req.payload = payload

	// validate api version
	if req.API == 0 || req.API > apiVersion {
		tracelog.Info("pzconnect", "handleRequest", "API version (%d) is not supported", req.API)
		return false
	}

	// handle the message
	tracelog.Info("handlers", "handleRequest", "request : %+v", req)
	return handleDispatchRequest(conn, &req)
}

func handleDispatchRequest(conn IConnection, req *pzMessage) bool {
	switch req.EventType {
	case eventConnect:
		handleEventConnect(conn, req)
		return true
	case eventSend:
		if receiveCallback == nil {
			tracelog.Info("pzconnect", "handleRequest", "no receive callback registered")
			return false
		}

		receiveCallback(GameMessage{
			SenderID:      req.SenderID,
			Payload:       req.payload,
			DestinationID: req.DestinationID,
			Metadata:      req.Metadata,
		})
		return true
	}

	return false
}

func handleEventConnect(conn IConnection, req *pzMessage) {
	tracelog.Info("pzconnect", "handleEventConnect", "Event handleEventConnect triggered")

	connAsserted := conn.(*connection)
	connAsserted.client.ClientID = req.SenderID
	ok := clientsHub.ForceSet(req.SenderID, conn)
	if !ok {
		tracelog.Info("pzconnect", "NewConnection", "Some error storing connection in hub, closing it")
		conn.SendClose(1000, "Force set failure")
		return
	}

	if connectCallback == nil {
		tracelog.Info("pzconnect", "NewConnection", "no connect callback registered")
	} else {
		err := connectCallback(ConnectMessage{
			SenderID: req.SenderID,
			Payload:  req.payload,
			Metadata: req.Metadata,
		})
		if err != nil {
			tracelog.Errorf(err, "pzconnect", "NewConnection", "connect callback rejected connection with error")
			conn.SendClose(1000, "connection rejected")
			return
		}
	}

	sendResponse(req.SenderID, req.EventType, req.GameServerID, req.RequestUUID)
}

func sendResponse(receiverID uint64, method eventType, gameServerID string, requestUUID uint64) bool {
	tracelog.Trace("pzconnect", "sendResponse", "Target Client Id: %#x", receiverID)
	conn := clientsHub.Get(receiverID)
	if conn == nil {
		tracelog.Trace("pzconnect", "sendResponse", "receiver: %#x is not available", receiverID)
		return false
	}

	// all good, send out the response to target client
	message := pzMessage{
		API:          apiVersion,
		EventType:    eventResponse,
		ResponseTo:   method,
		Status:       reqStatusSuccess,
		GameServerID: gameServerID,
	}
	if requestUUID != 0 {
		message.RequestUUID = requestUUID
	}
	tracelog.Trace("pzconnect", "sendResponse", "message: %+v", message)
	messagePack, err := msgpackEncode(message)
	if err != nil {
		tracelog.Errorf(err, "pzconnect", "sendResponse", "Encode err")
		return false
	}

	buf, err := prepareBuf(messagePack, nil)
	if err != nil {
		tracelog.Errorf(err, "pzconnect", "sendResponse", "binary.Write failed")
		return false
	}

	tracelog.Trace("pzconnect", "sendResponse", "sending buff: %s", string(buf))
	conn.SendMessage(buf)
	return true
}
