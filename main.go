package main

import (
	"encoding/binary"
	"encoding/hex"
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

type DNSFlags struct {
	QR     byte // (1bit) question/response
	Opcode byte // (4bit) query type
	AA     byte // (1bit) authority thing
	TC     byte // (1bit) truncated
	RD     byte // (1bit) recursion desired
	RA     byte // (1bit) recursion available
	Z      byte // (3bit) reserved
	RCODE  byte // (4bit) errors
}
type DNSHeader struct {
	ID      []byte   // (16bits) request id
	QDCOUNT []byte   // (16bits) number of entries of the question
	ANCOUNT []byte   // (16bits) number of entries of the answer
	NSCOUNT []byte   // (16bits) authority records
	ARCOUNT []byte   // (16bits) additionnal records
	flags   DNSFlags // (16bits) flags
}
type DNSQuestion struct {
	QNAME  []byte // domain
	QTYPE  []byte // (16bits) type
	QCLASS []byte // (16bits) class
}
type DNSAnswer struct {
	NAME     []byte // (16bits)
	TYPE     []byte // (16bits) same as qtype
	CLASS    []byte // (16bits) same as qclass
	TTL      []byte // (32bits unsigned) time to live in seconds
	RDLENGTH []byte // (16bits) length of rddata
	RDDATA   []byte // (32bits) ip address
}
type DNSMessage struct {
	header   DNSHeader
	question DNSQuestion
	answer   DNSAnswer
}

func (dnsheader *DNSHeader) ApplyTo(message []byte) []byte {
	copy(message[0:2], dnsheader.ID)

	// TODO: flags

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

	// idk how it works but it does so
	firstByteNum, _ := strconv.ParseUint(firstByte, 2, 8)
	firstByteByte := make([]byte, 8)
	binary.LittleEndian.PutUint64(firstByteByte, firstByteNum)

	lastByteNum, _ := strconv.ParseUint(lastByte, 2, 8)
	lastByteByte := make([]byte, 8)
	binary.LittleEndian.PutUint64(lastByteByte, lastByteNum)

	//fmt.Printf("\n%d, %d, %d, %d, %d, %d, %d, %d\n", dnsheader.flags.QR, dnsheader.flags.Opcode, dnsheader.flags.AA, dnsheader.flags.TC, dnsheader.flags.RD, dnsheader.flags.RA, dnsheader.flags.Z, dnsheader.flags.RCODE)

	message[2] = firstByteByte[0]
	message[3] = lastByteByte[0]

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
		message := make([]byte, 512)
		rlen, remote, err := ln.ReadFromUDP(message[:])
		if err != nil {
			panic(err)
		}

		var dnsMessage DNSMessage

		dnsMessage.Set(message, rlen)

		fmt.Printf("\n%s: type %s, class %s From %v\n", dnsMessage.question.GetDomain(), getType(dnsMessage.question.QTYPE), getClass(dnsMessage.question.QCLASS), remote)
		fmt.Println("QR (0: Question 1: Answer):", dnsMessage.header.flags.QR)
		fmt.Println("Opcode:", dnsMessage.header.flags.Opcode)
		fmt.Println("RD (1: recursion):", dnsMessage.header.flags.RD)
		fmt.Println("RA:", dnsMessage.header.flags.RA)
		fmt.Println("Erreur:", dnsMessage.header.flags.RCODE)

		//fmt.Println(insertNth(hex.EncodeToString(message[:rlen]), 2))
		fmt.Println("QUESTION", message[:rlen])

		go sendResponse(ln, remote, dnsMessage, rlen)
	}
}

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr, question DNSMessage, rlen int) {

	message := make([]byte, rlen)
	question.header.ANCOUNT = []byte{0, 1}
	question.header.flags.QR = byte(1)
	question.header.flags.AA = byte(0)
	question.header.flags.RA = byte(1)
	message = question.header.ApplyTo(message)
	message = question.question.ApplyTo(message)

	fmt.Println("ANSWER  ", message)

	data, err := hex.DecodeString("00028180000100010000000006676f6f676c6503636f6d0000010001c00c00010001000000b10004d83ad64e")
	if err != nil {
		panic(err)
	}
	//data[6:8] = byte(192)
	data[len(data)-4] = byte(192)
	data[len(data)-3] = byte(168)
	data[len(data)-2] = byte(1)
	data[len(data)-1] = byte(16)

	rlen2 := len(data)
	transactionid := data[0:2]

	flags := data[2:4]

	questions := data[4:6]
	data[6] = byte(0)
	data[7] = byte(1)
	answer_rrs := data[6:8]
	authority_rrs := data[8:10]
	additional_rrs := data[10:12]

	domain := data[12 : rlen-4]

	types := data[rlen2-4 : rlen2-2]
	class := data[rlen2-2 : rlen2]

	Use(flags, transactionid, questions, answer_rrs, authority_rrs, additional_rrs, domain, types, class, data)

	fmt.Println("HARDCODE", data)

	_, err = conn.WriteToUDP(message, addr)
	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
	}
}
