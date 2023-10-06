package src

import (
	"fmt"
)

// 把二进制报文转换成packet
type ParsePacket struct {
	binaryPacket []byte //???????????????????????????
	parsedPacket packet
}

func (p *ParsePacket) ParsePacket() (ParsePacket, error) {
	var res packet
	var err1 error
	// header
	res.header, err1 = p.ParseHeader(p.binaryPacket[0:12])
	if err1 != nil {
		return ParsePacket{}, err1
	}
	// queries
	// 通常QC是1，如果QC不是1如何处理？？？？？？？？？？？？？？？
	// for i := 1; i < int(res.header.QC); i++ {
	// 	ParseQuery()
	// }
	start := 12
	var currentIndex int
	query, currentIndex, err2 := p.ParseQuery(start)
	if err2 != nil {
		return ParsePacket{}, err2
	}
	res.queries = append(res.queries, query)
	// record
	res.answers, currentIndex = p.ParseResourceRecords(int(res.header.AC), currentIndex)
	res.authority, currentIndex = p.ParseResourceRecords(int(res.header.NSC), currentIndex)
	res.additional, currentIndex = p.ParseResourceRecords(int(res.header.AR), currentIndex)
	if currentIndex != len(p.binaryPacket)-1 {
		err3 := fmt.Errorf(ERROR_RECORD_LENGTH)
		return ParsePacket{}, err3
	}
	p.parsedPacket = res
	return *p, nil
}

// header 12 bytes-----------------------------------------
func (p *ParsePacket) ParseHeader(data []byte) (packetHeader, error) {
	if len(data) != 12 {
		return packetHeader{}, fmt.Errorf(ERROR_HEADER_LENGTH)
	}
	var err1, err2, err3, err4, err5, err6 error
	var header packetHeader
	header.ID, err1 = Byte2ToUint16(data[0:2])
	header.Flags, err2 = p.ParseHeaderFlags(data[2:4])
	header.QC, err3 = Byte2ToUint16(data[4:6])
	header.AC, err4 = Byte2ToUint16(data[6:8])
	header.NSC, err5 = Byte2ToUint16(data[8:10])
	header.AR, err6 = Byte2ToUint16(data[10:12])
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil {
		err := fmt.Errorf("error: At least one error occurred when parse header")
		return packetHeader{}, err
	} else {
		fmt.Println(header)
		return header, nil
	}
}

func (p *ParsePacket) ParseHeaderFlags(data []byte) (packetHeaderFlag, error) {
	if len(data) != 2 {
		err := fmt.Errorf(ERROR_HEADER_FLAG_LENGTH)
		return packetHeaderFlag{}, err
	}
	var flag packetHeaderFlag
	a := uint8(data[0])
	b := uint8(data[1])
	fmt.Println(a, data[0])
	flag.QR = bool((a&(1<<7))>>7 == 1)
	flag.OpCode = OpCodeType(a & ((0b1111) << 3) >> 3)
	flag.AA = bool(a&(1<<2) == 1)
	flag.TC = bool(a&(1<<1) == 1)
	flag.RD = bool(a&1 == 1)
	flag.RA = bool(b&(1<<7)>>7 == 1)
	flag.Z = bool(b&(1<<6)>>6 == 1)
	flag.AD = bool(b&(1<<5)>>5 == 1)
	flag.CD = bool(b&(1<<4)>>4 == 1)
	flag.RCode = RCodeType(b & (0b1111))
	return flag, nil
}

// queries----------------------------------------------
func (p *ParsePacket) ParseQuery(start int) (packetQuery, int, error) {
	var query packetQuery
	var err1 error
	var currentIndex int
	query.QName, currentIndex, err1 = p.ParseDomainName(start)
	if err1 != nil {
		return packetQuery{}, -1, err1
	}
	qTypeUint16, err1 := Byte2ToUint16(p.binaryPacket[currentIndex : currentIndex+2])
	if err1 != nil {
		return packetQuery{}, -1, err1
	}
	query.QType = QueryType(qTypeUint16)
	qClassUint16, err2 := Byte2ToUint16(p.binaryPacket[currentIndex+2 : currentIndex+4])
	if err2 != nil {
		return packetQuery{}, -1, err2
	}
	query.QClass = QueryClassType(qClassUint16)
	return query, currentIndex, nil
}

func (p *ParsePacket) ParseDomainName(offset int) (string, int, error) {
	var qname string
	currentIndex := offset

	for {
		labelLength := int(p.binaryPacket[currentIndex])
		if labelLength == 0 {
			// 遇到长度为0的标签，表示域名结束
			currentIndex++
			break
		}

		//DNS报文中的域名可以使用指针来压缩，以减小报文的大小，用于减少 DNS 报文的冗余。
		// if (labelLength & 0xC0) == 0xC0 {
		//     // 如果标签长度的两个最高位是1，表示这是一个指针
		//     // 解析指针并跳转到指定位置
		//     pointerOffset := int(binary.BigEndian.Uint16([]byte{p.binaryPacket[currentIndex], p.binaryPacket[currentIndex+1] & 0x3F}))
		//     currentIndex = pointerOffset
		// } else {
		// 非压缩标签，读取标签内容
		currentIndex++
		labelBytes := p.binaryPacket[currentIndex : currentIndex+labelLength]
		qname += string(labelBytes) + "."
		currentIndex += labelLength
		//}

		if currentIndex >= len(p.binaryPacket) {
			return "", currentIndex, fmt.Errorf(ERROR_QNAME_END_MISSING)
		}
	}

	// 移除最后一个点号，得到完整的域名
	if len(qname) > 0 {
		qname = qname[:len(qname)-1]
	}
	return qname, currentIndex, nil
}

