package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

type usernameClientMap struct {
	sync.Mutex
	usernameClient map[string]uint64
}

var usernameClientMapStore = "send_client.json"

func newUsernameClientMap() *usernameClientMap {
	var ucm map[string]uint64

	// load data from file
	data, err := ioutil.ReadFile(usernameClientMapStore)
	if err != nil {
		//log.Println("Info: send_client.json file not found")
	}

	if err := json.Unmarshal(data, &ucm); len(data) > 0 && err != nil {
		fmt.Println("Invalid data:", err)
	}

	if ucm == nil {
		ucm = make(map[string]uint64)
	}

	return &usernameClientMap{
		usernameClient: ucm,
	}
}

func (m *usernameClientMap) set(key string, value uint64) bool {
	if key == "" {
		return false
	}

	if _, ok := m.usernameClient[key]; ok {
		return false
	}

	m.usernameClient[key] = value

	// save to file
	saveUsernameClient(m.usernameClient)

	return true
}

func (m *usernameClientMap) get(key string) uint64 {
	if key == "" {
		return 0
	}

	if value, ok := m.usernameClient[key]; ok {
		return value
	}

	return 0
}

func (m *usernameClientMap) getKey(value uint64) string {
	if value == 0 {
		return ""
	}

	for k, v := range m.usernameClient {
		if v == value {
			return k
		}
	}

	return ""
}

func (m *usernameClientMap) delete(key string) bool {
	if key == "" {
		return false
	}

	if _, ok := m.usernameClient[key]; ok {
		delete(m.usernameClient, key)
		saveUsernameClient(m.usernameClient)
		return true
	}

	return false
}

func (m *usernameClientMap) Set(key string, value uint64) bool {
	m.Lock()
	defer m.Unlock()
	return m.set(key, value)
}

func (m *usernameClientMap) Get(key string) uint64 {
	m.Lock()
	defer m.Unlock()
	return m.get(key)
}

func (m *usernameClientMap) GetKey(value uint64) string {
	m.Lock()
	defer m.Unlock()
	return m.getKey(value)
}

func (m *usernameClientMap) Delete(key string) bool {
	m.Lock()
	defer m.Unlock()

	return m.delete(key)
}

func saveUsernameClient(usernameClient map[string]uint64) bool {
	// save to file
	v, err := json.Marshal(usernameClient)
	if err != nil {
		fmt.Println("Error marshalling usernameClient map :", err)
	}
	ioutil.WriteFile(usernameClientMapStore, v, 0644)

	return true
}
