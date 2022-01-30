package music

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/rbren/midi/pkg/config"
	"github.com/rbren/midi/pkg/generators"
	"github.com/rbren/midi/pkg/input"
	"github.com/rbren/midi/pkg/logger"
	"github.com/rbren/midi/pkg/output"
)

const msPerTick = 10

type MusicPlayer struct {
	Instrument     generators.Instrument
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
	logrus.Info("samples per Ms", samplesPerMs)
	logrus.Info("samples per tick", samplesPerTick)
	logrus.Info("output", out.GetCapacity())
	return MusicPlayer{
		Output:         out,
		Instrument:     generators.GetDefaultInstrument(),
		Generators:     generators.NewRegistry(),
		samplesPerTick: samplesPerTick,
		silence:        make([]float32, samplesPerTick),
	}
}

func (m *MusicPlayer) Clear() {
	m.Generators = generators.NewRegistry()
}

func (m *MusicPlayer) Start(notes chan input.InputKey) {
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
				logrus.Debug("tick")
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
				logrus.Info("note", note)
				g := generators.SetFrequency(m.Instrument, note.Frequency)
				if note.Action == "channel.NoteOn" {
					m.Generators.Attack(note.Key, m.CurrentSample, g)
				} else if note.Action == "channel.NoteOff" {
					m.Generators.Release(note.Key, m.CurrentSample, g)
				} else {
					logrus.Info("No action for " + note.Action)
				}
			}
		}
	}()
}

func (m *MusicPlayer) nextBytes() {
	logrus.Debug("active keys", len(m.Generators.Events))
	samples := m.Generators.GetSamples(m.CurrentSample, m.samplesPerTick)
	_, err := m.Output.WriteAudio(samples, samples)
	if err != nil {
		panic(err)
	}
	m.CurrentSample += uint64(m.samplesPerTick)
	logrus.Debug("pos", m.CurrentSample, samples[0], samples[len(samples)-1])
}
