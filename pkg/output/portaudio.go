package output

import (
	"github.com/gordonklaus/portaudio"
)

const sampleRate = 44100

type PortAudioOutput struct {
	stream *portaudio.Stream
	Buffer *CircularAudioBuffer
}

func NewPortAudioOutput(sampleRate int) (*PortAudioOutput, error) {
	p := PortAudioOutput{}
	p.Buffer = NewCircularAudioBuffer(sampleRate) // 1 second of samples
	portaudio.Initialize()
	var err error
	p.stream, err = portaudio.OpenDefaultStream(0, 2, float64(sampleRate), 0, p.Buffer.ReadAudio)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (p PortAudioOutput) Start() error {
	return p.stream.Start()
}

func (p PortAudioOutput) Close() error {
	if p.stream != nil {
		err := p.stream.Stop()
		if err != nil {
			return err
		}
	}
	portaudio.Terminate()
	return nil
}
