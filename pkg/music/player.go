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
	// byteSeqs := make([][]byte, len(m.ActiveKeys))
	silence := make([]byte, m.samplesPerTick)
	byteSeqs := [][]byte{}
	logger.Log("active keys", len(m.ActiveKeys))
	logger.Log("  delay", m.Output.GetBufferDelay())
	idx := 0
	for _, key := range m.ActiveKeys {
		byteSeqs = append(byteSeqs, GenerateFrequency(key.Frequency, m.SampleRate, m.samplesPerTick))
		idx++
	}
	toWrite := silence
	if len(byteSeqs) > 0 {
		toWrite = byteSeqs[0]
		logger.Log("  music")
	} else {
		logger.Log("  silence")
	}
	n, err := m.Output.Write(toWrite)
	if err != nil {
		panic(err)
	}
	logger.Log(fmt.Sprintf("  wrote %d of %d", n, len(toWrite)))
	logger.Log("  delay", *m.Output.ReadPos, *m.Output.WritePos)
}
