package pzconnect

import (
	"reflect"
	"sync"

	"github.com/goinggo/tracelog"
)

// IHub is an interface to hub struct
type IHub interface {
	Get(key uint64) interface{}
	Set(key uint64, value interface{}) bool
	Delete(key uint64) bool
}

type hub struct {
	sync.Mutex
	clients map[uint64]IConnection
}

func newHub() *hub {
	return &hub{
		clients: make(map[uint64]IConnection),
	}
}

func (h *hub) set(clientID uint64, value IConnection) bool {
	if clientID == 0 {
		return false
	}
	if _, ok := h.clients[clientID]; ok {
		return false
	}

	h.clients[clientID] = value
	return true
}

func (h *hub) forceSet(clientID uint64, value IConnection) bool {
	if clientID == 0 {
		return false
	}
	_, ok := h.clients[clientID]
	if !ok {
		h.clients[clientID] = value
	} else {
		tracelog.Info("pzconnect", "reader", "force setting new connection: %d", clientID)
		conn := h.clients[clientID]

		// TODO: not efficient solution, better assingn uniqueid to each connection and compare
		if reflect.DeepEqual(conn, value) == true {
			tracelog.Info("pzconnect", "reader", "same connection, so not overwriting. skipping!!! %d", clientID)
			return true
		}

		conn.SetShouldRemoveFromHub(false)
		conn.SendClose(1000, "connection overwrite")

		delete(h.clients, clientID)
		h.clients[clientID] = value
	}

	return true
}

func (h *hub) get(clientID uint64) IConnection {
	if clientID == 0 {
		return nil
	}
	if value, ok := h.clients[clientID]; ok {
		return value
	}

	return nil
}

func (h *hub) delete(clientID uint64) bool {
	if clientID == 0 {
		return false
	}
	if _, ok := h.clients[clientID]; ok {
		delete(h.clients, clientID)
		return true
	}

	return false
}

func (h *hub) deleteAll() bool {
	for k := range h.clients {
		delete(h.clients, k)
	}

	return true
}

func (h *hub) Set(clientID uint64, value IConnection) bool {
	h.Lock()
	defer h.Unlock()
	return h.set(clientID, value)
}

func (h *hub) ForceSet(clientID uint64, value IConnection) bool {
	h.Lock()
	defer h.Unlock()
	return h.forceSet(clientID, value)
}

func (h *hub) Get(clientID uint64) IConnection {
	h.Lock()
	defer h.Unlock()
	return h.get(clientID)
}

func (h *hub) Delete(clientID uint64) bool {
	h.Lock()
	defer h.Unlock()

	return h.delete(clientID)
}

func (h *hub) DeleteAll() bool {
	h.Lock()
	defer h.Unlock()

	return h.deleteAll()
}

func (h *hub) GetClients() map[uint64]IConnection {
	return h.clients
}
