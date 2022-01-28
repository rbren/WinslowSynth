package output

import (
	"math"
	"time"

	"github.com/gordonklaus/portaudio"
)

const sampleRate = 44100

type PortAudioOutput struct {
	stream *portaudio.Stream
	Buffer *CircularAudioBuffer
}

func (p PortAudioOutput) Start(sampleRate int) error {
	p.Buffer = NewCircularAudioBuffer(sampleRate) // 1 second of samples
	portaudio.Initialize()
	var err error
	p.stream, err = portaudio.OpenDefaultStream(0, 2, float64(sampleRate), 0, p.Buffer.ReadAudio)
	if err != nil {
		return err
	}
	return p.stream.Start()
}

func (p PortAudioOutput) Close() error {
	err := p.stream.Stop()
	if err != nil {
		return err
	}
	portaudio.Terminate()
	return nil
}

// Stuff below here gets deleted

func PlaySine() {
	s := NewStereoSine(256, 320, sampleRate)
	defer s.Close()
	chk(s.Start())
	time.Sleep(2 * time.Second)
	chk(s.Stop())
}

type stereoSine struct {
	*portaudio.Stream
	stepL, phaseL float64
	stepR, phaseR float64
}

func NewStereoSine(freqL, freqR, sampleRate float64) *stereoSine {
	s := &stereoSine{nil, freqL / sampleRate, 0, freqR / sampleRate, 0}
	var err error
	chk(err)
	return s
}

func (g *stereoSine) processAudio(out [][]float32) {
	for i := range out[0] {
		out[0][i] = float32(math.Sin(2 * math.Pi * g.phaseL))
		_, g.phaseL = math.Modf(g.phaseL + g.stepL)
		out[1][i] = float32(math.Sin(2 * math.Pi * g.phaseR))
		_, g.phaseR = math.Modf(g.phaseR + g.stepR)
	}
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
