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

func init() {
	flag.StringVar(&domain, "domain", "", "域名")
}

var domain string
var cout = color.New(color.FgGreen)
var wg sync.WaitGroup

//go:embed subdomains-top100.txt
var dict string

func GetHeadHtppAvaliable(url string) (int, []error) {
	request := gorequest.New().Timeout(3 * time.Second)
	resp, _, err := request.Head(url).End()
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, err
}

func main() {
	flag.Parse()
	if domain == "" {
		flag.Usage()
		return
	}
	var buffer bytes.Buffer
	buffer.WriteString(dict)
	t := time.Now()

	process := 0
	var currentURL string

	for {
		line, err := buffer.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		var reqHTTP = fmt.Sprintf("http://%s.%s", line, domain)
		var reqHTTPS = fmt.Sprintf("https://%s.%s", line, domain)

		wg.Add(2)
		go func(d string) {
			defer wg.Done()
			if status, err := GetHeadHtppAvaliable(d); err == nil && (status == 200 || status == 403) {
				cout.Println(fmt.Sprintf("存活-> %s", d))
			}
		}(reqHTTP)

		go func(d string) {
			defer wg.Done()
			if status, err := GetHeadHtppAvaliable(d); err == nil && (status == 200 || status == 403) {
				cout.Println(fmt.Sprintf("存活-> %s", d))
			}
		}(reqHTTPS)

		process += 1
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
		fmt.Printf("%s\r", currentURL)
	}

	wg.Wait()
	color.Red("扫描结束,耗时->%s\n", time.Now().Sub(t).String())
	fmt.Println("")
}
