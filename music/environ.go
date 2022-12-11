package music

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// NatsPortOffset is the offset from the http port defined by env ADDRESS or PORT
const NatsPortOffset = 10

// NatsWebsocketPortOffset is the offset from the nats port to the websocket port
const NatsWebsocketPortOffset = 10

func listenAddress() string {

	if address := os.Getenv("ADDRESS"); address != "" {
		return address
	}

	if port := os.Getenv("PORT"); port == "" {
		return "localhost:8080"
	} else {
		return "localhost:" + port
	}

}

func NatsAddress() (host string, port int, err error) {

	addressParts := strings.Split(listenAddress(), ":")

	var p int64

	if p, err = strconv.ParseInt(addressParts[1], 10, 16); err != nil {
		return
	}

	return addressParts[0], int(p) + NatsPortOffset, nil

}

func NatsWebsocketURL() string {
	// server won't start without correctly parsed address
	// so this error here can be ignored
	host, port, _ := NatsAddress()

	return fmt.Sprintf("ws://%s:%d", host, port+NatsWebsocketPortOffset)

}
