package audio

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
	"runtime"
	"math"
	"fmt"
)

var (
	current_volume = 1.0
	audio_players []*MixedAudioPlayerBackend
	cam_x, cam_y float64 // listener position, its probably the camera so we just call it that
)

type MixedAudioPlayerBackend struct {
	position_x, position_y float64
	positional bool
	positional_modifier float64
	player *audio.Player
}

type MixedAudioPlayer struct {
	player *MixedAudioPlayerBackend
}

func NewMixedAudioPlayer(bytes []byte) *MixedAudioPlayer {
	new := MixedAudioPlayer{}
	player := audio.CurrentContext().NewPlayerFromBytes(bytes)
	player.SetVolume(current_volume)
	backend := MixedAudioPlayerBackend{}
	backend.player = player
	backend.positional = false
	backend.positional_modifier = 1.0
	new.player = &backend
	audio_players = append(audio_players, &backend)

	runtime.SetFinalizer(&new, func(f *MixedAudioPlayer) {
		fmt.Println("Cleaning up audio player ", f)
		for i, p := range audio_players {
			if p == new.player {
				audio_players[i] = audio_players[len(audio_players)-1]
				audio_players[len(audio_players)-1] = nil
				audio_players = audio_players[:len(audio_players)-1]
				break
			}
		}
	})

	return &new
}

func SetVolume(volume float64) {
	current_volume = volume
	for _, player := range audio_players {
		player.SetVolume(current_volume)
	}
}

func SetCameraPos(x, y float64) {
	cam_x = x
	cam_y = y
	for _, player := range audio_players {
		if player.positional {
			player.positional_modifier = positionalModifier(x, y)
			player.SetVolume(current_volume)
		}
	}
}

func (player MixedAudioPlayerBackend) SetVolume(volume float64) {
	player.player.SetVolume(current_volume * player.positional_modifier)
}

func PlayBytes(bytes []byte) {
	player := audio.CurrentContext().NewPlayerFromBytes(bytes)
	player.SetVolume(current_volume)
	player.Play()
}

func PlayBytesAtPosition(bytes []byte, x, y float64) {
	player := audio.CurrentContext().NewPlayerFromBytes(bytes)
	modifier := positionalModifier(x, y)
	player.SetVolume(current_volume * modifier)
	player.Play()
}

func (player MixedAudioPlayerBackend) Play() {
	player.player.Play()
}

func (player MixedAudioPlayerBackend) Rewind() {
	player.player.Rewind()
}

func (player MixedAudioPlayerBackend) SetPosition(x, y float64) {
	player.positional = true
	player.position_x = x
	player.position_y = y
	player.positional_modifier = positionalModifier(x, y)
	player.player.SetVolume(current_volume * player.positional_modifier)
}

func (player MixedAudioPlayer) Play() {
	player.player.Play()
}

func (player MixedAudioPlayer) Rewind() {
	player.player.Rewind()
}

func (player MixedAudioPlayer) SetPosition(x, y float64) {
	player.player.SetPosition(x, y)
}

func squaredLength(x, y float64) (float64) {
	return (x*x)+(y*y)
}

func positionalModifier(x, y float64) float64 { // TODO: configurable audio distance and such
	modifier := (1000000.0 - squaredLength(x - cam_x, y - cam_y)) / 1000000.0
	return math.Max(math.Min(modifier, 1.0), 0.0)
}
