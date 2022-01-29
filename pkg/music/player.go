package music

import (
	"fmt"
	"time"

	"github.com/rbren/midi/pkg/buffers"
	"github.com/rbren/midi/pkg/config"
	"github.com/rbren/midi/pkg/generators"
	"github.com/rbren/midi/pkg/input"
	"github.com/rbren/midi/pkg/logger"
	"github.com/rbren/midi/pkg/output"
)

const msPerTick = 10

type MusicPlayer struct {
	ActiveKeys     map[int64]input.InputKey
	Output         *output.CircularAudioBuffer
	CurrentSample  uint64
	samplesPerTick int
	silence        []float32
}

func NewMusicPlayer(out *output.CircularAudioBuffer) MusicPlayer {
	samplesPerSec := config.MainConfig.SampleRate
	samplesPerMs := samplesPerSec / 1000
	samplesPerTick := samplesPerMs * msPerTick
	logger.Log("samples per Ms", samplesPerMs)
	logger.Log("samples per tick", samplesPerTick)
	logger.Log("output", out.GetCapacity())
	return MusicPlayer{
		Output:         out,
		ActiveKeys:     map[int64]input.InputKey{},
		samplesPerTick: samplesPerTick,
		silence:        make([]float32, samplesPerTick),
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
				logger.Log("tick")
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
				logger.Log("note", note)
				if note.Action == "channel.NoteOn" {
					m.ActiveKeys[note.Key] = note
				} else if note.Action == "channel.NoteOff" {
					delete(m.ActiveKeys, note.Key)
				} else {
					logger.Log("No action for " + note.Action)
				}
			}
		}
	}()
}

func (m *MusicPlayer) nextBytes() {
	logger.Log("active keys", len(m.ActiveKeys))

	samples := m.silence
	for _, key := range m.ActiveKeys {
		keySamples := generators.GenerateSine(key.Frequency, m.samplesPerTick, m.CurrentSample)
		samples = buffers.MixBuffers([][]float32{samples, keySamples})
	}
	_, err := m.Output.WriteAudio(samples, samples)
	if err != nil {
		panic(err)
	}
	m.CurrentSample += uint64(m.samplesPerTick)
	logger.Log("pos", m.CurrentSample, samples[0], samples[len(samples)-1])
}
