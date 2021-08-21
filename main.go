package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
)

func (dnsheader *DNSHeader) ApplyTo(message []byte) []byte {
	copy(message[0:2], dnsheader.ID)
	message[2], message[3] = getFlagsBytes(dnsheader)
	copy(message[4:6], dnsheader.QDCOUNT)
	copy(message[6:8], dnsheader.ANCOUNT)
	copy(message[8:10], dnsheader.NSCOUNT)
	copy(message[10:12], dnsheader.ARCOUNT)
	return message
}
func (dnsquestion *DNSQuestion) ApplyTo(message []byte) []byte {
	copy(message[12:], dnsquestion.QNAME)
	lenQNAME := len(dnsquestion.QNAME)
	copy(message[12+lenQNAME:], dnsquestion.QTYPE)
	copy(message[12+lenQNAME+2:], dnsquestion.QCLASS)
	return message
}
func (dnsanswer *DNSAnswer) ApplyTo(message []byte) []byte {
	message = append(message, []byte{192, 12}...)
	message = append(message, []byte{0, 0}...)
	copy(message[len(message)-2:], dnsanswer.TYPE)
	message = append(message, []byte{0, 0}...)
	copy(message[len(message)-2:], dnsanswer.CLASS)
	message = append(message, []byte{0, 0, 0, 0}...)
	copy(message[len(message)-4:], dnsanswer.TTL)
	message = append(message, []byte{0, 0}...)
	copy(message[len(message)-2:], dnsanswer.RDLENGTH)
	rdlen := binary.BigEndian.Uint16(dnsanswer.RDLENGTH)
	message = append(message, make([]byte, rdlen)...)
	copy(message[len(message)-int(rdlen):], dnsanswer.RDDATA)

	return message
}
func (dnsmsg *DNSMessage) Set(message []byte, rlen int) {
	dnsmsg.header.ID = message[0:2]
	dnsmsg.header.flags.Set(message[2:4])
	dnsmsg.header.QDCOUNT = message[4:6]
	dnsmsg.header.ANCOUNT = message[6:8]
	dnsmsg.header.NSCOUNT = message[8:10]
	dnsmsg.header.ARCOUNT = message[10:12]
	dnsmsg.question.QNAME = message[12 : rlen-4]
	dnsmsg.question.QTYPE = message[rlen-4 : rlen-2]
	dnsmsg.question.QCLASS = message[rlen-2 : rlen]
}
func (flags *DNSFlags) Set(b []byte) {
	flags.QR = convert8BitsToByte(byte(0), byte(0), byte(0), byte(0), byte(0), byte(0), byte(0), b[0]&128)   //query response
	flags.Opcode = convert8BitsToByte(byte(0), byte(0), byte(0), byte(0), b[0]&64, b[0]&32, b[0]&16, b[0]&8) //type of query or response (0, 1, 2)
	flags.AA = convert8BitsToByte(byte(0), byte(0), byte(0), byte(0), byte(0), byte(0), byte(0), b[0]&4)     //authoritative answer
	flags.TC = convert8BitsToByte(byte(0), byte(0), byte(0), byte(0), byte(0), byte(0), byte(0), b[0]&2)     //truncated
	flags.RD = convert8BitsToByte(byte(0), byte(0), byte(0), byte(0), byte(0), byte(0), byte(0), b[0]&1)     //recursion desired
	flags.RA = convert8BitsToByte(byte(0), byte(0), byte(0), byte(0), byte(0), byte(0), byte(0), b[1]&128)   //recursion available
	flags.Z = convert8BitsToByte(byte(0), byte(0), byte(0), byte(0), byte(0), b[1]&64, b[1]&32, b[1]&16)     //reserved (should be 000)
	flags.RCODE = convert8BitsToByte(byte(0), byte(0), byte(0), byte(0), b[1]&8, b[1]&4, b[1]&2, b[1]&1)     //status error
}
func (flags DNSFlags) Get() (byte, byte, byte, byte, byte, byte, byte, byte) {
	return flags.QR, flags.Opcode, flags.AA, flags.TC, flags.RD, flags.RA, flags.Z, flags.RCODE
}
func (question DNSQuestion) GetDomain() string {
	b := question.QNAME
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
	return strings.Join(domainSlice, ".")
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
		var dnsMessage DNSMessage
		message := make([]byte, 512)

		rlen, remote, err := ln.ReadFromUDP(message[:])
		if err != nil {
			panic(err)
		}

		dnsMessage.Set(message, rlen)

		fmt.Printf("\n%s: type %s, class %s From %v\n", dnsMessage.question.GetDomain(), getType(dnsMessage.question.QTYPE), getClass(dnsMessage.question.QCLASS), remote)

		//fmt.Println("Opcode:", dnsMessage.header.flags.Opcode)
		//fmt.Println("RD (1: recursion):", dnsMessage.header.flags.RD)
		//fmt.Println("Erreur:", dnsMessage.header.flags.RCODE)

		//fmt.Println(insertNth(hex.EncodeToString(message[:rlen]), 2))
		//fmt.Println("QUESTION", message[:rlen])

		go sendResponse(ln, remote, dnsMessage, rlen)
	}
}

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr, questionmsg DNSMessage, rlen int) {

	message := make([]byte, rlen)
	questionmsg.header.ANCOUNT = []byte{0, 1}
	questionmsg.header.flags.QR = byte(1)
	questionmsg.header.flags.AA = byte(0)
	questionmsg.header.flags.RA = byte(1)

	//answer
	questionmsg.answer.NAME = questionmsg.question.QNAME
	questionmsg.answer.TYPE = questionmsg.question.QTYPE
	questionmsg.answer.CLASS = questionmsg.question.QCLASS
	questionmsg.answer.TTL = []byte{0, 0, 0, 177}

	if getType(questionmsg.question.QTYPE) == "PTR" {
		questionmsg.answer.RDDATA = []byte{127, 0, 0, 69}
	}

	if getType(questionmsg.question.QTYPE) == "A" {
		questionmsg.answer.RDDATA = []byte{192, 168, 1, 16}
	}
	/*if getType(questionmsg.question.QTYPE) == "AAAA" {
		questionmsg.answer.RDDATA = []byte{0x2a, 0x00, 0x14, 0x50, 0x40, 0x07, 0x80, 0x0b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x20, 0x0e}
	}*/
	if questionmsg.question.GetDomain() == "google.com" {
		if getType(questionmsg.question.QTYPE) == "A" {
			questionmsg.answer.RDDATA = []byte{216, 58, 213, 78}
		}
		if getType(questionmsg.question.QTYPE) == "AAAA" {
			questionmsg.answer.RDDATA = []byte{0x2a, 0x00, 0x14, 0x50, 0x40, 0x07, 0x80, 0x0b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x20, 0x0e}
		}
	}
	if questionmsg.question.GetDomain() == "discord.com" ||
		questionmsg.question.GetDomain() == "remote-auth-gateway.discord.gg" ||
		questionmsg.question.GetDomain() == "gateway.discord.gg" ||
		questionmsg.question.GetDomain() == "status.discord.com" ||
		questionmsg.question.GetDomain() == "cdn.discordapp.com" ||
		questionmsg.question.GetDomain() == "media.discordapp.net" {
		if getType(questionmsg.question.QTYPE) == "A" {
			questionmsg.answer.RDDATA = []byte{162, 159, 136, 232}
		}
	}

	rdlenI := len(questionmsg.answer.RDDATA)
	questionmsg.answer.RDLENGTH = []byte{uint8(rdlenI >> 8), uint8(rdlenI & 0xff)}

	message = questionmsg.header.ApplyTo(message)
	message = questionmsg.question.ApplyTo(message)
	message = questionmsg.answer.ApplyTo(message)

	//fmt.Println("ANSWER  ", message)

	data, err := hex.DecodeString("00028180000100010000000006676f6f676c6503636f6d0000010001c00c00010001000000b10004d83ad64e")
	if err != nil {
		panic(err)
	}
	Use(data)

	//fmt.Println("HARDCODE", data)

	_, err = conn.WriteToUDP(message, addr)
	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
	}
}
