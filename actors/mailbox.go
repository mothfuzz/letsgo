package actors

import (
	//"fmt"
	"reflect"
)

//actors message passing

type Mailbox chan interface{}

func (m Mailbox) Read() (interface{}, bool) {
	select {
	case i := <-m:
		return i, true
	default:
		return nil, false
	}
}

func (m Mailbox) HandleMessages(f func(interface{})) {
	for {
		if message, ok := m.Read(); ok {
			f(message)
		} else {
			break
		}
	}
}

var mailboxes = map[ActorId]Mailbox{}
var messageTypes = map[reflect.Type][]ActorId{}

func Listen(a Actor, es ...interface{}) Mailbox {
	if id, ok := GetId(a); ok {
		c := make(chan interface{})
		mailboxes[id] = c
		for _, e := range es {
			t := reflect.TypeOf(e)
			messageTypes[t] = append(messageTypes[t], id)
		}
		return c
	}
	return nil
}

func Send(a Actor, e interface{}) {
	if id, ok := GetId(a); ok {
		go func(c chan interface{}) { c <- e }(mailboxes[id])
	}
}
func SendAll(e interface{}) {
	t := reflect.TypeOf(e)
	for _, id := range messageTypes[t] {
		go func(c chan interface{}) { c <- e }(mailboxes[id])
	}
}

func DestroyListener(id ActorId) {
	for m := range mailboxes {
		if m == id {
			delete(mailboxes, m)
			for t, ms := range messageTypes {
				for i := range ms {
					if ms[i] == m {
						ms[i] = ms[len(ms)-1]
						messageTypes[t] = ms[:len(ms)-1]
					}
				}
			}
			return
		}
	}
}

func AllListeners(e interface{}, f func(Actor)) {
	t := reflect.TypeOf(e)
	for _, id := range messageTypes[t] {
		if a, ok := GetActor(id); ok {
			f(a)
		}
	}
}
