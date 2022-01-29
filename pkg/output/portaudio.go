package output

import (
	"github.com/gordonklaus/portaudio"

	"github.com/rbren/midi/pkg/config"
)

type PortAudioOutput struct {
	stream *portaudio.Stream
	Buffer *CircularAudioBuffer
}

func NewPortAudioOutput() (*PortAudioOutput, error) {
	p := PortAudioOutput{}
	p.Buffer = NewCircularAudioBuffer(config.MainConfig.SampleRate) // 1 second of samples
	portaudio.Initialize()
	var err error
	p.stream, err = portaudio.OpenDefaultStream(0, 2, float64(config.MainConfig.SampleRate), 0, p.Buffer.ReadAudio)
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
