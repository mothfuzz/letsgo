package actors

import "github.com/mothfuzz/dyndraw/framework/transform"

func SpawnAt(a Actor, at transform.Transform) ActorId {
	if t, ok := a.(transform.HasTransform); ok {
		*t.GetTransform() = at
	}
	return Spawn(a)
}
