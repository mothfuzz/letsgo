package actors

import (
	//"fmt"
)

type ActorId struct {
	id int64
}
type Actor interface {
	Init()
	Update()
	Draw()
	Destroy()
}

var base_id = int64(1)
var actorsMap = make(map[ActorId]Actor)

func Spawn(a Actor) ActorId {
	id := ActorId{base_id}
	base_id++
	actorsMap[id] = a
	a.Init()
	return id
}

func Update() {
	for _, a := range(actorsMap) {
		a.Update()
	}
}
func Draw() {
	for _, a := range(actorsMap) {
		a.Draw()
	}
}

func Destroy(a Actor) {
	if id, ok := GetId(a); ok {
		DestroyId(id)
	}
}
func DestroyId(id ActorId) {
	if id.id != 0 {
		actorsMap[id].Destroy()
		delete(actorsMap, id)
	}
}

func Quit() {
	for _, a := range(actorsMap) {
		a.Destroy()
	}
}

func GetId(a Actor) (ActorId, bool) {
	for id, aa := range(actorsMap) {
		if aa == a {
			return id, true
		}
	}
	return ActorId{0}, false
}

func GetActor(id ActorId) (Actor, bool) {
	a, ok := actorsMap[id]
	return a, ok
}
