package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

var (
	typeMap = map[int]string{
		0:     "Reserved",
		1:     "A",
		2:     "NS",
		3:     "MD",
		4:     "MF",
		5:     "CNAME",
		6:     "SOA",
		7:     "MB",
		8:     "MG",
		9:     "MR",
		10:    "NULL",
		11:    "WKS",
		12:    "PTR",
		13:    "HINFO",
		14:    "MINFO",
		15:    "MX",
		16:    "TXT",
		17:    "RP",
		18:    "AFSDB",
		19:    "X25",
		20:    "ISDN",
		21:    "RT",
		22:    "NSAP",
		23:    "NSAP-PTR",
		24:    "SIG",
		25:    "KEY",
		26:    "PX",
		27:    "GPOS",
		28:    "AAAA",
		29:    "LOC",
		30:    "NXT",
		31:    "EID",
		32:    "NIMLOC",
		33:    "SRV",
		34:    "ATMA",
		35:    "NAPTR",
		36:    "KX",
		37:    "CERT",
		38:    "A6",
		39:    "DNAME",
		40:    "SINK",
		41:    "OPT",
		42:    "APL",
		43:    "DS",
		44:    "SSHFP",
		45:    "IPSECKEY",
		46:    "RRSIG",
		47:    "NSEC",
		48:    "DNSKEY",
		49:    "DHCID",
		50:    "NSEC3",
		51:    "NSEC3PARAM",
		52:    "TLSA",
		53:    "SMIMEA",
		54:    "Unassigned",
		55:    "HIP",
		56:    "NINFO",
		57:    "RKEY",
		58:    "TALINK",
		59:    "CDS",
		60:    "CDNSKEY",
		61:    "OPENPGPKEY",
		62:    "CSYNC",
		63:    "ZONEMD",
		64:    "SVCB",
		65:    "HTTPS",
		99:    "SPF",
		100:   "UINFO",
		101:   "UID",
		102:   "GID",
		103:   "UNSPEC",
		104:   "NID",
		105:   "L32",
		106:   "L64",
		107:   "LP",
		108:   "EUI48",
		109:   "EUI64",
		249:   "TKEY",
		250:   "TSIG",
		251:   "IXFR",
		252:   "AXFR",
		253:   "MAILB",
		254:   "MAILA",
		255:   "*",
		256:   "URI",
		257:   "CAA",
		258:   "AVC",
		259:   "DOA",
		260:   "AMTRELAY",
		32768: "TA",
		32769: "DLV",
		65535: "Reserved",
	}
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
