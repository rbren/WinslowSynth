package output

import (
	"time"

	oto "github.com/hajimehoshi/oto/v2"

	"github.com/rbren/midi/pkg/logger"
)

type OutputLine struct {
	sampleRate int
	Line       *CircularAudioBuffer
	Player     oto.Player
}

func NewOutputLine(sampleRate int) (*OutputLine, error) {
	line := NewCircularAudioBuffer(sampleRate * 4 * 10)
	logger.Log("create output", sampleRate)
	ctx, _, err := oto.NewContext(sampleRate, 2, 2)
	if err != nil {
		return nil, err
	}
	player := ctx.NewPlayer(line)

	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				logger.Log("Unplayed:", player.UnplayedBufferSize())
				if err := player.Err(); err != nil {
					logger.Log("player had an error:", err)
				}
			}
		}
	}()

	return &OutputLine{
		sampleRate: sampleRate,
		Line:       line,
		Player:     player,
	}, nil
}
