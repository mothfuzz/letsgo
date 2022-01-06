package events

import (
	//"fmt"
	"reflect"
)

type Channel = chan interface{}

var channels = make(map[reflect.Type][]Channel)

func Subscribe(typeVal interface{}) Channel {
	c := make(chan interface{})
	t := reflect.TypeOf(typeVal)
	channels[t] = append(channels[t], c)
	return c
}

func Publish(value interface{}) {
	if v, ok := channels[reflect.TypeOf(value)]; ok {
		for _, c := range v {
			go func(c chan interface{}) { c <- value }(c)
		}
	}
}

func Unsubscribe(ch Channel) {
	for t, v := range channels {
		for i, c := range v {
			if c == ch {
				close(ch)
				l := len(channels[t]) - 1
				//fmt.Printf("channels len: %d\n", l)
				//replace element with last element
				channels[t][i] = channels[t][l]
				//shorten list by 1
				channels[t] = channels[t][:l]
				return
			}
		}
	}
}
