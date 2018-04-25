package main

import (
	"github.com/packetzoom/PZConnect/pzconnect"
)

var usernameClientMapRegistry *usernameClientMap = newUsernameClientMap()

func main() {
	// register receive message callback
	pzconnect.RegisterReceieveMessageCallback(handleReceivedMessages)

	//register connect callback
	pzconnect.RegisterConnectCallback(handleConnect)

	// start the server
	pzconnect.Start(8080, pzconnect.LogLevelTrace)
}
