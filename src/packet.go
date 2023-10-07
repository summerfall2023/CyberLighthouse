package src

import (
	"fmt"
	"strings"
)

// byte 8 bits, uint16 1位 4 bits
type packet struct {
	header     packetHeader // 12 bytes
	queries    []packetQuery
	answers    []packetResource
	authority  []packetResource
	additional []packetResource
}

// header-----------------------------------------------------
type packetHeader struct {
	ID    uint16           // 2 bytes
	Flags packetHeaderFlag // 2 bytes
	QC    uint16           // query count 2 bytes
	AC    uint16           // answer count 2 bytes
	NSC   uint16           // name server count 2 bytes
	AR    uint16           // additional record count 2 bytes
}
type packetHeaderFlag struct {
	QR     bool       // 1 bit,query = 0,reaponse = 1
	OpCode OpCodeType // 4 bits
	AA     bool       // A 1 bit
	TC     bool       // TrunCaton 1 bit
	RD     bool       // recursive desired 1 bit
	RA     bool       // A recursive available 1 bit
	Z      bool       // 保留字段 1 bit
	AD     bool       // A answer athenticated
	CD     bool       // A unacceptable
	RCode  RCodeType  //4 bits
}

type OpCodeType uint16

const (
	STANDARD_QUERY OpCodeType = 0
	INVERSE_QUERY  OpCodeType = 1
	STATUS         OpCodeType = 2
	RESERVED       OpCodeType = 3
)

type RCodeType uint16

const (
	NO_ERROR        RCodeType = 0
	FORMAT_ERROR    RCodeType = 1
	SERVER_FAILURE  RCodeType = 2
	NAME_ERROR      RCodeType = 3
	NOT_IMPLEMENTED RCodeType = 4
	REFUSED         RCodeType = 5
)

// query----------------------------------------------------
type packetQuery struct {
	QName  string         // 可变长，以00结尾
	QType  QueryType      // 2 bytes 4 uint16 （一个16进制数，四个16进制位）
	QClass QueryClassType // 2 bytes 4 uint16
}

type QueryType uint16

const (
	A                  QueryType = 1
	NS                 QueryType = 2
	CNAME              QueryType = 5
	MX                 QueryType = 15
	TXT                QueryType = 16
	AAAA               QueryType = 28
	NOT_SUPPORTED_TYPE QueryType = 0
)

type QueryClassType uint16 // (0-255)
const (
	IN                  QueryClassType = 1 // internet
	NOT_SUPPORTED_CLASS QueryClassType = 0
)

// record---------------------------------------------------
type packetResource struct {
	Name     string //
	Type     QueryType
	Class    QueryClassType
	TTL      uint32
	ReLength uint16           // 8 bits
	RData    packetRecordData // 可变长
}

type packetRecordData struct {
	A_IP       [4]byte
	NS_Name    string
	CNAME_Name string
	MX         MXRecordData
	AAAA_IP    [16]byte
	Others     []byte
}

type MXRecordData struct { // todo
	MX_Preference uint16
	MX_Name       string
}

const f string = "(PRINTF)"

