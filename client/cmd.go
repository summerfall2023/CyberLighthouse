package client

import (
	// "client/client"
	"flag"
	"fmt"
	"strconv"
)
func ParseCommandLine() (*string,*string,*string,[]byte,error){
	// 定义命令行标志
	dataType := flag.String("type", "string", "Data type: int or string")//标志，默认值，说明
	serverAddr := flag.String("addr", "127.0.0.1:5099", "Server address")
	dataContent := flag.String("data", "Hello, the world!", "Data content")
	flag.Parse()

	// 解析数据类型
	var dataToSend []byte
	var err error
	var intData int
	switch *dataType {
	case "int":
		// 将数据内容解析为整数
		intData, err = strconv.Atoi(*dataContent)
		if err != nil {
			fmt.Println("Invalid integer data:", err)
			return nil,nil,nil,nil,err
		}
		dataToSend = []byte(strconv.Itoa(intData))
	case "string":
		dataToSend = []byte(*dataContent)
	default:
		err = fmt.Errorf("Invalid data type. Supported types: int or string")
	}
	// 检查是否有错误发生
	if err != nil {
	// 返回错误
	return nil, nil, nil, nil, err
	}
	// // 发送TCP包
	// _,err := SendTCPRequest(*serverAddr, string(dataToSend))
	// // _,err := client.SendTCPRequest(*serverAddr, string(dataToSend))
	// if err != nil {
	// 	fmt.Println("TCP client error:", err)
	// } else {
	// 	fmt.Println("TCP client sent:", *dataContent)
	// }
	return dataType,serverAddr,dataContent,dataToSend,nil
}