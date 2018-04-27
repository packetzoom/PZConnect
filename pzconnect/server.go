package pzconnect

import (
	"fmt"
	"log"
	"net/http"

	"github.com/goinggo/tracelog"
	"github.com/gorilla/websocket"
	"github.com/ugorji/go/codec"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 1024,
	WriteBufferSize: 1024 * 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clientsHub = newHub()
var msgpackHandle = new(codec.MsgpackHandle)

// Start function intizalie and start the server to accept the connection from proxy/clients
func Start(port int, logLevel LogLevel) {
	// setup log level
	switch logLevel {
	case LogLevelError:
		tracelog.Start(tracelog.LevelError)
	case LogLevelInfo:
		tracelog.Start(tracelog.LevelInfo)
	default:
		tracelog.Start(tracelog.LevelTrace)
	}

	// define routes
	http.HandleFunc("/echo", serveEcho)
	http.HandleFunc("/api", serveAPI)

	// turn off the legacy support
	msgpackHandle.WriteExt = true
	msgpackHandle.RawToString = true

	// start the server
	tracelog.Info("pzconnect", "start", "starting websocket server on port %d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	log.Fatal("ListenAndServe: ", err)

	tracelog.Stop()
}

func serveEcho(w http.ResponseWriter, r *http.Request) {
	tracelog.Info("pzconnect", "serveAPI", "got websocket request to %s", r.URL)
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		tracelog.Errorf(err, "pzconnect", "serveEcho", "Error in upgrading to websocket")
		return
	}
	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			tracelog.Errorf(err, "pzconnect", "serveEcho", "read in upstream")
			break
		}
		tracelog.Info("pzconnect", "serveAPI", "recv in upstream: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			tracelog.Errorf(err, "pzconnect", "serveEcho", "write in upstream:")
			break
		}
	}
}

// handle a request to /api
func serveAPI(w http.ResponseWriter, r *http.Request) {
	tracelog.Info("pzconnect", "serveAPI", "got websocket request to %s", r.URL)

	if r.Method != "GET" {
		tracelog.Info("pzconnect", "serveAPI", "invalid method %s", r.Method)
		httpError(w, http.StatusMethodNotAllowed)
		return
	}

	params := r.URL.Query()
	params.Set("timeout", "0")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		tracelog.Errorf(err, "pzconnect", "serveAPI", "Error in upgrading to websocket")
		return
	}

	c := newConnection(ws, r, params)
	if c != nil {
		c.Start()
	}
}

func httpError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
