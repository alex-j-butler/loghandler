package loghandler

import (
	"net"
	"regexp"
)

// RegexFunc is the handler function type for a RegexMatcher.
type RegexFunc func(*LogHandler, *net.UDPAddr, []string)

// RegexMatcher is a implementation of the Matcher interface
// that matches the data by compiled regex.
type RegexMatcher struct {
	// Regex object to match against.
	Regexp *regexp.Regexp
	// Function to handle the match.
	MatchFunc RegexFunc
	// Cache the previous match.
	previousMatches []string
}

// Match checks whether the correct input data has been received for a match.
func (r *RegexMatcher) Match(data string) bool {
	matches := r.Regexp.FindStringSubmatch(data)
	r.previousMatches = matches
	if len(matches) > 0 {
		return true
	}
	return false
}

// Handle calls the RegexMatcher's MatchFunc when a match was found.
func (r *RegexMatcher) Handle(lh *LogHandler, addr *net.UDPAddr) {
	r.MatchFunc(lh, addr, r.previousMatches)
}
