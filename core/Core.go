package core

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"golang.org/x/sync/semaphore"
	"golang.org/x/net/context"
)

type Config struct {
	MaxConcurrency int
	Timeout        time.Duration
	Protocols      []string
	UserAgent      string
}

func Start(config Config, targets []string) {
    sem := semaphore.NewWeighted(int64(config.MaxConcurrency))
    ctx := context.Background()
    var wg sync.WaitGroup
    results := make(chan string, len(targets)*len(config.Protocols))

    for _, target := range targets {
        parts := strings.Split(target, ":")
        if len(parts) != 2 {
            fmt.Printf("Invalid target format: %s. Expected IP:port\n", target)
            continue
        }
        ip := parts[0]
        port := parts[1]

        for _, protocol := range config.Protocols {
            wg.Add(1)
            go func(ip, port, protocol string) {
                defer wg.Done()
                if err := sem.Acquire(ctx, 1); err != nil {
                    fmt.Printf("Failed to acquire semaphore: %v\n", err)
                    return
                }
                defer sem.Release(1)
                active, respCode := sendRequest(protocol, ip, port, config.Timeout, config.UserAgent)
                if active {
                    results <- fmt.Sprintf("\033[32m[%s] %s://%s:%s\033[0m", respCode, protocol, ip, port)
                } 
            }(ip, port, protocol)
        }
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    for result := range results {
        fmt.Println(result)
    }
}
