package pzconnect

const (
	apiVersion = 1
)

// LogLevel type
type LogLevel uint8

// Log levels exposed to game server
const (
	LogLevelTrace LogLevel = 1
	LogLevelInfo  LogLevel = 2
	LogLevelError LogLevel = 3
)

// EventType defines the list of event types supported
type eventType uint8

// Events exposed to game server
const (
	eventConnect  eventType = 0
	eventSend     eventType = 1
	eventResponse eventType = 2
)

// EventType defines the list of event types supported
type reqStatus uint16

// Events exposed to game server
const (
	reqStatusSuccess reqStatus = 0
	reqStatusFailure reqStatus = 1
)

// pzMessage is the msgpack encoded incoming message
type pzMessage struct {
	_struct       bool      `codec:",uint"`
	API           uint32    `codec:"0"`
	EventType     eventType `codec:"1"`
	SenderID      uint64    `codec:"2"`
	DestinationID string    `codec:"3"`
	PayloadLength uint32    `codec:"4"`
	AppID         uint64    `codec:"5"`
	ResponseTo    eventType `codec:"6"`
	RequestUUID   uint64    `codec:"7"`
	GameServerID  string    `codec:"8"`
	Status        reqStatus `codec:"9"`
	Metadata      uint32    `codec:"10"`
	payload       []byte
}

// GameMessage represents the message passed to the game server
type GameMessage struct {
	SenderID      uint64
	DestinationID string
	Metadata      uint32
	Payload       []byte
}

// ConnectMessage represents the connect method, passed to the game server
type ConnectMessage struct {
	SenderID uint64
	Metadata uint32
	Payload  []byte
}

// ReceiveHandlerFunc is a function that has a GameMessage as input
type ReceiveHandlerFunc func(GameMessage)

// ConnectHandlerFunc is a function that has a ConnectMessage as input
type ConnectHandlerFunc func(ConnectMessage) error

// global callback variables
var receiveCallback ReceiveHandlerFunc
var connectCallback ConnectHandlerFunc

// RegisterReceieveMessageCallback sets the receive callback func
func RegisterReceieveMessageCallback(callback ReceiveHandlerFunc) {
	receiveCallback = callback
}

// RegisterConnectCallback sets the connect callback func
func RegisterConnectCallback(callback ConnectHandlerFunc) {
	connectCallback = callback
}
