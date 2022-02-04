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
	Sequence       Sequence
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
		Sequence:       NewSequence(),
		samplesPerTick: samplesPerTick,
		silence:        make([]float32, samplesPerTick),
	}
}

func (m *MusicPlayer) Clear() {
	m.Sequence = NewSequence()
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
				m.Sequence.Add(note, m.CurrentSample) // TODO: can we get more accurate than CurrentSample?
			}
		}
	}()
}

func (m *MusicPlayer) nextBytes() {
	logrus.Debug("active keys", len(m.Sequence.Events))
	samples := m.Sequence.GetSamples(m.Instrument, m.CurrentSample, m.samplesPerTick)
	_, err := m.Output.WriteAudio(samples, samples)
	if err != nil {
		panic(err)
	}
	m.CurrentSample += uint64(m.samplesPerTick)
	logrus.Debug("pos", m.CurrentSample, samples[0], samples[len(samples)-1])
}