// 整合输出
func (p *packet) OutputPacket() (string, error) {
	res := ""
	res = fmt.Sprintf("Domain Name System (%s)\n", p.OutputQR())
	fmt.Println(f, res)
	var err1, err2 error
	var res0 string
	res0, err1 = p.OutputHeader()
	if err1 != nil {
		return "", err1
	}
	res += res0
	fmt.Println(f, res)
	var res1 string
	res1, err2 = p.OutputQuery()
	if err2 != nil {
		return "", err2
	}
	res += res1
	fmt.Println(f, res)
	// RECORD
	// ANSWERS
	// p.header.AC != 0x0
	if len(p.answers) != 0 {
		res += ("	ANSWERS:\n")
		r3, err3 := p.OutputResources(p.answers, int(p.header.AC))
		if err3 != nil {
			return "", err3
		}
		res += r3
	} else {
		res += "	No Answer\n"
	}
	fmt.Println(f, res)
	// AUTHORITY
	// p.header.NSC != 0x0
	if len(p.authority) != 0 {
		res += ("	AUTHORITY:\n")
		r4, err4 := p.OutputResources(p.authority, int(p.header.NSC))
		if err4 != nil {
			return "", err4
		}
		res += r4
	} else {
		res += "	No Authority\n"
	}
	fmt.Println(f, res)
	// ADDITIONAL
	// p.header.AR != 0x0
	if len(p.additional) != 0 {
		res += ("	ADDITIONAL:\n")
		r5, err5 := p.OutputResources(p.additional, int(p.header.AR))
		if err5 != nil {
			return "", err5
		}
		res += r5
	} else {
		res += "	No Additional\n"
	}
	fmt.Println(f, res)
	return res, nil
}

func (p *packet) OutputQR() string {
	var res string
	switch p.header.Flags.QR {
	case false:
		res += "query"
	case true:
		res += "response"
	default:
		res += "wrong"
	}
	return res
}
func (p *packet) OutputHeader() (string, error) {
	res := ""
	res += fmt.Sprintf("	Transaction ID: %d\n", p.header.ID)
	res += ("	Flags: ")
	res += ("\n")
	res += fmt.Sprintf("		Response: Message is a %s\n", p.OutputQR())

	var opcode string
	switch p.header.Flags.OpCode {
	case STANDARD_QUERY:
		opcode = "standard query"
	case INVERSE_QUERY:
		opcode = "inverse query"
	case STATUS:
		opcode = "status query"
	default:
		opcode = "reserved"
	}
	res += fmt.Sprintf("		Opcode: %s (%d)\n", opcode, p.header.Flags.OpCode)

	switch p.header.Flags.AA {
	case true:
		res += "		Authoritative: Server is an authority for domain\n"
	case false:
		res += ("		Authoritative: Server is not an authority for domain\n")
	default:
		res += ""
	}

	switch p.header.Flags.TC {
	case true:
		res += ("		Truncated: Message is truncated\n")
	case false:
		res += ("		Truncated: Message is not truncated\n")
	default:
		res += ""
	}

	switch p.header.Flags.RD {
	case true:
		res += ("		Recursion desired: Do query recursively\n")
	case false:
		res += ("		Recursion not desired: Do not query recursively\n")
	default:
		res += ""
	}

	switch p.header.Flags.RA {
	case true:
		res += ("		Recursion available: Server can do recursive queries\n")
	case false:
		res += ("		Recursion unavailable: Server can not do recursive queries\n")
	default:
		res += ""
	}

	switch p.header.Flags.Z {
	case false:
		res += ("		Z:0\n")
	default:
		err := fmt.Errorf(ERROR_HEADER_FLAG_Z)
		return "", err
	}

	switch p.header.Flags.AD {
	case true:
		res += ("		Answer authenticated: Answer/authority portion was authenticated by the server\n")
	case false:
		res += ("		Answer authenticated: Answer/authority portion was not authenticated by the server\n")
	default:
		res += ""
	}

	switch p.header.Flags.CD {
	case true:
		res += ("		Authenticated data: Acceptable\n")
	case false:
		res += ("		Non-authenticated data: Unacceptable\n")
	default:
		res += ""
	}

	switch p.header.Flags.RCode {
	case NO_ERROR:
		res += "		Reply code: No error (0)\n"
	case FORMAT_ERROR:
		res += "		Reply code: Format error (1)\n"
	case SERVER_FAILURE:
		res += "		Reply code: Server failure (2)\n"
	case NAME_ERROR:
		res += "		Reply code: Name Error (3)\n"
	case NOT_IMPLEMENTED:
		res += "		Reply code: Not Implemented (4)\n"
	case REFUSED:
		res += "		Reply code: Refused (5)\n"
	default:
		res += ""
	}

	res += fmt.Sprintf("	Query: %d\n", p.header.QC)
	res += fmt.Sprintf("	Answer: %d\n", p.header.AC)
	res += fmt.Sprintf("	Authority: %d\n", p.header.NSC)
	res += fmt.Sprintf("	Additional: %d\n", p.header.AR)

	return res, nil
}

