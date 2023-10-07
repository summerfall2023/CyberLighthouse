package src

import "strings"

type GeneratePacket struct {
	OriginData Packet
	BinaryData []byte
}

func (p *GeneratePacket) GeneratePacket() GeneratePacket {
	var res []byte
	//var err1 error
	header := p.GenerateHeader()
	res = append(res, header...)
	query := p.GenerateQuery()
	res = append(res, query...)
	answers := p.GenerateResources(p.OriginData.answers, int(p.OriginData.header.AC))
	res = append(res, answers...)
	authority := p.GenerateResources(p.OriginData.authority, int(p.OriginData.header.NSC))
	res = append(res, authority...)
	additional := p.GenerateResources(p.OriginData.additional, int(p.OriginData.header.AR))
	res = append(res, additional...)
	p.BinaryData = res
	return *p
}

func (p *GeneratePacket) GenerateHeader() []byte {
	var res []byte
	res = Uint16ToByte2AndAppend(res, p.OriginData.header.ID)
	res = Uint16ToByte2AndAppend(res, p.FlagsToBinary())
	res = Uint16ToByte2AndAppend(res, p.OriginData.header.QC)
	res = Uint16ToByte2AndAppend(res, p.OriginData.header.AC)
	res = Uint16ToByte2AndAppend(res, p.OriginData.header.NSC)
	res = Uint16ToByte2AndAppend(res, p.OriginData.header.AR)
	return res
}

// uint16(16bits)=>byte,byte=>append
func Uint16ToByte2AndAppend(res []byte, a uint16) []byte {
	b1, b2 := Uint16ToByte2(a)
	res = append(res, b1, b2)
	return res
}

// 把bool换成uint16,再换成[]byte
func (p *GeneratePacket) FlagsToBinary() uint16 {
	// 构建16位的标志字段
	flags := uint16(0)

	// QR位：第0位
	if p.OriginData.header.Flags.QR {
		flags |= 1 << 15
	}

	// Opcode位：第1-4位
	flags |= uint16(p.OriginData.header.Flags.OpCode) << 11

	// AA位：第5位
	if p.OriginData.header.Flags.AA {
		flags |= 1 << 10
	}

	// TC位：第6位
	if p.OriginData.header.Flags.TC {
		flags |= 1 << 9
	}

	// RD位：第7位
	if p.OriginData.header.Flags.RD {
		flags |= 1 << 8
	}

	// RA位：第8位
	if p.OriginData.header.Flags.RA {
		flags |= 1 << 7
	}

	// Z位：第9位（保留字段，设置为0）
	if p.OriginData.header.Flags.Z {
		flags |= 1 << 4
	}
	// AD: 10
	if p.OriginData.header.Flags.AD {
		flags |= 1 << 3
	}
	if p.OriginData.header.Flags.CD {
		flags |= 1 << 2
	}
	// Response Code位：第11-15位
	flags |= uint16(p.OriginData.header.Flags.RCode)

	return flags
}

func (p *GeneratePacket) GenerateQuery() []byte {
	var res []byte
	if len(p.OriginData.queries) == 0 {
		return nil
	}
	res = append(res, p.GenerateDomainName(p.OriginData.queries[0].QName)...) // 默认1个查询
	res = Uint16ToByte2AndAppend(res, uint16(p.OriginData.queries[0].QType))
	res = Uint16ToByte2AndAppend(res, uint16(p.OriginData.queries[0].QClass))
	return res
}

func (p *GeneratePacket) GenerateDomainName(domain string) []byte {
	var res []byte

	labels := strings.Split(domain, ".") // 分割域名字符串为标签

	for _, label := range labels {
		labelLength := byte(len(label))
		res = append(res, labelLength)      // 添加标签长度
		res = append(res, []byte(label)...) // 添加标签内容
	}

	// 添加域名末尾的零字节
	res = append(res, 0)

	return res
}

// 多个记录生成
func (p *GeneratePacket) GenerateResources(arg []packetResource, count int) []byte {
	var res []byte
	i := 0
	for i < count {
		res = append(res, p.GenerateResource(arg[i])...)
		i += 1
	}
	return res
}

// 单个记录生成
func (p *GeneratePacket) GenerateResource(arg packetResource) []byte {
	var res []byte
	name := p.GenerateDomainName(arg.Name)
	res = append(res, name...)
	res = Uint16ToByte2AndAppend(res, uint16(arg.Type))
	res = Uint16ToByte2AndAppend(res, uint16(arg.Class))
	res = append(res, Uint32ToByte4(arg.TTL)...)
	res = Uint16ToByte2AndAppend(res, arg.ReLength)
	res = append(res, p.GenerateRecordData(arg.RData)...)
	return res
}

func (p *GeneratePacket) GenerateRecordData(arg packetRecordData) []byte {
	var res []byte
	res = append(res, arg.A_IP[0:3]...)
	res = append(res, []byte(arg.NS_Name)...)
	res = append(res, []byte(arg.CNAME_Name)...)
	res = append(res, p.GenerateMXRecordData(arg)...)
	res = append(res, arg.AAAA_IP[0:15]...)
	return res
}
func (p *GeneratePacket) GenerateMXRecordData(arg packetRecordData) []byte {
	var res []byte
	res = Uint16ToByte2AndAppend(res, arg.MX.MX_Preference)
	res = append(res, []byte(arg.MX.MX_Name)...)
	return res
}
