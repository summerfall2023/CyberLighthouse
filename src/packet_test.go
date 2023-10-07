package src_test

import (
	"fmt"
	"os"
	"reflect"
	"src"
	"testing"
)

// var parserTestBinFile []string = []string{"../../bin/parse_test/response_A_google.com2.bin"}

var parserTestBinFile []string = []string{
	//"../../bin/parse_test/answer_A_google.com.bin",
	//"../../bin/parse_test/response_A_google.com2.bin",
	//"../../bin/parse_test/query_A_google.com.bin",
	"../../bin/parse_test/query_A_google.com2.bin",
	//"../../bin/parse_test/response_A_google.com.bin",
}

// go test -v
func TestParser(t *testing.T) {
	for _, f := range parserTestBinFile {
		dir, _ := os.Getwd()
		f = dir + f
		file, err := os.Open(f)
		if err != nil {
			t.Error(err, fmt.Sprintf("Read file %s failed.", f))

			continue
		}
		defer file.Close()

		buff := make([]byte, 1024)
		file.Read(buff)
		var pk src.ParsePacket
		pk.BinaryPacket = buff
		_, err = pk.ParsePacket()
		if err != nil {
			t.Error(err, "Parse failed.")
		} else {
			r, _ := pk.ParsedPacket.OutputPacket()
			fmt.Println("------------------------------")
			fmt.Println(r)
		}
	}
}

func TestGenerator(t *testing.T) {
	for _, f := range parserTestBinFile {
		dir, _ := os.Getwd()
		f = dir + f
		file, err := os.Open(f)
		if err != nil {
			t.Error(err, fmt.Sprintf("Read file %s failed.", f))
			continue
		}
		defer file.Close()

		buff := make([]byte, 1024)
		n, _ := file.Read(buff)
		var pk src.ParsePacket
		pk.BinaryPacket = buff
		_, err = pk.ParsePacket()
		if err != nil {
			t.Error(err, "Parse failed.")
		}

		var ge src.GeneratePacket
		ge.BinaryData = pk.BinaryPacket
		_ = ge.GeneratePacket()

		var pk2 src.ParsePacket
		pk2.BinaryPacket = ge.BinaryData
		_, err = pk2.ParsePacket()
		if err != nil {
			t.Error(err, "Parse2 failed.")
		}
		fmt.Println("------------------------------------------------")
		s1, _ := pk.ParsedPacket.OutputPacket()
		s2, _ := pk2.ParsedPacket.OutputPacket()
		fmt.Println(s1)
		fmt.Println(s2)
		fmt.Println("------------------------------------------------")
		if reflect.DeepEqual(pk.ParsedPacket, pk2.ParsedPacket) || s1 != s2 {
			t.Error(err, "Generate check failed. Result not the same.\n", buff[:n], "\n", ge.OriginData)
		}
	}
}
