package main

import (
	"fmt"

	"github.com/mothfuzz/letsgo/actors"
	"github.com/mothfuzz/letsgo/app"
)

type MyMessage struct{}
type MyMessage2 struct{}
type MyMessage3 struct{}

type Receiver struct {
	actors.Mailbox
}

func (r *Receiver) Init() {
	//listen on a specific set of message types, or just actors.Listen(r) for whatever
	r.Mailbox = actors.Listen(r, MyMessage{}, MyMessage2{}, MyMessage3{})
}

func (r *Receiver) MessageHandler(message interface{}) {
	switch m := message.(type) {
	case MyMessage:
		fmt.Println("woohoo!")
		fmt.Println(m)
	case MyMessage2:
		fmt.Println("woohoo!!")
		fmt.Println(m)
	case MyMessage3:
		fmt.Println("woohoo!!!")
		fmt.Println(m)
	default:
		fmt.Println("do not respond")
	}
}

func (r *Receiver) Update() {
	//lots of ways to read messages!
	//you can use a traditional select, the Read() function, or a handler function
	/*for {
		select {
		case m := <-g.Mailbox:
			g.MessageHandler(m)
		default:
			return
		}
	}*/
	/*for {
		if m, ok := g.Mailbox.Read(); ok {
			g.MessageHandler(m)
		} else {
			return
		}
	}*/
	r.Mailbox.HandleMessages(r.MessageHandler)
}

func (r *Receiver) Destroy() {}

func main() {
	app.Init()
	defer app.Quit()

	for i := 0; i < 5; i++ {
		actors.Spawn(&Receiver{})
	}
	//can send to individual actors or all of them
	//note: SendAll(T) will not send to actors who don't listen to T in the first place.
	//Send(a, T), however, will.
	actors.SendAll(MyMessage{})
	actors.SendAll(MyMessage2{})
	actors.SendAll(MyMessage3{})

	//can poll through all listeners of a particular type (good for mxn scanning)
	actors.AllListeners(MyMessage{}, func(a actors.Actor) {
		actors.Send(a, MyMessage{})
	})

	for app.PollEvents() {
		app.Update()
		app.Draw()
	}
}
