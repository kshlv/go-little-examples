package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignalFromClosedPubsub(t *testing.T) {
	ps := NewPubSub()
	ch := ps.Subscribe("dummy-topic")
	ps.Close()
	msg := <-ch
	assert.Equal(t, "closing the channel, good-bye", msg)
}

func TestPublishOnClosedPubsub(t *testing.T) {
	ps := NewPubSub()
	ps.Close()
	err := ps.Publish("closed-topic", "msg")
	assert.Error(t, err)
}

func TestPublish(t *testing.T) {
	ps := NewPubSub()
	ch := ps.Subscribe("test-topic")
	err := ps.Publish("test-topic", "test-message")
	assert.NoError(t, err)
	assert.Equal(t, "test-message", <-ch)
}
