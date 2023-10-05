package src

import (
	"fmt"
)

const BYTE_TO_UINT16_OFFSET = 8
const BYTE_TO_UINT32_OFFSET = 16

// 把一个byte转换成uint16时偏移量是8个二进制位

// 两个字节（ 8 bits/byte）转换成一个uint16 (4 bits/1位)（两位）
func Byte2ToUint16(a []byte) (uint16, error) {
	if len(a) == 0 {
		return 0, nil
	}
	if len(a) == 1 {
		return uint16(a[0]), nil
	}
	if len(a) == 2 {
		return (uint16(a[0]) << BYTE_TO_UINT16_OFFSET) + uint16(a[1]), nil
	}
	err := fmt.Errorf("The lenth of assingment should be less than 2")
	return 0, err
}

// byte 8 bits  Uint32 8 bits
func Byte2ToUint32(a []byte) (uint32, error) {
	if len(a) == 2 {
		return uint32(a[0])<<BYTE_TO_UINT16_OFFSET*2 + uint32(a[1]), nil
	} else {
		err := fmt.Errorf("ByteToUint32 error:The lenth expected 2")
		return 0, err
	}
}
func Byte4ToUint32(a []byte) (uint32, error) {
	if len(a) == 0 {
		return 0, nil
	}
	if len(a) == 1 {
		return uint32(a[0]), nil
	}
	if len(a) == 2 {
		return uint32(a[0])<<BYTE_TO_UINT32_OFFSET + uint32(a[1]), nil
	}
	if len(a) == 3 {
		b, err1 := Byte2ToUint32(a[0:2])
		if err1 != nil {
			return 0, err1
		}
		return b<<BYTE_TO_UINT32_OFFSET + uint32(a[2]), nil
	}
	if len(a) == 4 {
		c, err2 := Byte2ToUint32(a[0:2])
		if err2 != nil {
			return 0, err2
		}
		d, err3 := Byte2ToUint32(a[2:4])
		if err3 != nil {
			return 0, err3
		}
		return c<<BYTE_TO_UINT32_OFFSET + d, nil
	}
	err4 := fmt.Errorf("Byte4ToUint32 error:The length expected less than 4")
	return 0, err4
}
