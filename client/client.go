package client

// 将package client注释后用于测试
// package main

import (
	//"fmt"
	"net"
)

// func main() {
// 	// 启动TCP客户端
// 	tcpData := "Hello, TCP server!"
// 	tcpResult, tcpErr := sendTCPRequest("127.0.0.1:12345", tcpData)
// 	//error
// 	if tcpErr != nil {
// 		fmt.Println("TCP client error:", tcpErr)
// 	} else {
// 		fmt.Println("TCP client received:", tcpResult)//tcp server recieve:
// 	}

// 	// 启动UDP客户端
// 	udpData := "Hello, UDP server!"
// 	udpResult, udpErr := sendUDPRequest("127.0.0.1:12346", udpData)
// 	//error
// 	if udpErr != nil {
// 		fmt.Println("UDP client error:", udpErr)
// 	} else {
// 		fmt.Println("UDP client received:", udpResult)
// 	}
// }

// 发送TCP请求
func SendTCPRequest(serverAddr string, data string) (string, error) {
	conn, err := net.Dial("tcp", serverAddr)// 连接服务器
	// error
	if err != nil {
		return "", err
	}
	defer conn.Close()// 压入栈，推迟到后续函数返回后执行

	_, err = conn.Write([]byte(data))// 数据写入连接，发送给服务器
	// error
	if err != nil {
		return "", err
	}

	buffer := make([]byte, 1024)// 数组切片unit8,buffer相当于缓存区
	n, err := conn.Read(buffer)
	if err != nil {
		return "", err
	}

	return string(buffer[:n]), nil// 获取读取到数据的部分
}

// 发送UDP请求
func SendUDPRequest(serverAddr string, data string) (string, error) {
	addr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		return "", err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(data))
	if err != nil {
		return "", err
	}

	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return "", err
	}

	return string(buffer[:n]), nil
}
