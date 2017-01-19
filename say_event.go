package loghandler

import (
	"net"
	"regexp"
)

const SayEventRegex = "\"(.+)<(\\d+)><(.+)><(Blue|Red|Unassigned|Spectator)>\" say \"(.+)\""

// SayEvent represents a 'SayText' event sent by the TF2 logging client.
type SayEvent struct {
	userEvent
	Message string
}

// Internal function to create an instance of the RegexMatcher that is used to match this event.
func sayEventMatcher() *RegexMatcher {
	regex, err := regexp.Compile(SayEventRegex)
	if err != nil {
		panic(err)
	}

	return &RegexMatcher{
		Regexp:    regex,
		MatchFunc: handleSayEvent,
	}
}

// Internal function to handle SayEvent.
func handleSayEvent(lh *LogHandler, addr *net.UDPAddr, matches []string) {
	lh.handle(addr, &SayEvent{
		userEvent: userEvent{
			Username: matches[1],
			UserID:   matches[2],
			SteamID:  matches[3],
			Team:     matches[4],
		},
		Message: matches[5],
	})
}
