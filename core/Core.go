package core

import (
	"fmt"
	"strings"
	"sync"
	"time"
    "strconv"
    "net"
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
                active, respCode, hostname := sendRequest(protocol, ip, port, config.Timeout, config.UserAgent)
                if active {
                    if hostname == "" {
                        // Perform reverse DNS lookup if hostname is empty
                        names, err := net.LookupAddr(ip)
                        if err == nil && len(names) > 0 {
                            hostname = names[0]
                        }
                    }
                    var colorCode string
                    statusCode, _ := strconv.Atoi(respCode)
                    switch {
                    case statusCode >= 200 && statusCode < 300:
                        colorCode = "\033[32m" // Green
                    case statusCode >= 300 && statusCode < 400:
                        colorCode = "\033[34m" // Blue
                    case statusCode >= 400:
                        colorCode = "\033[31m" // Red
                    default:
                        colorCode = "\033[0m" // Default (no color)
                    }
                    
                    hostnameInfo := ""
                    if hostname != "" {
                        hostnameInfo = fmt.Sprintf(" - Hostname: %s", hostname)
                    }
                    
                    results <- fmt.Sprintf("%s[%s] %s://%s:%s%s\033[0m", colorCode, respCode, protocol, ip, port, hostnameInfo)
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
