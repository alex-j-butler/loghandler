package loghandler

import (
	"net"
	"regexp"
	"strconv"
)

const ConnectEventRegex = "\"(.+)<(\\d+)><(.+)><(Blue|Red|Unassigned|Spectator|)>\" connected, address \"(.+):(\\d+)\""

type ConnectEvent struct {
	userEvent
	Address string
	Port    int
}

func connectEventMatcher() *RegexMatcher {
	regex, err := regexp.Compile(ConnectEventRegex)
	if err != nil {
		panic(err)
	}

	return &RegexMatcher{
		Regexp:    regex,
		MatchFunc: handleConnectEvent,
	}
}

func handleConnectEvent(lh *LogHandler, addr *net.UDPAddr, matches []string) {
	port64, _ := strconv.ParseInt(matches[6], 10, 32)
	port := int(port64)

	lh.handle(addr, &ConnectEvent{
		userEvent: userEvent{
			Username: matches[1],
			UserID:   matches[2],
			SteamID:  matches[3],
			Team:     matches[4],
		},
		Address: matches[5],
		Port:    port,
	})
}
