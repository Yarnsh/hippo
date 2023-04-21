package audio

import (
	"log"
	"strings"
	"os"
	"io"
	"io/fs"
	"bytes"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

var (
	Sounds map[string][]byte
	FileSystem = os.DirFS(".")
)

func Init() {
	audio.NewContext(44100)
	Sounds = make(map[string][]byte)
}

func LoadSoundPath(path string) { // TODO: a func to clean up the sound map
	if strings.HasSuffix(path, ".wav") {
		// Load it to bytes so we can make multiple streams from the same source
		dat, err := fs.ReadFile(FileSystem, path)
    	if err != nil {
			log.Fatal(err)
		}
		s, err := wav.DecodeWithSampleRate(audio.CurrentContext().SampleRate(), bytes.NewReader(dat))
    	if err != nil {
			log.Fatal(err)
		}
		result_bytes, err := io.ReadAll(s)
    	if err != nil {
			log.Fatal(err)
		}
		Sounds[path] = result_bytes
	} else if strings.HasSuffix(path, ".ogg") {
		// Load it to bytes so we can make multiple streams from the same source
		dat, err := fs.ReadFile(FileSystem, path)
    	if err != nil {
			log.Fatal(err)
		}
		s, err := vorbis.DecodeWithSampleRate(audio.CurrentContext().SampleRate(), bytes.NewReader(dat))
    	if err != nil {
			log.Fatal(err)
		}
		result_bytes, err := io.ReadAll(s)
    	if err != nil {
			log.Fatal(err)
		}
		Sounds[path] = result_bytes
	} else if strings.HasSuffix(path, ".mp3") {
		dat, err := fs.ReadFile(FileSystem, path)
    	if err != nil {
			log.Fatal(err)
		}
		s, err := mp3.DecodeWithSampleRate(audio.CurrentContext().SampleRate(), bytes.NewReader(dat))
    	if err != nil {
			log.Fatal(err)
		}
		result_bytes, err := io.ReadAll(s)
    	if err != nil {
			log.Fatal(err)
		}
		Sounds[path] = result_bytes
	}
}

func PlayAt(sound_path string, x, y float64) {
	sfx, ok := Sounds[sound_path]
	if !ok {
		log.Fatal("Missing Sound " + sound_path) // TODO: load and play from the file directly, making sure to close it and clean up the reference correctly
	}
	PlayBytesAtPosition(sfx, x, y)
}

func Play(sound_path string) {
	sfx, ok := Sounds[sound_path]
	if !ok {
		log.Fatal("Missing Sound " + sound_path) // TODO: load and play from the file directly, making sure to close it and clean up the reference correctly
	}
	PlayBytes(sfx)
}

func GetPlayingPlayer(sound_path string) (*audio.Player) {
	sfx, ok := Sounds[sound_path]
	if !ok {
		log.Fatal("Missing Sound " + sound_path) // TODO: load and play from the file directly, making sure to close it and clean up the reference correctly
	}
	player := audio.CurrentContext().NewPlayerFromBytes(sfx)
	player.SetVolume(current_volume)
	player.Play()
	return player
}
