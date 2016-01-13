// Copyright © 2015 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package pubsub

import (
	"fmt"
	"github.com/thethingsnetwork/ttn/core"
	"net/http"
	"time"
)

var ErrConnectionLost = fmt.Errorf("Connection has been lost")

type regAckNacker struct {
	response chan regRes // A channel dedicated to send back a response
}

// Ack implements the core.Acker interface
func (r regAckNacker) Ack(p core.Packet) error {
	select {
	case r.response <- regRes{statusCode: http.StatusOK}:
		return nil
	case <-time.After(time.Millisecond * 50):
		return ErrConnectionLost
	}
}

// Nack implements the core.Nacker interface
func (r regAckNacker) Nack(p core.Packet) error {
	select {
	case r.response <- regRes{
		statusCode: http.StatusConflict,
		content:    []byte("Unable to register the given device"),
	}:
		return nil
	case <-time.After(time.Millisecond * 50):
		return ErrConnectionLost
	}
}
