package actors

import (
	//"fmt"
	"github.com/mothfuzz/dyndraw/framework/events"
	"reflect"
)

type Channel = events.Channel

var channels = map[reflect.Type]map[ActorId]Channel{}

func ensure(t reflect.Type) {
	if _, ok := channels[t]; !ok {
		channels[t] = map[ActorId]Channel{}
	}
}

func Listen(a Actor, e interface{}) events.Channel {
	if id, ok := GetId(a); ok {
		t := reflect.TypeOf(e)
		c := events.Subscribe(e)
		ensure(t)
		channels[t][id] = c
		return c
	}
	return nil
}

func Send(a Actor, e interface{}) {
	if id, ok := GetId(a); ok {
		t := reflect.TypeOf(e)
		ensure(t)
		go func(c chan interface{}) { c <- e }(channels[t][id])
	}
}
func SendAll(e interface{}) {
	events.Publish(e)
}

func DestroyListener(id ActorId) {
	for t := range channels {
		if c, ok := channels[t][id]; ok {
			events.Unsubscribe(c)
			delete(channels[t], id)
		}
	}
}

func AllListeners(e interface{}, f func(Actor)) {
	t := reflect.TypeOf(e)
	ensure(t)
	actors := channels[t]
	for id := range actors {
		if a, ok := GetActor(id); ok {
			f(a)
		}
	}
}
