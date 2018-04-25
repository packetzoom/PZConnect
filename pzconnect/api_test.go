package pzconnect

import (
	"testing"

	"github.com/goinggo/tracelog"
)

func TestSendMessage(t *testing.T) {
	// log strat and stop
	tracelog.Start(tracelog.LevelError)
	defer func() { tracelog.Stop() }()

	var getTests = []struct {
		receiverIDs []uint64
		payload     []byte
		expected    int // expected failed ids count
	}{
		{
			[]uint64{456},
			[]byte("other user"),
			0,
		},
		{
			[]uint64{123},
			[]byte("self"),
			0,
		},
		{
			[]uint64{0},
			[]byte("wrong user"),
			1,
		},
	}

	setUp()
	defer func() {
		tearDown()
	}()

	for _, tt := range getTests {
		actual := SendMessage(tt.receiverIDs, tt.payload)
		if tt.expected != len(actual) {
			t.Errorf("SendMessage(\"%v\", \"%v\"): expected %v, actual %v", tt.receiverIDs, tt.payload, tt.expected, actual)
		}
	}
}

func TestBroadcast(t *testing.T) {
	// log strat and stop
	tracelog.Start(tracelog.LevelError)
	defer func() { tracelog.Stop() }()

	var getTests = []struct {
		payload  []byte
		expected int // expected failed ids count
	}{
		{
			[]byte("other user"),
			0,
		},
	}

	setUp()
	defer func() {
		tearDown()
	}()

	for _, tt := range getTests {
		actual := Broadcast(tt.payload)
		if tt.expected != len(actual) {
			t.Errorf("SendMessage(\"%v\"): expected %v, actual %v", tt.payload, tt.expected, actual)
		}
	}
}
