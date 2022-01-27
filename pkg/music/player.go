package music

import (
	"fmt"
	"time"

	"github.com/rbren/midi/pkg/logger"
	"github.com/rbren/midi/pkg/output"
)

const msPerTick = 100

type Note struct {
	Frequency float64
	Velocity  int64
}

type MusicPlayer struct {
	SampleRate     int
	ActiveKeys     map[int64]Note
	Output         output.AudioReaderWriter
	samplesPerTick int
}

func NewMusicPlayer(sampleRate int, out output.AudioReaderWriter) MusicPlayer {
	return MusicPlayer{
		SampleRate:     sampleRate,
		Output:         out,
		ActiveKeys:     map[int64]Note{},
		samplesPerTick: msPerTick * (sampleRate / 1000),
	}
}

func (m MusicPlayer) Start() {
	ticker := time.NewTicker(msPerTick * time.Millisecond)
	go func() {
		for {
			select {
			case <-ticker.C:
				m.nextBytes()
			}
		}
	}()
}

func (m MusicPlayer) nextBytes() {
	logger.Log("active keys", len(m.ActiveKeys))
	logger.Log("  delay", m.Output.GetBufferDelay())

	samples := make([]float64, m.samplesPerTick) // silence
	for _, key := range m.ActiveKeys {
		// TODO: don't just take the last one
		samples = GenerateFrequency(key.Frequency, m.SampleRate, m.samplesPerTick)
	}
	n, err := m.Output.WriteAudio(samples, samples)
	if err != nil {
		panic(err)
	}
	logger.Log(fmt.Sprintf("  wrote %d of %d", n, len(samples)*4))
	logger.Log("  delay", *m.Output.ReadPos, *m.Output.WritePos)
}
