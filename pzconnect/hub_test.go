package pzconnect

import (
	"testing"
)

type MockConnection struct{}

func (s *MockConnection) Start()                               {}
func (s *MockConnection) SendMessage(body []byte)              {}
func (s *MockConnection) SendClose(closeCode int, text string) {}
func (s *MockConnection) SetShouldRemoveFromHub(value bool)    {}

func setUp() {
	clientsHub.DeleteAll()
	clientsHub.Set(123, &MockConnection{})
	clientsHub.Set(456, &MockConnection{})
}

func tearDown() {
	clientsHub.Delete(123)
	clientsHub.Delete(456)
}

func TestGetMethod(t *testing.T) {
	var getTests = []struct {
		clientID uint64
		expected interface{} // expected result
	}{
		{123, &MockConnection{}},
		{456, &MockConnection{}},
		{12312, nil},
		{0, nil},
	}

	setUp()
	defer func() {
		tearDown()
	}()

	for _, tt := range getTests {
		actual := clientsHub.Get(tt.clientID)
		if actual != tt.expected {
			t.Errorf("clientsHub.Get(\"%d\"): expected %v, actual %v", tt.clientID, tt.expected, actual)
		}
	}
}

func TestSetMethod(t *testing.T) {
	var setTests = []struct {
		clientID uint64
		value    *MockConnection // value
		expected bool            //expected result
	}{
		{123, &MockConnection{}, false},
		{456, &MockConnection{}, false},
		{12312, &MockConnection{}, true},
		{0, &MockConnection{}, false},
	}

	setUp()
	defer func() {
		tearDown()
	}()

	for _, tt := range setTests {
		actual := clientsHub.Set(tt.clientID, tt.value)
		if actual != tt.expected {
			t.Errorf("clientsHub.Get(\"%d\"): expected %v, actual %v", tt.clientID, tt.expected, actual)
		}
	}
}

func TestDeleteMethod(t *testing.T) {
	var deleteTests = []struct {
		clientID uint64
		expected bool //expected result
	}{
		{123, true},
		{456, true},
		{12312, false},
		{0, false},
	}

	setUp()
	defer func() {
		tearDown()
	}()

	for _, tt := range deleteTests {
		actual := clientsHub.Delete(tt.clientID)
		if actual != tt.expected {
			t.Errorf("clientsHub.Get(\"%d\"): expected %v, actual %v", tt.clientID, tt.expected, actual)
		}
	}
}

func TestDeleteAllMethod(t *testing.T) {
	setUp()
	defer func() {
		tearDown()
	}()

	actual := clientsHub.DeleteAll()
	if actual != true {
		t.Errorf("clientsHub.DeleteAll(): failed (%v)", actual)
	}
}
