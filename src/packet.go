package src

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

// type QRType bool
// const (
// 	QUERY QRType = true
// 	answer QRType = false
// )

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
	QType  QueryType      // 4 bytes 2 uint16 （一个16进制数，两个16进制位，0-255）
	QClass QueryClassType // 4 bytes 2 uint16
}

type QueryType uint16 // (0-255)
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
	AAAA_IP    [8]uint16
	originData []byte // todo
}

type MXRecordData struct { // todo
	Preference uint16
	Name       string
}
