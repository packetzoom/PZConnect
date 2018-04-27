

## Prerequisites
- go version `1.7.5`

## Usage

**Register Receive Message Callbacks**

Register a callback function for receiving messages.
```
//type GameMessage struct {
//	SenderID      uint64
//	DestinationID string
//	Metadata      uint32
//	Payload       []byte
//}

func handleReceivedMessages(gm pzconnect.GameMessage) {
	// process the message
}

pzconnect.RegisterReceieveMessageCallback(handleReceivedMessages)
```


**Register connect callback**

You will get a client ID when connection is established and reaady to use.
```

//type ConnectMessage struct {
//	SenderID      uint64
//	Metadata      uint32
//	Payload       []byte
//}

func handleConnect(cm pzconnect.ConnectMessage) error {
	fmt.Printf("client: %v just connected, with req: %v!", cm.SenderID, cm.Payload)
	return nil
}

pzconnect.RegisterConnectCallback(handleConnect)
```



**Start Server**

Starts the websocket server on specified port (make sure port is accesbile to internet), websocket server handles connections from both mode (PZConnect or Bypass)
```
var port = 8080

// start the server with port number and log level 
pzconnect.Start(port, pzconnect.LogLevelTrace)
```


**Send Message**

You can send message to any client but need to know his client ID.
```
SenderID := 1
ReceiverIDs := []uint64{123, 234}
Payload := []byte{1,2,3,4}

failedIDs := pzconnect.SendMessage(SenderID, ReceiverIDs, Payload)
fmt.Printf("failed ids : %v\n", failedIDs)
```



