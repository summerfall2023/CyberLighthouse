package server
// package main

import (
	"fmt"
	"net"
)

// func main() {
// 	// 启动TCP和UDP服务器
// 	go startTCPServer()
// 	go startUDPServer()

// 	// 防止程序退出
// 	select {}
// }

// 启动TCP服务器
func StartTCPServer() {
	listen, err := net.Listen("tcp", "127.0.0.1:12345")
	if err != nil {
		fmt.Println("TCP server failed to start:", err)
		return
	}
	defer listen.Close()
	fmt.Println("TCP server is listening on :12345")

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("TCP server failed to accept connection:", err)
			continue
		}

		go handleTCPConnection(conn)
	}
}

// 处理TCP连接
func handleTCPConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("TCP server failed to read from connection:", err)
		return
	}

	data := string(buffer[:n])
	fmt.Println("TCP server received:", data)

	// 原路返回数据
	_, err = conn.Write([]byte("TCP server received: " + data))
	if err != nil {
		fmt.Println("TCP server failed to write to connection:", err)
	}
}

// 启动UDP服务器
func StartUDPServer() {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:12346")
	if err != nil {
		fmt.Println("UDP server failed to resolve address:", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("UDP server failed to start:", err)
		return
	}
	defer conn.Close()
	fmt.Println("UDP server is listening on :12346")

	buffer := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("UDP server failed to read from connection:", err)
			continue
		}

		data := string(buffer[:n])
		fmt.Printf("UDP server received from %s: %s\n", addr, data)

		// 原路返回数据
		_, err = conn.WriteToUDP([]byte("UDP server received: "+data), addr)
		if err != nil {
			fmt.Println("UDP server failed to write to connection:", err)
		}
	}
}
