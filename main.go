package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func Use(vals ...interface{}) {
	for _, val := range vals {
		_ = val
	}
}

func main() {
	ln, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: 53,
		IP:   net.ParseIP("0.0.0.0"),
	})
	if err != nil {
		fmt.Println("listen error: " + err.Error())
		return
	}

	fmt.Println("launching server on " + ln.LocalAddr().String())

	for {
		message := make([]byte, 512)
		rlen, remote, err := ln.ReadFromUDP(message[:])
		if err != nil {
			panic(err)
		}
		transactionid := message[0:2]

		flags := message[2:4]

		questions := message[4:6]
		answer_rrs := message[6:8]
		authority_rrs := message[8:10]
		additional_rrs := message[10:12]

		domain := message[12 : rlen-4]

		types := message[rlen-4 : rlen-2]
		class := message[rlen-2 : rlen]
		data := strings.TrimSpace(string(message[:rlen]))
		//fmt.Printf("\nreceived: %s from %s\n", data, remote)
		//fmt.Printf("\ntransactionid: %d, flags: %b, questions: %b, answer_rrs: %b, authority_rrs: %b, additional_rrs: %b, domain: %s, type: %s, class: %s\n", transactionid, flags, questions, answer_rrs, authority_rrs, additional_rrs, domain, getType(types), getClass(class))

		Use(remote, transactionid, questions, answer_rrs, authority_rrs, additional_rrs, data)
		fmt.Printf("\n%s: type %s, class %s\n", strings.Join(parseDomain(domain), "."), getType(types), getClass(class))
		fmt.Println(parseFlags(flags))

		//fmt.Printf("%08b", message[:rlen])
	}
}

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
