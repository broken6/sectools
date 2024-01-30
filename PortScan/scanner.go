package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/lcvvvv/gonmap"
)

var cout = color.New(color.FgRed) //定义扫描结果颜色
var mencolor = color.New(color.FgGreen)
var wg sync.WaitGroup

func scan(host string, port int) {
	defer wg.Done()
	nmap := gonmap.New() //定义nmap
	status, response := nmap.Scan(host, port)
	if status.String() == "Closed" { //如果端口关闭return
		return
	} else if response != nil && response.FingerPrint != nil { //确保指针不为空，防止访问无效内存
		cout.Println("[+]端口开放", "  ", port, "  ", response.FingerPrint.Service)
	} else {
		return
	}
}

func getUserInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func main() {
	var menStr string
	var host string

	for {
		mencolor.Println("===================================")
		mencolor.Println("        欢迎使用 PortScanner        ")
		mencolor.Println("         Author    br0ken        ")
		mencolor.Println("===================================")
		mencolor.Println("          0.   退出")
		mencolor.Println("          1. 开始扫描")

		menStr = getUserInput("请选择操作（输入对应数字）: ")
		men, err := strconv.Atoi(menStr)
		if err != nil {
			fmt.Println("无效选择")
			continue
		}

		switch men {
		case 1:
			host = getUserInput("请输入要扫描的主机:")

			if host == "" {
				fmt.Println("主机名不能为空")
				continue
			}

			for i := 1; i <= 65535; i++ {
				wg.Add(1)
				go scan(host, i)
			}
			wg.Wait() // 等待所有goroutine完成后打印结果

		case 0:
			fmt.Println("退出程序")
			return
		default:
			fmt.Println("无效选择")
		}
	}
}
