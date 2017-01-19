package loghandler

import (
	"net"
	"regexp"
)

const DisconnectEventRegex = "\"(.+)<(\\d+)><(.+)><(Blue|Red|Unassigned|Spectator|)>\" disconnected \\(reason \"(.+)\"\\)"

type DisconnectEvent struct {
	userEvent
	Reason string
}

func disconnectEventMatcher() *RegexMatcher {
	regex, err := regexp.Compile(DisconnectEventRegex)
	if err != nil {
		panic(err)
	}

	return &RegexMatcher{
		Regexp:    regex,
		MatchFunc: handleDisconnectEvent,
	}
}

func handleDisconnectEvent(lh *LogHandler, addr *net.UDPAddr, matches []string) {
	lh.handle(addr, &DisconnectEvent{
		userEvent: userEvent{
			Username: matches[1],
			UserID:   matches[2],
			SteamID:  matches[3],
			Team:     matches[4],
		},
		Reason: matches[5],
	})
}
