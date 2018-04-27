package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/packetzoom/PZConnect/pzconnect"
)

// handleReceivedMessages : this is simple example of handling send event
func handleReceivedMessages(gm pzconnect.GameMessage) {
	// get receiver ids from payload
	var payload payload

	if err := json.Unmarshal(gm.Payload, &payload); err != nil {
		fmt.Println("Invalid request:", err)
	}
	fmt.Printf("payload : %#v\n", payload)

	switch strings.ToLower(payload.Method) {
	case "send":
		handleSend(gm, payload)
	case "register":
		handleRegister(gm, payload)
	}
}

// handleConnect : gets called on new connections
func handleConnect(cm pzconnect.ConnectMessage) error {
	fmt.Printf("client: %#x just connected, with req: %v!", cm.SenderID, cm.Payload)
	return nil
}

func handleSend(gm pzconnect.GameMessage, message payload) {
	// send message to intended receiptents
	a := message.Data.ReceiverIDs
	b := make([]uint64, len(a))
	for i, typedValue := range a {
		value, err := strconv.ParseUint(typedValue, 10, 64)
		if err != nil {
			clientID := usernameClientMapRegistry.Get(typedValue)
			fmt.Println("clientID: %#x", clientID)
			if clientID != 0 {
				value = clientID
			}
		}
		b[i] = value
	}

	failedIDs := pzconnect.SendMessage(b, gm.Payload)
	fmt.Printf("failed ids : %#x\n", failedIDs)
}

func handleRegister(gm pzconnect.GameMessage, request payload) {
	// register the senderId
	username := register(gm.SenderID, request)

	// prepare list to all clients available
	var senders []string
	for _, key := range pzconnect.GetOnlineClientIds() {
		senderID := usernameClientMapRegistry.GetKey(key)
		if senderID != "" {
			senders = append(senders, senderID)
		}
	}
	data := map[string]interface{}{
		"senders": senders,
	}

	// prepare the response and broadcast
	event := request.Data.Event
	registerResp := &registerMessage{
		Sender: username,
		Event:  event,
		Data:   data,
	}
	resp, err := json.Marshal(registerResp)
	if err != nil {
		fmt.Println("Error marshalling:", err)
		return
	}
	pzconnect.Broadcast(resp)

	// also send the new name to current client
	resp2, err := json.Marshal(registerResponse{
		ID: request.ID,
		Result: &map[string]interface{}{
			"client_id": username,
		},
	})
	if err != nil {
		fmt.Println("Error marshalling:", err)
		return
	}

	failedIDs := pzconnect.SendMessage([]uint64{gm.SenderID}, resp2)
	fmt.Printf("failed ids : %#x\n", failedIDs)

}

func register(currentClientID uint64, request payload) string {
	username := request.Data.Username
	i := 0
	for {
		if i > 0 {
			username = fmt.Sprintf("%s%d", username, i)
		}

		i++
		clientID := usernameClientMapRegistry.Get(username)
		if clientID == 0 { // ok username is not assigned to anyone
			// check whether this client has existing username, if so replace it?
			currentSenderID := usernameClientMapRegistry.GetKey(currentClientID)
			if currentSenderID != "" {
				usernameClientMapRegistry.Delete(currentSenderID)
			}

			// first time registration, just assign it and break
			ok := usernameClientMapRegistry.Set(username, currentClientID)
			if !ok {
				fmt.Println("username already exists:", username)
			} else {
				break
			}
		} else {
			// check whether client id is same or not
			if clientID == currentClientID {
				break
			} else {
				fmt.Println("username already exists to someone:", username)
			}
		}
	}

	return username
}
