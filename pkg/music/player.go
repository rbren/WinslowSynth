package music

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/rbren/midi/pkg/config"
	"github.com/rbren/midi/pkg/input"
	"github.com/rbren/midi/pkg/logger"
	"github.com/rbren/midi/pkg/output"
)

const msPerSampleTick = 10
const msPerClearTick = 250

type MusicPlayer struct {
	Sequence       Sequence
	Output         *output.CircularAudioBuffer
	CurrentSample  uint64
	samplesPerTick int
	silence        []float32
}

func NewMusicPlayer(out *output.CircularAudioBuffer) MusicPlayer {
	samplesPerSec := config.MainConfig.SampleRate
	samplesPerMs := samplesPerSec / 1000
	samplesPerTick := samplesPerMs * msPerSampleTick
	logrus.Info("samples per Ms", samplesPerMs)
	logrus.Info("samples per tick", samplesPerTick)
	logrus.Info("output", out.GetCapacity())
	return MusicPlayer{
		Output:         out,
		Sequence:       NewSequence(),
		samplesPerTick: samplesPerTick,
		silence:        make([]float32, samplesPerTick),
	}
}

func (m *MusicPlayer) Clear() {
	m.Sequence = NewSequence()
}

func (m *MusicPlayer) startSampling() {
	defer func() {
		if e := recover(); e != nil {
			logger.Recover("player tick", e)
		}
	}()
	ticker := time.NewTicker(msPerSampleTick * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			logrus.Debug("tick")
			m.nextBytes()
		}
	}
}

func (m *MusicPlayer) startSequencing(notes chan input.InputKey) {
	defer func() {
		if e := recover(); e != nil {
			logger.Recover("while sequencing", e)
		}
	}()
	for {
		select {
		case note := <-notes:
			logrus.Info("note", note)
			m.Sequence.Add(note, m.CurrentSample) // TODO: can we get more accurate than CurrentSample?
		}
	}
}

func (m *MusicPlayer) startClearingSequence() {
	defer func() {
		if e := recover(); e != nil {
			logger.Recover("player tick", e)
		}
	}()
	ticker := time.NewTicker(msPerClearTick * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			m.Sequence.ClearOldEvents(m.CurrentSample)
		}
	}
}

func (m *MusicPlayer) Start(notes chan input.InputKey) {
	go m.startSampling()
	go m.startSequencing(notes)
	go m.startClearingSequence()
}

func (m *MusicPlayer) nextBytes() {
	logrus.Debug("active keys", len(m.Sequence.Events))
	samples := m.Sequence.GetSamples(m.CurrentSample, m.samplesPerTick)
	_, err := m.Output.WriteAudio(samples, samples)
	if err != nil {
		panic(err)
	}
	m.CurrentSample += uint64(m.samplesPerTick)
	logrus.Debug("pos", m.CurrentSample, samples[0], samples[len(samples)-1])
}
