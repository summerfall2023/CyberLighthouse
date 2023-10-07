package src

import (
	"fmt"
)

// 把二进制报文转换成packet
type ParsePacket struct {
	BinaryPacket []byte
	ParsedPacket Packet
}

func (p *ParsePacket) ParsePacket() (ParsePacket, error) {
	fmt.Println("START================================================")
	var res Packet
	var err1 error
	// header
	res.header, err1 = p.ParseHeader(p.BinaryPacket[0:12])
	if err1 != nil {
		return ParsePacket{}, err1
	}
	start := 12
	var currentIndex int
	query, currentIndex, err2 := p.ParseQuery(start)
	//fmt.Println("TTTTTTTTTT:cur = ", currentIndex)
	if err2 != nil {
		return ParsePacket{}, err2
	}
	res.queries = append(res.queries, query)
	// record
	var err3 error
	//fmt.Println("PARSEPACKET: index1", currentIndex, " AC: ", int(res.header.AC))
	res.answers, currentIndex, err3 = p.ParseResourceRecords(int(res.header.AC), currentIndex)
	//fmt.Println("PASREPACKET: index2", currentIndex)
	if err3 != nil {
		fmt.Println("PARSEPACKET: err3 ,", err3)
		return ParsePacket{}, err3
	}
	// fmt.Println("PASEPACKET: parser res.answer:", res.answers, " res.answer[0].", res.answers[0].RData.A_IP[0], " ", res.answers[0].RData.A_IP[1], " ", res.answers[0].RData.A_IP[2], " ", res.answers[0].RData.A_IP[3])
	res.authority, currentIndex, _ = p.ParseResourceRecords(int(res.header.NSC), currentIndex)
	res.additional, _, _ = p.ParseResourceRecords(int(res.header.AR), currentIndex)
	fmt.Println("PASEPACKET: parser res.additional:", res.additional, " res.additional[0].", res.additional[0])
	p.ParsedPacket = res
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
		// fmt.Println("PARSEHEADER:header-----------------")
		fmt.Println("PARSEHEADER", header)
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
	if currentIndex > len(p.BinaryPacket) {
		fmt.Println("ParseQuery expired")
		err3 := fmt.Errorf("parseQuery expired")
		return packetQuery{}, 0, err3
	}
	query.QName, currentIndex, err1 = p.ParseDomainName(start)
	if err1 != nil {
		return packetQuery{}, -1, err1
	}
	qTypeUint16, err1 := Byte2ToUint16(p.BinaryPacket[currentIndex : currentIndex+2])
	if err1 != nil {
		return packetQuery{}, -1, err1
	}
	query.QType = QueryType(qTypeUint16)
	qClassUint16, err2 := Byte2ToUint16(p.BinaryPacket[currentIndex+2 : currentIndex+4])
	if err2 != nil {
		return packetQuery{}, -1, err2
	}
	query.QClass = QueryClassType(qClassUint16)
	currentIndex += 4
	fmt.Println("QUETY: cur ", currentIndex)
	return query, currentIndex, nil
}

func (p *ParsePacket) ParseDomainName(offset int) (string, int, error) {
	var qname string
	currentIndex := offset

	for {
		labelLength := int(p.BinaryPacket[currentIndex])
		if labelLength == 0 {
			// 遇到长度为0的标签，表示域名结束
			currentIndex++
			break
		}

		currentIndex++
		labelBytes := p.BinaryPacket[currentIndex : currentIndex+labelLength]
		qname += string(labelBytes) + "."
		currentIndex += labelLength
		//}

		if currentIndex >= len(p.BinaryPacket) {
			return "", currentIndex, fmt.Errorf(ERROR_QNAME_END_MISSING)
		}
	}

	// 移除最后一个点号，得到完整的域名
	if len(qname) > 0 {
		qname = qname[:len(qname)-1]
	}
	return qname, currentIndex, nil
}

