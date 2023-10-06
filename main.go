// package main

// import(
// 	"client"
// 	"fmt"
// 	"os"
// )

// func main(){
// 	for {
// 		fmt.Print("Enter command (or 'exit' to quit): ")
// 		var command string
// 		_, err := fmt.Scan(&command)
// 		if err != nil {
// 			fmt.Println("Error reading input:", err)
// 			os.Exit(1)
// 		}

// 		if command == "exit" {
// 			fmt.Println("Exiting the program.")
// 			break
// 		}

//			dataType,serverAddr,dataContent,dataToSend, err := client.ParseCommandLine()
//			if err != nil {
//				fmt.Println("Error:", err)
//				continue // 继续等待下一次输入
//			}
//			client.SendTCPRequest(*serverAddr,*dataContent)
//			fmt.Println("Data to send:", *dataType," ",string(dataToSend))
//		}
//	}





// package main

// import (
// 	"client"
// 	"flag"
// 	"fmt"
// 	"server"
// 	//"os"
// 	//"strconv"
// )

// func main() {
// 	server.CheckFrequency()
// 	for {
// 		fmt.Print("输入命令（或'exit'退出）: ")
// 		var command string
// 		fmt.Scanln(&command)

// 		if command == "exit" {
// 			fmt.Println("退出程序。")
// 			break
// 		}

// 		args := parseCommand(command)
// 		if len(args) == 0 {
// 			fmt.Println("无效命令。")
// 			continue
// 		}

// 		dataType := flag.String("type", "string", "数据类型: int 或 string")
// 		serverAddr := flag.String("addr", "127.0.0.1:5099", "服务器地址")
// 		dataContent := flag.String("data", "Hello, the world!", "数据内容")

// 		flag.Parse()

// 		if len(*dataType) == 0 || len(*serverAddr) == 0 || len(*dataContent) == 0 {
// 			fmt.Println("缺少必要参数。")
// 			continue
// 		}

// 		client.SendTCPRequest(*serverAddr, *dataContent)
// 		fmt.Println("待发送数据:", *dataType, " ", *dataContent)
// 	}
// }

// func parseCommand(command string) []string {
// 	args := []string{}
// 	arg := ""
// 	insideQuotes := false

// 	for _, char := range command {
// 		if char == ' ' && !insideQuotes {
// 			if len(arg) > 0 {
// 				args = append(args, arg)
// 				arg = ""
// 			}
// 		} else if char == '"' {
// 			insideQuotes = !insideQuotes
// 		} else {
// 			arg += string(char)
// 		}
// 	}

// 	if len(arg) > 0 {
// 		args = append(args, arg)
// 	}

// 	return args
// }

package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	// DNS服务器地址和端口（本地DNS服务器地址和端口）
	dnsServer := "127.0.0.1:53"

	// 要查询的域名
	domain := "github.com"

	// 创建UDP连接
	conn, err := net.Dial("udp", dnsServer)
	if err != nil {
		fmt.Println("Error creating UDP connection:", err)
		return
	}
	defer conn.Close()

	// 构造DNS查询报文
	query := buildDNSQuery(domain)

	// 发送DNS查询报文到DNS服务器
	_, err = conn.Write(query)
	if err != nil {
		fmt.Println("Error sending DNS query:", err)
		return
	}

	// 接收DNS响应报文
	buffer := make([]byte, 1024)
	_, err = conn.Read(buffer)
	if err != nil {
		fmt.Println("Error receiving DNS response:", err)
		return
	}

	// 解析DNS响应报文
	answers := parseDNSResponse(buffer)
	if len(answers) == 0 {
		fmt.Println("No A records found.")
		return
	}

	// 打印A记录信息
	fmt.Println("A Records:")
	for _, answer := range answers {
		fmt.Printf("Name: %s, IPv4 Address: %s\n", answer.Name, answer.IP)
	}
}

// 构造DNS查询报文
func buildDNSQuery(domain string) []byte {
	query := make([]byte, 12+len(domain)+5)

	// 设置DNS报文头部
	query[2] = 1 // 标志位：查询
	query[5] = 1 // 查询个数

	// 设置查询域名
	labels := []string{}
	for _, label := range splitLabels(domain) {
		labels = append(labels, string([]byte{byte(len(label))})+label)
	}
	domainName := fmt.Sprintf("%s\x00", strings.Join(labels, ""))
	copy(query[12:], []byte(domainName))

	// 设置查询类型（A记录）
	query[12+len(domain)+1] = 0
	query[12+len(domain)+2] = 1

	// 设置查询类（IN - Internet）
	query[12+len(domain)+4] = 1

	return query
}

// 解析DNS响应报文
func parseDNSResponse(response []byte) []DNSAnswer {
	// 解析逻辑...
	// 这里你需要编写代码来解析DNS响应报文，并提取A记录的信息
	// 请注意，实际的解析逻辑可能更复杂，需要处理多个响应、不同类型的DNS记录等情况

	// 示例代码返回一个空的DNSAnswer切片
	return []DNSAnswer{}
}

// 切分域名标签
func splitLabels(domain string) []string {
	labels := strings.Split(domain, ".")
	return labels
}

// DNSAnswer 表示DNS回答中的A记录
type DNSAnswer struct {
	Name string
	IP   string
}
