package actors

//"fmt"

type ActorId struct {
	id int64
}
type Actor interface {
	Update()
}
type OnInit interface {
	Init()
}
type OnDestroy interface {
	Destroy()
}

var base_id = int64(1)
var actorsMap = make(map[ActorId]Actor)

func Spawn(a Actor) ActorId {
	id := ActorId{base_id}
	base_id++
	actorsMap[id] = a
	if a, ok := a.(OnInit); ok {
		a.Init()
	}
	return id
}

func Update() {
	for _, a := range actorsMap {
		a.Update()
	}
}

//run systems across all actors.
func All(f func(Actor)) {
	for _, a := range actorsMap {
		f(a)
	}
}

func Destroy(a Actor) {
	if id, ok := GetId(a); ok {
		DestroyId(id)
	}
}
func DestroyId(id ActorId) {
	if id.id != 0 {
		DestroyListener(id)
		if a, ok := actorsMap[id].(OnDestroy); ok {
			a.Destroy()
		}
		delete(actorsMap, id)
	}
}

func DestroyAll() {
	for id := range actorsMap {
		DestroyId(id)
	}
}

func GetId(a Actor) (ActorId, bool) {
	for id, aa := range actorsMap {
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
