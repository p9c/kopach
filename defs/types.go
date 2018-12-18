package defs

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// URL is the data parsed from a URL specification for an RPC connection
type URL struct {
	Username string
	Password string
	Protocol string
	Address  string
	Port     uint32
}

// ParseURL takes a string in the format user:pa55word@http://127.0.0.1:11348 and returns a URL struct. This parsing is without any safety checks but should get caught upstream or cause a panic
func ParseURL(url string) (out URL) {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Fprintf(os.Stderr, "ERROR: malformed URL %s\n", r)
	// 		fmt.Println(out)
	// 		out = URL{}
	// 	}
	// }()
	s1 := strings.Split(url, "@")
	var creds []string
	var address string
	switch len(s1) {
	case 1:
		creds = strings.Split(s1[0], ":")
	case 2:
		creds = strings.Split(s1[0], ":")
		address = s1[1]
	default:
		fmt.Fprintln(os.Stderr, "ERROR: malformed URL")
		return URL{}
	}
	out.Username, out.Password = creds[0], creds[1]
	if address == "" {
		out.Protocol, out.Address, out.Port = "http", "127.0.0.1", 11048
	}
	bits := strings.Split(address, "://")
	out.Protocol = "http"
	switch len(bits) {
	case 1:
		adr := strings.Split(bits[0], ":")
		out.Address = adr[0]
		var pr int
		if out.Address == "" {
			out.Address = "127.0.0.1"
		}
		if len(adr) > 1 {
			if adr[1] != "" {
				pr, _ = strconv.Atoi(adr[1])
			}
		}
		if pr > 1024 && pr < 65536 {
			out.Port = uint32(pr)
		} else {
			out.Port = 11048
		}
	case 2:
		out.Protocol = bits[0]
		ip := strings.Split(bits[1], ":")
		out.Address = ip[0]
		out.Address = "127.0.0.1"
		switch len(ip) {
		case 1:
			if out.Address != "" {
				out.Address = ip[0]
			}
			out.Port = 11048
		case 2:
			if ip[0] == "" {
				out.Address = "127.0.0.1"
			}
			port, err := strconv.Atoi(ip[1])
			if err != nil {
				out.Port = 11048
			} else {
				out.Port = uint32(port)
			}
		}
	default:
		fmt.Fprintln(os.Stderr, "ERROR: malformed URL")
		return URL{}
	}

	return
}
