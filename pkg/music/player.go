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
	silence        []float64
	sampleData     []float64
}

func NewMusicPlayer(sampleRate int, out output.AudioReaderWriter) MusicPlayer {
	samplesPerSec := sampleRate
	samplesPerMs := samplesPerSec / 1000
	samplesPerTick := samplesPerMs * msPerTick
	return MusicPlayer{
		SampleRate:     sampleRate,
		Output:         out,
		ActiveKeys:     map[int64]Note{},
		samplesPerTick: samplesPerTick,
		silence:        make([]float64, samplesPerTick),
		sampleData:     GenerateFrequency(440.0, sampleRate, samplesPerTick),
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

	samples := m.silence
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
	logger.Log("  delay", m.Output.GetBufferDelay())
}
