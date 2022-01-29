package music

import (
	"fmt"
	"time"

	"github.com/rbren/midi/pkg/input"
	"github.com/rbren/midi/pkg/logger"
	"github.com/rbren/midi/pkg/output"
)

const msPerTick = 100

type Note struct {
	Frequency float32
	Velocity  int64
}

type MusicPlayer struct {
	SampleRate     int
	ActiveKeys     map[int64]Note
	Output         *output.CircularAudioBuffer
	samplesPerTick int
	silence        []float32
	sampleData     []float32
}

func NewMusicPlayer(sampleRate int, out *output.CircularAudioBuffer) MusicPlayer {
	samplesPerSec := sampleRate
	samplesPerMs := samplesPerSec / 1000
	samplesPerTick := samplesPerMs * msPerTick
	fmt.Println("samples per Ms", samplesPerMs)
	fmt.Println("samples per tick", samplesPerTick)
	fmt.Println("output", out.GetCapacity())
	return MusicPlayer{
		SampleRate:     sampleRate,
		Output:         out,
		ActiveKeys:     map[int64]Note{},
		samplesPerTick: samplesPerTick,
		silence:        make([]float32, samplesPerTick),
		sampleData:     GenerateFrequency(440.0, sampleRate, samplesPerTick),
	}
}

func (m MusicPlayer) Start(notes chan input.InputKey) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("PANIC!!! %v \n", r)
			}
		}()
		ticker := time.NewTicker(msPerTick * time.Millisecond)
		for {
			select {
			case <-ticker.C:
				fmt.Println("tick")
				m.nextBytes()
			}
		}
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("PANIC!!! %v \n", r)
			}
		}()
		for {
			select {
			case note := <-notes:
				fmt.Println("note", note)
				if note.Action == "channel.NoteOn" {
					m.ActiveKeys[note.Key] = Note{
						Frequency: 440.0,
						Velocity:  note.Velocity,
					}
				} else if note.Action == "channel.NoteOff" {
					delete(m.ActiveKeys, note.Key)
				} else {
					fmt.Println("No action for " + note.Action)
				}
			}
		}
	}()
}

func (m MusicPlayer) nextBytes() {
	logger.Log("active keys", len(m.ActiveKeys))

	samples := m.silence
	fmt.Println("check keys")
	for _, _ = range m.ActiveKeys {
		// TODO: don't just take the last one
		//samples = GenerateFrequency(key.Frequency, m.SampleRate, m.samplesPerTick)
		samples = m.sampleData
		fmt.Println("  send music!")
	}
	n, err := m.Output.WriteAudio(samples, samples)
	if err != nil {
		panic(err)
	}
	logger.Log(fmt.Sprintf("  wrote %d of %d", n, len(samples)*4))
	fmt.Println("done bytes")
}
