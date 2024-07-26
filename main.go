package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/bug5y/iprobe/core"
)

func main() {
	var config core.Config
	defaultTimeout := 10 * time.Second
	defaultThreads := 20
	defaultProtocols := []string{"http", "https"}
	defaultUserAgent := "Mozilla/5.0 (Windows NT 10.0; Android 13; Mobile; rv:120.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 OPR/95.0.0.0"

	input := flag.String("i", "", "Input (file path or IP:port)")
	threads := flag.Int("t", defaultThreads, "Max threads")
	flag.DurationVar(&config.Timeout, "timeout", defaultTimeout, "Timeout duration")
	protocols := flag.String("p", strings.Join(defaultProtocols, ","), "Comma-separated list of protocols")
	userAgent := flag.String("header", defaultUserAgent, "Custom user agent to add to each request")
	help := flag.Bool("h", false, "Show help")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	var targets []string

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		targets = readInput(os.Stdin)
	} else if *input != "" {
		if _, err := os.Stat(*input); err == nil {
			file, err := os.Open(*input)
			if err != nil {
				fmt.Printf("Error opening file: %v\n", err)
				os.Exit(1)
			}
			defer file.Close()
			targets = readInput(file)
		} else {
			targets = append(targets, *input)
		}
	} else {
		fmt.Println("Error: No input provided. Use -i flag to specify a file or IP:port, or pipe data to the program.")
		flag.Usage()
		os.Exit(1)
	}

	if *userAgent != "" {
		config.UserAgent = *userAgent
	}
	config.MaxConcurrency = *threads
	config.Protocols = strings.Split(*protocols, ",")

	core.Start(config, targets)
}

func readInput(r io.Reader) []string {
	var targets []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		targets = append(targets, strings.TrimSpace(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		os.Exit(1)
	}
	return targets
}