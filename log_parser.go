package loghandler

import (
	"fmt"
	"log"
	"net"
	"reflect"
	"sync"
)

// LogHandler represents a TF2 log handler server.
// Use the Dial function to create an instance of the server.
type LogHandler struct {
	Address    string
	Port       int
	conn       *net.UDPConn
	handlersMu sync.RWMutex
	handlers   map[interface{}][]reflect.Value
	// regexHandlers map[*regexp.Regexp]RegexFunc
	matchers []Matcher
}

// Dial creates an instance of the TF2 log handler server.
// Address specifies the IP address/hostname to listen on.
// Port represents the UDP port to listen on.
func Dial(address string, port int) (*LogHandler, error) {
	lh := &LogHandler{
		Address: address,
		Port:    port,
	}
	// Adding matches here so that they can be iterated over when a message is received.
	lh.matchers = []Matcher{
		connectEventMatcher(),
		disconnectEventMatcher(),
		sayEventMatcher(),
	}

	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", lh.Address, lh.Port))
	if err != nil {
		return nil, err
	}

	lh.conn, err = net.ListenUDP("udp", serverAddr)
	if err != nil {
		return nil, err
	}

	go lh.handleConn()

	return lh, nil
}

func (lh *LogHandler) handleConn() {
	// 1024 bytes should be enough for any message the TF2 server could send.
	buf := make([]byte, 1024)

	for {
		n, addr, err := lh.conn.ReadFromUDP(buf)
		if err != nil {
			log.Println("LogHandler error:", err)
		}

		data := string(buf[:n])
		log.Println("Read", n, "bytes")

		// Iterate over the matchers, and check whether they match the data just received.
		// If they do, then call the Handle function for that matcher.
		for _, v := range lh.matchers {
			if v.Match(data) {
				v.Handle(lh, addr)
				break
			}
		}
	}
}

// AddHandler adds an event handler to the server.
func (lh *LogHandler) AddHandler(handler interface{}) func() {
	lh.initialise()

	eventType := lh.validateHandler(handler)

	lh.handlersMu.Lock()
	defer lh.handlersMu.Unlock()

	h := reflect.ValueOf(handler)

	lh.handlers[eventType] = append(lh.handlers[eventType], h)

	return func() {
		lh.handlersMu.Lock()
		defer lh.handlersMu.Unlock()

		handlers := lh.handlers[eventType]
		for i, v := range handlers {
			if h == v {
				lh.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
				return
			}
		}
	}
}

func (lh *LogHandler) initialise() {
	lh.handlersMu.Lock()
	if lh.handlers != nil {
		lh.handlersMu.Unlock()
		return
	}

	lh.handlers = make(map[interface{}][]reflect.Value)
	lh.handlersMu.Unlock()
}

func (lh *LogHandler) validateHandler(handler interface{}) reflect.Type {
	handlerType := reflect.TypeOf(handler)

	if handlerType.NumIn() != 2 {
		panic("Unable to add event handler, handler must be of type func(*loghandler.LogHandler, *loghandler.EventType)")
	}

	if handlerType.In(0) != reflect.TypeOf(lh) {
		panic("Unable to add event handler, first argument must be of type *loghandler.LogHandler")
	}

	eventType := handlerType.In(1)

	if eventType.Kind() == reflect.Interface {
		eventType = nil
	}

	return eventType
}

func (lh *LogHandler) handle(addr *net.UDPAddr, event interface{}) {
	lh.handlersMu.RLock()
	defer lh.handlersMu.RUnlock()

	if lh.handlers == nil {
		return
	}

	handlerParameters := []reflect.Value{reflect.ValueOf(lh), reflect.ValueOf(event)}

	if handlers, ok := lh.handlers[nil]; ok {
		for _, handler := range handlers {
			go handler.Call(handlerParameters)
		}
	}

	if handlers, ok := lh.handlers[reflect.TypeOf(event)]; ok {
		for _, handler := range handlers {
			go handler.Call(handlerParameters)
		}
	}
}
