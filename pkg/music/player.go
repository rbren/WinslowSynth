package music

import (
	"fmt"
	"time"

	"github.com/rbren/midi/pkg/output"
)

const msPerTick = 5

type Note struct {
	Frequency float64
	Velocity  int64
}

type MusicPlayer struct {
	SampleRate     int
	ActiveKeys     map[int64]Note
	Output         output.MusicReaderWriter
	samplesPerTick int
}

func NewMusicPlayer(sampleRate int, out output.MusicReaderWriter) MusicPlayer {
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
	byteSeqs := [][]byte{}
	fmt.Println("active keys", len(m.ActiveKeys))
	idx := 0
	for _, key := range m.ActiveKeys {
		byteSeqs = append(byteSeqs, GenerateFrequency(key.Frequency, m.SampleRate, m.samplesPerTick))
		idx++
	}
	if len(byteSeqs) > 0 {
		fmt.Println("  note")
		m.Output.Write(byteSeqs[0])
	} else {
		fmt.Println("  silence")
		silence := make([]byte, m.samplesPerTick)
		m.Output.Write(silence)
	}
}