func (p *ParsePacket) ParseResourceRecords(count int, currentIndex int) ([]packetResource, int, error) {
	fmt.Println("RECORDS START```````````````````````````")
	i := 0
	var res []packetResource
	for i < count {
		var err error
		var packetResource packetResource
		packetResource, currentIndex, err = p.ParseResourceRecord(currentIndex)
		if err != nil {
			fmt.Println("PARSE_RESOURCE_RECORDS: ", err, " currentIndex", currentIndex)
			return nil, currentIndex, err
		}
		res = append(res, packetResource)
		i += 1
	}
	return res, currentIndex, nil
}
func (p *ParsePacket) ParseResourceRecord(offset int) (packetResource, int, error) {
	fmt.Println("RECORD START·············")
	var rr packetResource
	currentIndex := offset

	// 解析域名字段
	domainName, newOffset, err := p.ParseDomainName(currentIndex)
	if err != nil {
		fmt.Println("PARSE_DOMAIN_NAME: error ", err)
		return rr, currentIndex, err
	}
	rr.Name = domainName
	currentIndex = newOffset
	fmt.Println("RECORD_PARSE_DOMAIN_NAME 1: ", rr.Name)
	// 解析资源记录的类型、类、TTL和数据长度字段
	if currentIndex+10 > len(p.BinaryPacket) {
		err0 := fmt.Errorf("this is a resource record parsing error: unexpected end of data")
		fmt.Println("PARSE_RECORD: error ", err0, "\n		lenth of binary Packet: ", len(p.BinaryPacket))
		return rr, currentIndex, err0
	}
	//fmt.Println("TTTTTTTTTTTTTTTT:cur = ", currentIndex)
	Type, err1 := Byte2ToUint16(p.BinaryPacket[currentIndex : currentIndex+2])
	if err1 != nil {
		fmt.Println("PARSE_RECORD: err1 ", err1)
		return packetResource{}, -1, err1
	}
	rr.Type = QueryType(Type)
	fmt.Println("RECORD rr.Type: ", rr.Type)

	Class, err2 := Byte2ToUint16(p.BinaryPacket[currentIndex+2 : currentIndex+4])
	if err2 != nil {
		fmt.Println("PARSE_RECORD: err2 ", err2)
		return packetResource{}, -2, err2
	}
	rr.Class = QueryClassType(Class)
	TTL, err3 := Byte4ToUint32(p.BinaryPacket[currentIndex+4 : currentIndex+8])
	rr.TTL = TTL
	if err3 != nil {
		fmt.Println("PARSE_RECORD: err3 ", err3)
		return packetResource{}, -3, err3
	}

	DataLen, err4 := Byte2ToUint16(p.BinaryPacket[currentIndex+8 : currentIndex+10])
	if err4 != nil {
		fmt.Println("PARSE_RECORD: err4 ", err4)
		return packetResource{}, -4, err4
	}
	rr.ReLength = DataLen
	currentIndex += 10

	// 解析资源记录的数据字段
	if currentIndex+int(rr.ReLength) > len(p.BinaryPacket) {
		err6 := fmt.Errorf("resource record parsing error: unexpected end of data")
		//fmt.Println("PARSE_RECORD: ERROR ", err6, " CURRENTINDEX+INT(RR.RELENTH) = ", currentIndex+int(rr.ReLength), " CURRENTINDEX = ", currentIndex, " INT = ", int(rr.ReLength), " rr.RELENGT = ", rr.ReLength, " LEN(P.BINARYPACKET) = ", len(p.BinaryPacket))
		return rr, currentIndex, err6
	}
	fmt.Println("RDATA: cur ", currentIndex)
	RData, err5 := p.ParseRecourdData(rr, currentIndex)
	currentIndex += int(rr.ReLength)
	if err5 != nil {
		fmt.Println("PARSE_RECORD: err5 ", err5)
		return packetResource{}, -5, err5
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
			p.BinaryPacket[RDstartIndex],
			p.BinaryPacket[RDstartIndex+1],
			p.BinaryPacket[RDstartIndex+2],
			p.BinaryPacket[RDstartIndex+3],
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
			p.BinaryPacket[RDstartIndex],
			p.BinaryPacket[RDstartIndex+1],
			p.BinaryPacket[RDstartIndex+2],
			p.BinaryPacket[RDstartIndex+3],
			p.BinaryPacket[RDstartIndex+4],
			p.BinaryPacket[RDstartIndex+5],
			p.BinaryPacket[RDstartIndex+6],
			p.BinaryPacket[RDstartIndex+7],
			p.BinaryPacket[RDstartIndex+8],
			p.BinaryPacket[RDstartIndex+9],
			p.BinaryPacket[RDstartIndex+10],
			p.BinaryPacket[RDstartIndex+11],
			p.BinaryPacket[RDstartIndex+12],
			p.BinaryPacket[RDstartIndex+13],
			p.BinaryPacket[RDstartIndex+14],
			p.BinaryPacket[RDstartIndex+15],
		}
	case MX:

		preference := uint16(p.BinaryPacket[RDstartIndex])<<8 | uint16(p.BinaryPacket[RDstartIndex+1])
		RecordData.MX.MX_Preference = preference

		// MX记录中的RData包括一个域名
		name, _, err := p.ParseDomainName(RDstartIndex + 2) // 加2是因为优先级字段占用了2个字节
		if err != nil {
			return RecordData, err
		}
		RecordData.MX.MX_Name = name
	default:
	}
	return RecordData, nil
}
