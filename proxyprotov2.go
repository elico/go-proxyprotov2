package proxyprot

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

/// magic octet prefix for PROXY protocol version 1
const Proxy1p0magic = "PROXY "

/// magic octet prefix for PROXY protocol version 2
const Proxy2p0magic = "\x0D\x0A\x0D\x0A\x00\x0D\x0A\x51\x55\x49\x54\x0A"

func ParseProxy(header []byte) (src, dst *net.TCPAddr, err error) {
	headerString := string(header)

	if strings.HasPrefix(headerString, Proxy1p0magic) && strings.Contains(headerString, "\n") {
		src, dst, err := ProxyV1(headerString)
		if err != nil {
			fmt.Println("not a proxy v1")
			fmt.Println(err)
		}
		return src, dst, nil
	}

	if strings.HasPrefix(headerString, Proxy2p0magic) && len(headerString) > 107 {
		_ = 1
	} else {
		_ = err
	}

	return nil, nil, fmt.Errorf("Invalid PROXY PROTOCOL header")
}

func ProxyV1(header string) (src, dst *net.TCPAddr, err error) {
	if strings.HasPrefix(header, Proxy1p0magic) && strings.Contains(header, "\r\n") && strings.Index(header, "\r\n") < 108 {
		/*
				- TCP/IPv4 :
			      "PROXY TCP4 255.255.255.255 255.255.255.255 65535 65535\r\n"
			    => 5 + 1 + 4 + 1 + 15 + 1 + 15 + 1 + 5 + 1 + 5 + 2 = 56 chars

			  - TCP/IPv6 :
			      "PROXY TCP6 ffff:f...f:ffff ffff:f...f:ffff 65535 65535\r\n"
			    => 5 + 1 + 4 + 1 + 39 + 1 + 39 + 1 + 5 + 1 + 5 + 2 = 104 chars

			  - unknown connection (short form) :
			      "PROXY UNKNOWN\r\n"
			    => 5 + 1 + 7 + 2 = 15 chars

			  - worst case (optional fields set to 0xff) :
			      "PROXY UNKNOWN ffff:f...f:ffff ffff:f...f:ffff 65535 65535\r\n"
			    => 5 + 1 + 7 + 1 + 39 + 1 + 39 + 1 + 5 + 1 + 5 + 2 = 107 chars
		*/
		// There are couple edge cases which should not be here at all(port 0 and ips 0.0.0.0)
		// Also the case which composed of two IP addresses types while it can happen in real like on a v6 to v4 NAT/PROXY

		headerArr := strings.Split(header[:strings.Index(header, "\r\n")], " ")
		fmt.Println(headerArr)
		if len(headerArr) > 4 {
			switch headerArr[1] {
			case "TCP4":
			case "TCP6":
			default:
				return nil, nil, fmt.Errorf("Unhandled address type: %s", headerArr[1])
			}
			// Parse out the source address
			ip := net.ParseIP(headerArr[2])
			if ip == nil {
				return nil, nil, fmt.Errorf("Invalid source ip: %s", headerArr[2])
			}
			port, err := strconv.ParseUint(headerArr[4], 10, 16)
			if err != nil {
				return nil, nil, fmt.Errorf("Invalid source port: %s", headerArr[4])
			}

			srcAddr := &net.TCPAddr{IP: ip, Port: int(port)}

			// Parse out the destination address
			ip = net.ParseIP(headerArr[3])
			if ip == nil {
				return nil, nil, fmt.Errorf("Invalid destination ip: %s", headerArr[3])
			}
			port, err = strconv.ParseUint(headerArr[5], 10, 16)
			if err != nil {
				return nil, nil, fmt.Errorf("Invalid destination port: %s", headerArr[5])
			}
			dstAddr := &net.TCPAddr{IP: ip, Port: int(port)}

			return srcAddr, dstAddr, nil
		}
	}

	return nil, nil, err
}

/*
func ProxyV2(header string) (src, dst *net.TCPAddr ,err error){
	if strings.HasPrefix(headerString, Proxy1p0magic) && strings.Contains(headerString, "\n") {
		headerArr = strings.split(header)


		v = 1
	}
	return nil, nil, fmt.Errorf("Invalid ProxyProtocolV2")
}
*/
