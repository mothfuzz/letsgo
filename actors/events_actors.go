package actors

import (
	//"fmt"
	"reflect"
)

//actor-specific message passing built on top of pubsub

var actorChannels = map[reflect.Type]map[ActorId]Channel{}

func ensure(t reflect.Type) {
	if _, ok := actorChannels[t]; !ok {
		actorChannels[t] = map[ActorId]Channel{}
	}
}

func Listen(a Actor, e interface{}) Channel {
	if id, ok := GetId(a); ok {
		t := reflect.TypeOf(e)
		c := Subscribe(e)
		ensure(t)
		actorChannels[t][id] = c
		return c
	}
	return nil
}

func Send(a Actor, e interface{}) {
	if id, ok := GetId(a); ok {
		t := reflect.TypeOf(e)
		ensure(t)
		go func(c chan interface{}) { c <- e }(actorChannels[t][id])
	}
}
func SendAll(e interface{}) {
	Publish(e)
}

func DestroyListener(id ActorId) {
	for t := range actorChannels {
		if c, ok := actorChannels[t][id]; ok {
			Unsubscribe(c)
			delete(actorChannels[t], id)
		}
	}
}

func AllListeners(e interface{}, f func(Actor)) {
	t := reflect.TypeOf(e)
	ensure(t)
	for id := range actorChannels[t] {
		if a, ok := GetActor(id); ok {
			f(a)
		}
	}
}
