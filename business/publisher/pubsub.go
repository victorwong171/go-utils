// Package pubsub provides a high-performance publish-subscribe messaging system.
// It supports topic filtering, concurrent operations, and message expiration.
//
// Example usage:
//
//	// Create a new publisher with buffer size
//	pub := pubsub.NewPublisher(100)
//
//	// Subscribe to all messages
//	allMessages := pub.Subscribe()
//
//	// Subscribe with topic filter
//	filteredMessages := pub.SubscribeTopic(func(msg *Message) bool {
//		return msg.Event == "user_action"
//	})
//
//	// Publish a message
//	msg := &pubsub.Message{
//		Event:     "user_action",
//		Data:      "user clicked button",
//		Source:    "web",
//		TimeStamp: time.Now().Format(time.RFC3339),
//		Expire:    300,
//	}
//	pub.Publish(msg)
//
//	// Clean up
//	pub.Close()
package pubsub

import (
	"sync"
	"time"
)

type (
	// subscriber represents a message channel for a subscriber
	subscriber chan *Message

	// topicFunc is a filter function that determines if a message should be sent to a subscriber
	topicFunc func(v *Message) bool
)

// Publisher manages subscribers and message distribution.
// It is safe for concurrent use by multiple goroutines.
type Publisher struct {
	m           sync.RWMutex             // protects subscribers map
	buffer      int                      // channel buffer size for new subscribers
	subscribers map[subscriber]topicFunc // active subscribers with their filters
}

// NewPublisher creates a new Publisher with the specified buffer size for subscriber channels.
// The buffer size determines how many messages can be queued for each subscriber before blocking.
//
// Example:
//
//	pub := pubsub.NewPublisher(100) // 100 message buffer per subscriber
func NewPublisher(buffer int) *Publisher {
	return &Publisher{
		buffer:      buffer,
		subscribers: make(map[subscriber]topicFunc),
	}
}

// Subscribe creates a new subscriber that receives all messages.
// It returns a channel that will receive all published messages.
//
// Example:
//
//	ch := pub.Subscribe()
//	for msg := range ch {
//		fmt.Printf("Received: %+v\n", msg)
//	}
func (p *Publisher) Subscribe() chan *Message {
	return p.SubscribeTopic(nil)
}

// SubscribeTopic creates a new subscriber with a topic filter.
// The filter function determines which messages the subscriber will receive.
// If filter is nil, the subscriber receives all messages.
//
// Example:
//
//	// Subscribe to specific events only
//	ch := pub.SubscribeTopic(func(msg *Message) bool {
//		return msg.Event == "user_action"
//	})
func (p *Publisher) SubscribeTopic(topic topicFunc) chan *Message {
	ch := make(chan *Message, p.buffer)
	p.m.Lock()
	defer p.m.Unlock()
	p.subscribers[ch] = topic
	return ch
}

// Evict removes a specific subscriber and closes its channel.
// It is safe to call Evict multiple times on the same channel.
//
// Example:
//
//	ch := pub.Subscribe()
//	// ... use channel ...
//	pub.Evict(ch) // Remove and close the channel
func (p *Publisher) Evict(sub chan *Message) {
	p.m.Lock()
	defer p.m.Unlock()
	if _, exists := p.subscribers[sub]; exists {
		delete(p.subscribers, sub)
		// Use select to avoid closing an already closed channel
		select {
		case <-sub:
			// channel is already closed
		default:
			close(sub)
		}
	}
}

// Close removes all subscribers and closes their channels.
// After calling Close, the Publisher should not be used for publishing new messages.
//
// Example:
//
//	pub.Close() // Clean up all subscribers
func (p *Publisher) Close() {
	p.m.Lock()
	defer p.m.Unlock()
	for sub := range p.subscribers {
		delete(p.subscribers, sub)
		// Use select to avoid closing an already closed channel
		select {
		case <-sub:
			// channel is already closed
		default:
			close(sub)
		}
	}
}

// Publish sends a message to all subscribers that match their topic filters.
// It blocks until all subscribers have been notified or the message expires.
//
// Example:
//
//	msg := &pubsub.Message{
//		Event:     "user_action",
//		Data:      "user clicked button",
//		Source:    "web",
//		TimeStamp: time.Now().Format(time.RFC3339),
//		Expire:    300, // 5 minutes
//	}
//	pub.Publish(msg)
func (p *Publisher) Publish(v *Message) {
	p.m.Lock()
	defer p.m.Unlock()
	var wg sync.WaitGroup
	for sub, topic := range p.subscribers {
		wg.Add(1)
		go p.SendTopic(sub, topic, v, &wg)
	}
	wg.Wait()
}

// SendTopic sends a message to a specific subscriber if it matches the topic filter.
// It respects the message expiration time and will timeout if the subscriber
// channel is full and the message expires.
func (p *Publisher) SendTopic(sub subscriber, topic topicFunc, v *Message, wg *sync.WaitGroup) {
	defer wg.Done()
	if topic != nil && !topic(v) {
		return
	}
	select {
	case sub <- v:
	case <-time.After(time.Duration(v.Expire) * time.Second):
		return
	}
}