func (p *packet) OutputQuery() (string, error) {
	res := ("	Query:\n")
	qtype, err1 := p.OutputType(p.queries[0].QType)
	if err1 != nil {
		return "", err1
	}
	class, err2 := p.OutputClass(p.queries[0].QClass)
	if err2 != nil {
		return "", err2
	}
	res += fmt.Sprintf("		%s: type %s class %s\n", p.queries[0].QName, qtype, class)
	return res, nil
}
func (p *packet) OutputType(t QueryType) (string, error) {
	var qtype string
	switch t {
	case A:
		qtype = "A"
	case NS:
		qtype = "NS"
	case CNAME:
		qtype = "CNAME"
	case MX:
		qtype = "MX"
	case TXT:
		qtype = "TXT"
	case AAAA:
		qtype = "AAAA"
	default:
		qtype = fmt.Sprintf("%v", t)
	}
	return qtype, nil

}
func (p *packet) OutputClass(c QueryClassType) (string, error) {
	var class string
	switch c {
	case IN:
		class = "IN"
	default:
		err2 := fmt.Errorf(ERROR_QUERY_CLASS_UNSUPPORTED)
		return "", err2
	}
	return class, nil
}
func (p *packet) OutputResources(packetResource []packetResource, count int) (string, error) {
	fmt.Println("OUTPUT_RESOURCES START+++++++++++++++++++++++++++++")
	i := 0
	res := ""
	for i < count {
		fmt.Println("RESOURCE START++++++++++++++++++++++++++++")
		r, err := p.OutputResource(packetResource[i])
		if err != nil {
			return "", err
		}
		res += r
		fmt.Println(f, "RESOURCE", res)
		i += 1
	}
	return res, nil
}

func (p *packet) OutputResource(packetResource packetResource) (string, error) {
	var res string
	if packetResource.Name == " " {
		res = "<root>"
	} else {
		res = fmt.Sprintf("		Name: %s \n", packetResource.Name)
	}
	t, err1 := p.OutputType(packetResource.Type)
	if err1 != nil {
		return "", err1
	}
	res += fmt.Sprintf("		Type: %s \n", t)
	c, err2 := p.OutputClass(packetResource.Class)
	if err2 != nil {
		return "", err1
	}
	res += fmt.Sprintf("		Class: %s \n", c)
	res += fmt.Sprintf("		Time to live: %d\n", packetResource.TTL)
	res += fmt.Sprintf("		Data length: %d\n", packetResource.ReLength)
	res += fmt.Sprintf("		Resource: %s\n", p.OutputRecordData(packetResource))
	return res, nil
}
func (p *packet) OutputRecordData(packetResource packetResource) string {
	res := ""
	switch packetResource.Type {
	case NS:
		res += packetResource.RData.NS_Name
	case CNAME:
		res += packetResource.RData.CNAME_Name
	case A:
		res += bytesToDotSeparatedString(packetResource.RData.A_IP[0:3])
		//res += string(packetResource.RData.A_IP[0:3])
	case AAAA:
		res += bytesToDotSeparatedString(packetResource.RData.AAAA_IP[0:15])
		//res += string(packetResource.RData.AAAA_IP[0:15])
	case MX:
		res += fmt.Sprintf("%d", packetResource.RData.MX.MX_Preference)
		res += packetResource.RData.MX.MX_Name
	default:
		res += string(packetResource.RData.Others)
	}
	return res
}
func bytesToDotSeparatedString(data []byte) string {
	var parts []string
	for _, b := range data {
		parts = append(parts, fmt.Sprintf("%d", b))
	}
	return strings.Join(parts, ".")
}
