package sound

import (
	"github.com/mothfuzz/letsgo/resources"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

//sound is admittedly not that complex

var sounds = map[string]*mix.Chunk{}
var musics = map[string]*mix.Music{}

func PlaySound(name string) {
	if _, ok := sounds[name]; !ok {
		//so many ways this can go wrong
		bytes, err := resources.ReadResource("sounds/" + name)
		if err != nil {
			panic(err)
		}
		rw, err := sdl.RWFromMem(bytes)
		if err != nil {
			panic(err)
		}
		wav, err := mix.LoadWAVRW(rw, false)
		if err != nil {
			panic(err)
		}
		sounds[name] = wav
	}
	sounds[name].Play(-1, 0)
}

func PlayMusic(name string) {
	if _, ok := musics[name]; !ok {
		bytes, err := resources.ReadResource("sounds/" + name)
		if err != nil {
			panic(err)
		}
		rw, err := sdl.RWFromMem(bytes)
		if err != nil {
			panic(err)
		}
		mus, err := mix.LoadMUSRW(rw, 0)
		if err != nil {
			panic(err)
		}
		musics[name] = mus
	}
	musics[name].Play(-1)

}
