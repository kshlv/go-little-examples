package main

import (
	"fmt"
	"sync"
	"time"
)

// Pubsub is a naïve and super-simple implementation of a pub-sub message broker.
type Pubsub struct {
	// `subs` contains a mapping between topics and lists of channels
	// through which subscribers fetch messages.
	subs map[string][]chan string
	// `closed == true` means that publishing is prohibited.
	closed bool

	// We have to use a mutex here because maps in Go are not concurrent.
	// (Removing the mutex won't break the code though.)
	sync.RWMutex
}

// NewPubSub creates a new Pubsub object.
func NewPubSub() *Pubsub {
	return &Pubsub{
		subs:   make(map[string][]chan string),
		closed: false,
	}
}

// Subscribe returns a channel through which
// messages on a certain topic are being sent to a subscriber.
func (ps *Pubsub) Subscribe(topic string) <-chan string {
	ps.Lock()
	defer ps.Unlock()

	ch := make(chan string, 10)
	ps.subs[topic] = append(ps.subs[topic], ch)

	return ch
}

// Publish publishes a msg on a certain topic.
func (ps *Pubsub) Publish(topic string, msg string) error {
	ps.RLock()
	defer ps.RUnlock()

	// Sending to a closed channel produces a panic,
	// so we must check if it is permitted to publish a message.
	if ps.closed {
		return fmt.Errorf("the pubsub is closed, message '%s' is not sent", msg)
	}

	for _, ch := range ps.subs[topic] {
		go func(ch chan string) {
			fmt.Printf("publishing a message %s on a topic %s\n", msg, topic)
			ch <- msg
		}(ch)
	}

	return nil
}

// Close disables the ability to publish messages to ps.
func (ps *Pubsub) Close() {
	ps.Lock()
	defer ps.Unlock()

	for topic, subs := range ps.subs {
		for _, ch := range subs {
			fmt.Printf("closing a channel %v for topic: %s\n", ch, topic)
			ch <- "closing the channel, good-bye"
			close(ch)
		}
	}

	ps.closed = true
}

// Publisher has a range of various topics and
// sends messages on them.
type Publisher struct {
	topics []string
}

// Subscriber subscribes on a topic
// and gets messages that publishers send.
type Subscriber struct {
	topics []string
	chs    []<-chan string
	// It also has a name so that we had a human-readable way
	// to distinguish subscribers.
	name string
}

var pubTopics []string = []string{
	"hello-topic",
	"test-topic",
}

var subTopics []string = []string{
	"hello-topic",
	"test-topic",
	"goodbye-topic",
}

func main() {
	pubsub := NewPubSub()

	p := Publisher{
		topics: pubTopics,
	}

	s := Subscriber{
		topics: subTopics,
		chs:    nil,
		name:   "test-sub",
	}

	for _, t := range s.topics {
		ch := pubsub.Subscribe(t)
		s.chs = append(s.chs, ch)
	}

	// We can also create a (s *Subscriber) Start() method.
	// But we're doing it in the main() for the sake of explicitness.
	go func(s Subscriber) {
		fmt.Println("starting a sub:", s.name)

		for _, ch := range s.chs {
			go func(ch <-chan string) {
				for {
					msg, more := <-ch
					if !more {
						fmt.Printf("the channel %v is closed\n", ch)
						return
					}
					fmt.Printf("sub %s got a message '%s' from a channel %v\n", s.name, msg, ch)
				}
			}(ch)
		}
	}(s)

	for _, t := range p.topics {
		go func(t, msg string) {
			count := 0
			for {
				time.Sleep(100 * time.Millisecond)
				if err := pubsub.Publish(t, fmt.Sprintf("%s-%s-%v", msg, t, count)); err != nil {
					fmt.Println("error publishing the message:", err)
					return
				}
				count++
			}
		}(t, "hello")
	}

	time.Sleep(2 * time.Second)
	pubsub.Close()
	time.Sleep(1 * time.Second)
}
