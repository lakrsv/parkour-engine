package game

import (
	"fmt"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
	"log/slog"
	"time"
)

const (
	Rate beep.SampleRate = 44100
)

var pickupColorSound *beep.Buffer
var walkSound *beep.Buffer
var goalSound *beep.Buffer
var doorOpenSounds []beep.Buffer

func InitAudio() {
	if err := speaker.Init(Rate, Rate.N(time.Second/10)); err != nil {
		slog.Error("Failed initializing speaker", "error", err)
		return
	}
	pickupColorSound = loadSound("assets/audio/pickup_color.wav")
	walkSound = loadSound("assets/audio/walk.wav")
	goalSound = loadSound("assets/audio/goal.wav")
	doorOpenSounds = loadDoorOpenSounds()
}

func loadDoorOpenSounds() []beep.Buffer {
	var buffers []beep.Buffer
	for i := range 11 {
		sound := loadSound(fmt.Sprintf("assets/audio/door_open_%d.wav", i))
		if sound != nil {
			buffers = append(buffers, *sound)
		}
	}
	return buffers
}

func loadSound(path string) *beep.Buffer {
	f, err := content.Open(path)
	defer f.Close()
	if err != nil {
		slog.Error("Failed opening pickup color audio file", "error", err)
		return nil
	} else {
		streamer, format, err := wav.Decode(f)
		if err != nil {
			slog.Error("Failed decoding mp3", "error", err)
			return nil
		} else {
			defer streamer.Close()
			sound := beep.NewBuffer(format)
			sound.Append(streamer)
			return sound
		}
	}
}

func PlayBackgroundMusic() {
	idx := 0
	for {
		f, err := content.Open(fmt.Sprintf("assets/audio/background_%d.mp3", idx%2))
		if err != nil {
			slog.Error("Failed opening audio", "error", err)
			return
		}
		streamer, format, err := mp3.Decode(f)
		if err != nil {
			slog.Error("Failed decoding mp3", "error", err)
			return
		}
		defer streamer.Close()

		resampled := beep.Resample(4, format.SampleRate, Rate, streamer)
		volume := &effects.Volume{Streamer: resampled, Base: 2, Volume: -2}
		done := make(chan bool)
		speaker.Play(beep.Seq(volume, beep.Callback(func() {
			done <- true
		})))
		<-done
		time.Sleep(5 * time.Second)
		idx += 1
	}
}

func playPickupSound() {
	streamer := pickupColorSound.Streamer(0, pickupColorSound.Len())
	resampled := beep.Resample(4, pickupColorSound.Format().SampleRate, Rate, streamer)
	speaker.Play(resampled)
}

func playWalkSound() {
	streamer := walkSound.Streamer(0, walkSound.Len())
	resampled := beep.Resample(4, walkSound.Format().SampleRate, Rate, streamer)
	volume := &effects.Volume{Streamer: resampled, Base: 2, Volume: -1}
	speaker.Play(volume)
}

func playDoorOpenSound(num int) {
	soundIndex := num % len(doorOpenSounds)
	streamer := doorOpenSounds[soundIndex].Streamer(0, doorOpenSounds[soundIndex].Len())
	resampled := beep.Resample(4, doorOpenSounds[soundIndex].Format().SampleRate, Rate, streamer)
	volume := &effects.Volume{Streamer: resampled, Base: 2, Volume: -0.5}
	speaker.Play(volume)
}

func playGoalSound() {
	streamer := goalSound.Streamer(0, goalSound.Len())
	resampled := beep.Resample(4, goalSound.Format().SampleRate, Rate, streamer)
	volume := &effects.Volume{Streamer: resampled, Base: 2, Volume: -0.5}
	speaker.Play(volume)
}
