package main

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
