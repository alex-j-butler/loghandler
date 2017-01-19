package loghandler

import "net"

type Matcher interface {
	Match(string) bool
	Handle(*LogHandler, *net.UDPAddr)
}
