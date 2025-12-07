package pubsub

import (
	"sync"
	"testing"
	"time"
)

// BenchmarkPublisher_Subscribe benchmarks subscriber creation
func BenchmarkPublisher_Subscribe(b *testing.B) {
	pub := NewPublisher(100)
	defer pub.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ch := pub.Subscribe()
		_ = ch
	}
}

// BenchmarkPublisher_SubscribeTopic benchmarks filtered subscriber creation
func BenchmarkPublisher_SubscribeTopic(b *testing.B) {
	pub := NewPublisher(100)
	defer pub.Close()

	filter := func(msg *Message) bool {
		return msg.Event == "test"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ch := pub.SubscribeTopic(filter)
		_ = ch
	}
}

// BenchmarkPublisher_Publish benchmarks message publishing
func BenchmarkPublisher_Publish(b *testing.B) {
	pub := NewPublisher(100)
	defer pub.Close()

	// Create subscribers
	for i := 0; i < 10; i++ {
		pub.Subscribe()
	}

	msg := &Message{
		Event:     "test",
		Data:      "benchmark data",
		Source:    "benchmark",
		TimeStamp: time.Now().Format(time.RFC3339),
		Expire:    300,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pub.Publish(msg)
	}
}

// BenchmarkPublisher_PublishWithFilter benchmarks filtered publishing
func BenchmarkPublisher_PublishWithFilter(b *testing.B) {
	pub := NewPublisher(100)
	defer pub.Close()

	// Create subscribers with filters
	for i := 0; i < 5; i++ {
		pub.SubscribeTopic(func(msg *Message) bool {
			return msg.Event == "test"
		})
	}

	// Create subscribers without filters
	for i := 0; i < 5; i++ {
		pub.Subscribe()
	}

	msg := &Message{
		Event:     "test",
		Data:      "benchmark data",
		Source:    "benchmark",
		TimeStamp: time.Now().Format(time.RFC3339),
		Expire:    300,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pub.Publish(msg)
	}
}

// BenchmarkPublisher_ConcurrentPublish benchmarks concurrent publishing
func BenchmarkPublisher_ConcurrentPublish(b *testing.B) {
	pub := NewPublisher(100)
	defer pub.Close()

	// Create subscribers
	for i := 0; i < 10; i++ {
		pub.Subscribe()
	}

	msg := &Message{
		Event:     "test",
		Data:      "benchmark data",
		Source:    "benchmark",
		TimeStamp: time.Now().Format(time.RFC3339),
		Expire:    300,
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pub.Publish(msg)
		}
	})
}

// BenchmarkPublisher_Evict benchmarks subscriber removal
func BenchmarkPublisher_Evict(b *testing.B) {
	pub := NewPublisher(100)
	defer pub.Close()

	// Create subscribers
	subscribers := make([]chan *Message, 1000)
	for i := 0; i < 1000; i++ {
		subscribers[i] = pub.Subscribe()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i < len(subscribers) {
			pub.Evict(subscribers[i])
		}
	}
}

// BenchmarkPublisher_SendTopic benchmarks direct message sending
func BenchmarkPublisher_SendTopic(b *testing.B) {
	pub := NewPublisher(100)

	ch := make(chan *Message, 100)
	msg := &Message{
		Event:     "test",
		Data:      "benchmark data",
		Source:    "benchmark",
		TimeStamp: time.Now().Format(time.RFC3339),
		Expire:    300,
	}

	var wg sync.WaitGroup

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pub.SendTopic(ch, nil, msg, &wg)
		}()
	}
	wg.Wait()

	// Drain the channel
	go func() {
		for range ch {
		}
	}()
}

// BenchmarkPublisher_MessageCreation benchmarks message creation overhead
func BenchmarkPublisher_MessageCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg := &Message{
			Event:     "test",
			Data:      "benchmark data",
			Source:    "benchmark",
			TimeStamp: time.Now().Format(time.RFC3339),
			Expire:    300,
		}
		_ = msg
	}
}

// BenchmarkPublisher_MessageWithExpire benchmarks message with expiration
func BenchmarkPublisher_MessageWithExpire(b *testing.B) {
	pub := NewPublisher(100)
	defer pub.Close()

	// Create subscribers
	for i := 0; i < 10; i++ {
		pub.Subscribe()
	}

	msg := &Message{
		Event:     "test",
		Data:      "benchmark data",
		Source:    "benchmark",
		TimeStamp: time.Now().Format(time.RFC3339),
		Expire:    1, // Short expiration
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pub.Publish(msg)
	}
}
