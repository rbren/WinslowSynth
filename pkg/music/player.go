package music

import (
	"time"

	"github.com/rbren/midi/pkg/config"
	"github.com/rbren/midi/pkg/generators"
	"github.com/rbren/midi/pkg/input"
	"github.com/rbren/midi/pkg/logger"
	"github.com/rbren/midi/pkg/output"
)

const msPerTick = 10

type MusicPlayer struct {
	Generators     generators.Registry
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
		Generators:     generators.NewRegistry(),
		samplesPerTick: samplesPerTick,
		silence:        make([]float32, samplesPerTick),
	}
}

func (m MusicPlayer) Start(notes chan input.InputKey) {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				logger.Recover("player tick", e)
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
			if e := recover(); e != nil {
				logger.Recover("player notes", e)
			}
		}()
		for {
			select {
			case note := <-notes:
				logger.Log("note", note)
				if note.Action == "channel.NoteOn" {
					m.Generators.Attack(note, m.CurrentSample)
				} else if note.Action == "channel.NoteOff" {
					m.Generators.Release(note, m.CurrentSample)
				} else {
					logger.Log("No action for " + note.Action)
				}
			}
		}
	}()
}

func (m *MusicPlayer) nextBytes() {
	logger.Log("active keys", len(m.Generators.Events))
	samples := m.Generators.GetSamples(m.CurrentSample, m.samplesPerTick)
	_, err := m.Output.WriteAudio(samples, samples)
	if err != nil {
		panic(err)
	}
	m.CurrentSample += uint64(m.samplesPerTick)
	logger.Log("pos", m.CurrentSample, samples[0], samples[len(samples)-1])
}
