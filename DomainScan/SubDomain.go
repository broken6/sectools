package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/parnurzeal/gorequest"
	"strings"
	"sync"
	"time"
)

//go:embed subdomains-top100.txt
var d string
var domain string
var wg sync.WaitGroup
var green = color.New(color.FgGreen)

func scanWorker(ch <-chan string) {
	defer wg.Done()
	for x := range ch {
		http := fmt.Sprintf("http://%s", x)
		https := fmt.Sprintf("https://%s", x)
		request := gorequest.New()
		resp, _, err := request.Head(http).End()
		if err != nil {
			goto HTTPS
		}
		if resp.StatusCode == 200 || resp.StatusCode == 403 {
			green.Println("[+]存活->", http, "  ", resp.StatusCode)
		}
	HTTPS:
		resp, _, err = request.Head(https).End()
		if err != nil {
			continue
		}
		if resp.StatusCode == 200 || resp.StatusCode == 403 {
			green.Println("[+]存活->", https, "  ", resp.StatusCode)
		}
	}
}

func init() {
	flag.StringVar(&domain, "domain", "", "域名")
}

func main() {
	flag.Parse()
	if domain == "" {
		flag.Usage()
		return
	}

	startTime := time.Now()

	ch := make(chan string, 1024)

	// Start goroutines
	concurrency := 20
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go scanWorker(ch)
	}

	var buffer bytes.Buffer
	buffer.WriteString(d)

	// Send tasks to goroutines
	for {
		line, err := buffer.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		sprintf := fmt.Sprintf("%s.%s", line, domain)
		ch <- sprintf
	}

	// Close channel after sending all tasks
	close(ch)

	// Wait for all goroutines to finish
	wg.Wait()

	elapsed := time.Since(startTime)
	fmt.Println("扫描耗时:", elapsed)
}
