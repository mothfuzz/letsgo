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
	/*for {
		select {
		case m := <-g.Mailbox:
			g.MessageHandler(m)
		default:
			break
		}
	}*/
	/*for {
		if m, ok := g.Mailbox.Read(); ok {
			g.MessageHandler(m)
		} else {
			break
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
	actors.SendAll(MyMessage{})
	actors.SendAll(MyMessage2{})
	actors.SendAll(MyMessage3{})

	for app.PollEvents() {
		app.Update()
		app.Draw()
	}
}
