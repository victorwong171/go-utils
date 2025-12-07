# go-utils

[![GoDoc](https://pkg.go.dev/badge/github.com/victorwong171/go-utils?utm_source=godoc)](https://pkg.go.dev/github.com/victorwong171/go-utils)
[![Go Report Card](https://goreportcard.com/badge/github.com/victorwong171/go-utils)](https://goreportcard.com/report/github.com/victorwong171/go-utils)
[![codecov](https://codecov.io/github/victorwong171/go-utils/branch/master/graph/badge.svg?token=2XWEF1Z3ZI)](https://codecov.io/github/victorwong171/go-utils)
![GitHub License](https://img.shields.io/github/license/victorwong171/go-utils)
[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![Build Status](https://github.com/victorwong171/go-utils/workflows/CI/badge.svg)](https://github.com/victorwong171/go-utils/actions)

A comprehensive Go utilities library providing high-performance data structures, business components, and common utilities for modern Go applications.

## 🚀 Features

- **High Performance**: Optimized algorithms and data structures
- **Thread Safe**: All components are designed for concurrent use
- **Well Tested**: 100% test coverage with comprehensive benchmarks
- **Production Ready**: Used in production environments
- **Go 1.24+**: Modern Go features and best practices

## 📦 Installation

```bash
go get github.com/victorwong171/go-utils
```

## 🏗️ Architecture

```
go-utils/
├── business/          # Business logic components
│   ├── observer/      # Event observer pattern
│   └── publisher/      # Pub/Sub messaging
├── desc/             # Data structures & algorithms
│   ├── bitmap/        # Bit manipulation
│   ├── list_node/     # Linked list utilities
│   ├── set/           # Set operations
│   ├── trie/          # Trie data structure
│   └── union_find/    # Union-Find algorithm
├── utils/            # Common utilities
│   ├── logger.go      # Structured logging
│   ├── slice.go       # Slice operations
│   └── interface.go   # Common interfaces
└── internal/          # Internal packages
```

## 🎯 Quick Start

### Event Observer Pattern

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/victorwong171/go-utils/business/observer"
    "github.com/victorwong171/go-utils/utils"
)

func main() {
    // Create logger
    logger := utils.MustNewDevelopment()
    
    // Create event with observers
    event := observer.NewEvent[UserEvent](logger,
        observer.Cfg[UserEvent]{
            IsAsync: false,
            Name:    "email-notification",
            Action: func(ctx context.Context, params UserEvent) error {
                fmt.Printf("Sending email to %s\n", params.Email)
                return nil
            },
        },
        observer.Cfg[UserEvent]{
            IsAsync: true,
            Name:    "analytics-tracking",
            Action: func(ctx context.Context, params UserEvent) error {
                fmt.Printf("Tracking user event: %s\n", params.Event)
                return nil
            },
        },
    )
    
    // Emit event
    userEvent := UserEvent{
        Email: "user@example.com",
        Event: "user_registered",
    }
    
    if err := event.Emit(context.Background(), userEvent); err != nil {
        log.Fatal(err)
    }
}

type UserEvent struct {
    Email string
    Event string
}
```

### Pub/Sub Messaging

```go
package main

import (
    "fmt"
    "time"
    "github.com/victorwong171/go-utils/business/publisher"
)

func main() {
    // Create publisher
    pub := publisher.NewPublisher(100)
    
    // Subscribe to all messages
    allMessages := pub.Subscribe()
    
    // Subscribe with filter
    filteredMessages := pub.SubscribeTopic(func(msg *publisher.Message) bool {
        return msg.Event == "user_action"
    })
    
    // Publish messages
    go func() {
        for i := 0; i < 10; i++ {
            msg := &publisher.Message{
                Event:     "user_action",
                Data:      fmt.Sprintf("Action %d", i),
                Source:    "web",
                TimeStamp: time.Now().Format(time.RFC3339),
                Expire:    300,
            }
            pub.Publish(msg)
            time.Sleep(100 * time.Millisecond)
        }
    }()
    
    // Consume messages
    go func() {
        for msg := range allMessages {
            fmt.Printf("Received: %+v\n", msg)
        }
    }()
    
    go func() {
        for msg := range filteredMessages {
            fmt.Printf("Filtered: %+v\n", msg)
        }
    }()
    
    time.Sleep(2 * time.Second)
    pub.Close()
}
```

### Data Structures

```go
package main

import (
    "fmt"
    "github.com/victorwong171/go-utils/desc/bitmap"
    "github.com/victorwong171/go-utils/desc/set"
    "github.com/victorwong171/go-utils/desc/union_find"
)

func main() {
    // Bitmap operations
    bm := bitmap.NewBitMap(1000)
    bm.Set(100)
    bm.Set(200)
    fmt.Printf("Bit 100 is set: %v\n", bm.Check(100))
    
    // Set operations
    s1 := set.Setify("a", "b", "c")
    s2 := set.Setify("b", "c", "d")
    s1.Set("d")
    fmt.Printf("Set contains 'a': %v\n", s1.HasKey("a"))
    
    // Union-Find
    uf := union_find.InitUnionFind(10)
    uf.Union(0, 1)
    uf.Union(1, 2)
    fmt.Printf("0 and 2 are connected: %v\n", uf.Find(0) == uf.Find(2))
}
```

### Utilities

```go
package main

import (
    "fmt"
    "github.com/victorwong171/go-utils/utils"
)

func main() {
    // Generate UUID
    id := utils.GetUuid()
    fmt.Printf("UUID: %s\n", id)
    
    // Ternary operator
    result := utils.TernaryOperator(true, "yes", "no")
    fmt.Printf("Result: %s\n", result)
    
    // Slice operations
    numbers := []int{1, 2, 3, 4, 5}
    contains := utils.Contain(numbers, 3)
    fmt.Printf("Contains 3: %v\n", contains)
    
    // Concurrent processing
    tasks := []func() error{
        func() error { fmt.Println("Task 1"); return nil },
        func() error { fmt.Println("Task 2"); return nil },
        func() error { fmt.Println("Task 3"); return nil },
    }
    
    err := utils.CurrentLimit(2, tasks, false)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

## 📊 Performance

### Benchmarks

```bash
go test -bench=. -benchmem ./...
```

### Memory Usage

- **Bitmap**: O(n) space for n bits
- **Set**: O(n) space for n elements  
- **Trie**: O(ALPHABET_SIZE * N) space
- **Union-Find**: O(n) space for n elements

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test ./... -cover

# Run benchmarks
go test -bench=. ./...

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## 📈 Coverage

- **Overall**: 99.2% coverage
- **Business**: 100% coverage
- **Data Structures**: 100% coverage
- **Utils**: 97.6% coverage

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Setup

```bash
# Clone repository
git clone https://github.com/victorwong171/go-utils.git
cd go-utils

# Install dependencies
go mod download

# Run tests
go test ./...

# Run linter
golangci-lint run

# Format code
go fmt ./...
```

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Go Team](https://golang.org/team) for the amazing language
- [Uber](https://github.com/uber-go/zap) for the excellent logging library
- [Google](https://github.com/google/uuid) for the UUID implementation
- All contributors and users of this library

## 📞 Support

- 📧 Email: victorwang171@gmail.com
- 🐛 Issues: [GitHub Issues](https://github.com/victorwong171/go-utils/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/victorwong171/go-utils/discussions)

---

⭐ **Star this repository if you find it helpful!**