//record-------------------------------------------------
//todo:delete after finishing coding
/*A记录（IPv4地址记录）：RData 包含一个4字节的IPv4地址，通常表示为 []byte 或 net.IP 类型。

AAAA记录（IPv6地址记录）：RData 包含一个16字节的IPv6地址，通常表示为 []byte 或 net.IP 类型。

CNAME记录（规范名称记录）：RData 包含一个规范名称，通常是一个域名字符串。

MX记录（邮件交换记录）：RData 包含邮件服务器的域名和优先级，通常表示为一个结构体或自定义类型。

TXT记录（文本记录）：RData 包含文本数据，通常是一个字符串。

NS记录（域名服务器记录）：RData 包含域名服务器的域名，通常表示为一个域名字符串。

SRV记录（服务记录）：RData 包含服务的相关信息，通常表示为一个结构体或自定义类型。
*/
func (p *ParsePacket) ParseResourceRecords(count int, currentIndex int) ([]packetResource, int) {
	i := 0
	var res []packetResource
	for i < count {
		var err error
		var packetResource packetResource
		packetResource, currentIndex, err = p.ParseResourceRecord(currentIndex)
		if err != nil {
			return nil, -1
		}
		res = append(res, packetResource)
	}
	return res, currentIndex
}
func (p *ParsePacket) ParseResourceRecord(offset int) (packetResource, int, error) {
	var rr packetResource
	currentIndex := offset

	// 解析域名字段
	domainName, newOffset, err := p.ParseDomainName(currentIndex)
	if err != nil {
		return rr, currentIndex, err
	}
	rr.Name = domainName
	currentIndex = newOffset

	// 解析资源记录的类型、类、TTL和数据长度字段
	if currentIndex+10 > len(p.binaryPacket) {
		err1 := fmt.Errorf("resource record parsing error: unexpected end of data")
		return rr, currentIndex, err1
	}

	Type, err1 := Byte2ToUint16(p.binaryPacket[currentIndex : currentIndex+2])
	if err1 != nil {
		return packetResource{}, -1, err1
	}
	rr.Type = QueryType(Type)

	Class, err2 := Byte2ToUint16(p.binaryPacket[currentIndex+2 : currentIndex+4])
	if err2 != nil {
		return packetResource{}, -1, err2
	}
	rr.Class = QueryClassType(Class)
	TTL, err3 := Byte4ToUint32(p.binaryPacket[currentIndex+4 : currentIndex+8])
	rr.TTL = TTL
	if err3 != nil {
		return packetResource{}, -1, err3
	}

	DataLen, err4 := Byte2ToUint16(p.binaryPacket[currentIndex+8 : currentIndex+10])
	if err4 != nil {
		return packetResource{}, -1, err4
	}
	rr.ReLength = DataLen
	currentIndex += 10

	// 解析资源记录的数据字段
	if currentIndex+int(rr.ReLength) > len(p.binaryPacket) {
		return rr, currentIndex, fmt.Errorf("resource record parsing error: unexpected end of data")
	}
	RData, err5 := p.ParseRecourdData(rr, currentIndex)
	currentIndex += int(rr.ReLength)
	if err5 != nil {
		return packetResource{}, -1, err5
	}
	rr.RData = RData
	return rr, currentIndex, nil
}

func (p *ParsePacket) ParseRecourdData(r packetResource, RDstartIndex int) (packetRecordData, error) {
	//switch
	var RecordData packetRecordData
	switch r.Type {
	case A:
		RecordData.A_IP = [4]byte{
			p.binaryPacket[RDstartIndex],
			p.binaryPacket[RDstartIndex+1],
			p.binaryPacket[RDstartIndex+2],
			p.binaryPacket[RDstartIndex+3],
		}
	case NS:
		var err error
		RecordData.NS_Name, _, err = p.ParseDomainName(RDstartIndex)
		if err != nil {
			return packetRecordData{}, err
		}
	case CNAME:
		var err error
		RecordData.NS_Name, _, err = p.ParseDomainName(RDstartIndex)
		if err != nil {
			return packetRecordData{}, err
		}
	case AAAA:
		// 解析AAAA记录的RData（IPv6地址）
		RecordData.AAAA_IP = [16]byte{
			p.binaryPacket[RDstartIndex],
			p.binaryPacket[RDstartIndex+1],
			p.binaryPacket[RDstartIndex+2],
			p.binaryPacket[RDstartIndex+3],
			p.binaryPacket[RDstartIndex+4],
			p.binaryPacket[RDstartIndex+5],
			p.binaryPacket[RDstartIndex+6],
			p.binaryPacket[RDstartIndex+7],
			p.binaryPacket[RDstartIndex+8],
			p.binaryPacket[RDstartIndex+9],
			p.binaryPacket[RDstartIndex+10],
			p.binaryPacket[RDstartIndex+11],
			p.binaryPacket[RDstartIndex+12],
			p.binaryPacket[RDstartIndex+13],
			p.binaryPacket[RDstartIndex+14],
			p.binaryPacket[RDstartIndex+15],
		}
	case MX:

		preference := uint16(p.binaryPacket[RDstartIndex])<<8 | uint16(p.binaryPacket[RDstartIndex+1])
		RecordData.MX.MX_Preference = preference

		// MX记录中的RData包括一个域名
		name, _, err := p.ParseDomainName(RDstartIndex + 2) // 加2是因为优先级字段占用了2个字节
		if err != nil {
			return RecordData, err
		}
		RecordData.MX.MX_Name = name
	default:
		return RecordData, fmt.Errorf("parseRecordData unsupported record type: %d", r.Type)
	}
	return RecordData, nil
}
