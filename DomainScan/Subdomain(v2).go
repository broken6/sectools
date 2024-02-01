package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/parnurzeal/gorequest"
	"runtime"
	"strings"
	"sync"
	"time"
)

//go:embed domain_1800w.txt
var d string
var domain string
var thread int
var wg sync.WaitGroup
var pool sync.Pool
var current_url string
var green = color.New(color.FgGreen)
var red = color.New(color.FgRed)

func scan(ch chan string) {
	for x := range ch {
		defer wg.Done()
		http := fmt.Sprintf("http://%s", x)
		https := fmt.Sprintf("https://%s", x)
		request := gorequest.New()
		current_url = http
		resp, _, err := request.Head(http).Timeout(3 * time.Second).End()
		if err != nil {
			goto HTTPS
		}
		if resp.StatusCode == 200 || resp.StatusCode == 403 {
			pool.Put(http)
			green.Println("[+]存活->", http)
		}
	HTTPS:
		resp, _, err = request.Head(https).Timeout(3 * time.Second).End()
		current_url = https
		if err != nil {
			continue
		}
		if resp.StatusCode == 200 || resp.StatusCode == 403 {
			pool.Put(http)
			green.Println("[+]存活->", https)
		}
	}
}

func init() {
	flag.StringVar(&domain, "domain", "", "域名")
	flag.IntVar(&thread, "thread", 10, "线程 默认10 最大1024")
}

func main() {
	flag.Parse()
	if domain == "" {
		flag.Usage()
		return
	}
	if thread >= 1024 {
		thread = 1024
	}
	red.Println("环境检测")
	red.Println("当前系统", runtime.GOOS)
	red.Println("cpu核心数", runtime.NumCPU())
	red.Println("设置cpu核心数", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())
	red.Println("准备点火")
	time.Sleep(time.Second * 5)
	ch := make(chan string, 1024)
	for i := 0; i <= thread; i++ {
		go scan(ch)
	}
	t := time.Now()
	var buffer bytes.Buffer
	buffer.WriteString(d)
	var process int
	go func() {
		for {
			switch process % 4 {
			case 0:
				fmt.Print("\033[36m[/]\033[m")
			case 1:
				fmt.Print("\033[36m[-]\033[m")
			case 2:
				fmt.Print("\033[36m[\\]\033[m")
			case 3:
				fmt.Print("\033[36m[|]\033[m")
			}
			process += 1
			fmt.Printf("%s\r", current_url)
			time.Sleep(75 * time.Millisecond)
		}
	}()
	for {
		line, err := buffer.ReadString('\n')
		if err != nil {
			close(ch)
			break
		}
		wg.Add(1)
		line = strings.TrimSpace(line)
		sprintf := fmt.Sprintf("%s.%s", line, domain)
		ch <- sprintf
	}
	wg.Wait()
	red.Println("运行结束 总耗时->", time.Now().Sub(t).String())

}
