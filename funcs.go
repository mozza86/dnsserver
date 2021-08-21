package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

func getType(b []byte) string {
	var typeInt uint16 = binary.BigEndian.Uint16(b)
	value := typeMap[int(typeInt)]
	return value
}

func getClass(b []byte) string {
	switch b[1] {
	case 1:
		return "IN"
	case 2:
		return "Cs"
	case 3:
		return "Ch"
	case 4:
		return "Hs"
	default:
		return "Unknown"
	}
}

func getClassLong(b []byte) string {
	switch b[1] {
	case 1:
		return "Internet"
	case 2:
		return "Class Csnet"
	case 3:
		return "Chaos"
	case 4:
		return "Hesiod"
	default:
		return "Unknown"
	}
}
func Use(vals ...interface{}) {
	for _, val := range vals {
		_ = val
	}
}

func parseDomain(b []byte) []string {
	lenB := len(b)
	var domainSlice []string
	var curByte int = 0 //byte position of the first string length
	for {
		labelLen := int(b[curByte])                              //size of the label
		labelString := string(b[curByte+1 : curByte+1+labelLen]) //get value of the label
		curByte = curByte + 1 + labelLen                         //get the byte position of the next label length
		if curByte >= lenB {
			break
		}
		domainSlice = append(domainSlice, labelString)
	}
	return domainSlice
}
func convertByteSliceToStr(b []byte) string { // TODO: a check si ca renvoie pas de la merde
	s := make([]string, len(b))
	for i := range b {
		s[i] = strconv.Itoa(int(b[i]))
	}
	return strings.Join(s, "")
}
func convert8BitsToByte(b1 byte, b2 byte, b3 byte, b4 byte, b5 byte, b6 byte, b7 byte, b8 byte) byte { // TODO: a check si ca renvoie pas de la merde
	bitString := convertByteSliceToStr([]byte{b1, b2, b3, b4, b5, b6, b7, b8})
	bytestr, err := strconv.ParseUint(bitString, 2, 8)
	if err != nil {
		fmt.Println(err.Error())
	}
	return byte(bytestr)
}
func parseFlags(b []byte) (byte, byte, byte, byte, byte, byte, byte, byte) {

	QR := b[0] & 128
	Opcode := convert8BitsToByte(byte(0), byte(0), byte(0), byte(0), b[0]&64, b[0]&32, b[0]&16, b[0]&8)
	AA := b[0] & 4
	TC := b[0] & 2
	RD := b[0] & 1
	RA := b[1] & 128
	Z := convert8BitsToByte(byte(0), byte(0), byte(0), byte(0), byte(0), b[1]&64, b[1]&32, b[1]&16)
	rcode := convert8BitsToByte(byte(0), byte(0), byte(0), byte(0), b[1]&8, b[1]&4, b[1]&2, b[1]&1)

	return QR, Opcode, AA, TC, RD, RA, Z, rcode
}

func insertNth(s string, n int) string {
	var buffer bytes.Buffer
	var n_1 = n - 1
	var l_1 = len(s) - 1
	for i, rune := range s {
		buffer.WriteRune(rune)
		if i%n == n_1 && i != l_1 {
			buffer.WriteRune(' ')
		}
	}
	return buffer.String()
}

func getFlagsBytes(dnsheader *DNSHeader) (byte, byte) {
	var firstByte string
	var lastByte string

	if dnsheader.flags.QR == byte(1) {
		firstByte += "1"
	} else {
		firstByte += "0"
	}
	if dnsheader.flags.Opcode == byte(0) {
		firstByte += "0000"
	} else if dnsheader.flags.Opcode == byte(1) {
		firstByte += "0001"
	} else if dnsheader.flags.Opcode == byte(2) {
		firstByte += "0010"
	} else if dnsheader.flags.Opcode == byte(3) {
		firstByte += "0011"
	} else if dnsheader.flags.Opcode == byte(4) {
		firstByte += "0100"
	}
	if dnsheader.flags.AA == byte(1) {
		firstByte += "1"
	} else {
		firstByte += "0"
	}
	if dnsheader.flags.TC == byte(1) {
		firstByte += "1"
	} else {
		firstByte += "0"
	}
	if dnsheader.flags.RD == byte(1) {
		firstByte += "1"
	} else {
		firstByte += "0"
	}
	if dnsheader.flags.RA == byte(1) {
		lastByte += "1"
	} else {
		lastByte += "0"
	}
	if dnsheader.flags.Z == byte(1) {
		lastByte += "001"
	} else {
		lastByte += "000"
	}
	if dnsheader.flags.RCODE == byte(0) {
		lastByte += "0000"
	} else if dnsheader.flags.RCODE == byte(1) {
		lastByte += "0001"
	} else if dnsheader.flags.RCODE == byte(2) {
		lastByte += "0010"
	} else if dnsheader.flags.RCODE == byte(3) {
		lastByte += "0011"
	} else if dnsheader.flags.RCODE == byte(4) {
		lastByte += "0100"
	} else if dnsheader.flags.RCODE == byte(5) {
		lastByte += "0101"
	} else if dnsheader.flags.RCODE == byte(6) {
		lastByte += "0110"
	} else if dnsheader.flags.RCODE == byte(7) {
		lastByte += "0111"
	} else if dnsheader.flags.RCODE == byte(8) {
		lastByte += "1000"
	} else if dnsheader.flags.RCODE == byte(9) {
		lastByte += "1001"
	} else if dnsheader.flags.RCODE == byte(10) {
		lastByte += "1010"
	} else if dnsheader.flags.RCODE == byte(11) {
		lastByte += "1011"
	} else if dnsheader.flags.RCODE == byte(12) {
		lastByte += "1100"
	} else if dnsheader.flags.RCODE == byte(13) {
		lastByte += "1101"
	} else if dnsheader.flags.RCODE == byte(14) {
		lastByte += "1110"
	} else if dnsheader.flags.RCODE == byte(15) {
		lastByte += "1111"
	}

	//it works in a hacky way
	firstByteNum, _ := strconv.ParseUint(firstByte, 2, 8)
	firstByteByte := make([]byte, 8)
	binary.LittleEndian.PutUint64(firstByteByte, firstByteNum)

	lastByteNum, _ := strconv.ParseUint(lastByte, 2, 8)
	lastByteByte := make([]byte, 8)
	binary.LittleEndian.PutUint64(lastByteByte, lastByteNum)

	return firstByteByte[0], lastByteByte[0]
}
