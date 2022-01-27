package music

import (
	"fmt"
	"time"

	"github.com/rbren/midi/pkg/output"
)

type Note struct {
	Frequency float64
	Velocity  int64
}

type MusicPlayer struct {
	SampleRate int
	ActiveKeys map[int64]Note
	Output     output.MusicReaderWriter
}

func (m MusicPlayer) Start() {
	samplesPerMillisecond := m.SampleRate / 1000
	samplesPerGap := 5 * samplesPerMillisecond
	ticker := time.NewTicker(5 * time.Millisecond)
	go func() {
		for {
			select {
			case <-ticker.C:
				m.nextBytes(samplesPerGap)
			}
		}
	}()
}

func (m MusicPlayer) nextBytes(samples int) {
	// byteSeqs := make([][]byte, len(m.ActiveKeys))
	byteSeqs := [][]byte{}
	fmt.Println("active keys", len(m.ActiveKeys))
	idx := 0
	for _, key := range m.ActiveKeys {
		byteSeqs = append(byteSeqs, GenerateFrequency(key.Frequency, m.SampleRate, samples))
		idx++
	}
	if len(byteSeqs) > 0 {
		m.Output.Write(byteSeqs[0])
	} else {
		silence := make([]byte, samples)
		m.Output.Write(silence)
	}
}
