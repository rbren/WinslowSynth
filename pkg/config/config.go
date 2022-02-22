package config

type Config struct {
	SampleRate    int
	HistoryMs     int
	FrequencyBins int

	MaxSampleRateHandicap  float32
	SampleRateHandicapJump float32
	SampleRateRatioMin     float32
	SampleRateRatioMax     float32

	InstrumentSampleMs      int
	InstrumentClearEventsMs int

	MaxReleaseTimeMs    int
	MinZeroedTimeMs     int
	ZeroSampleThreshold float32

	ServerSendIntervalMs int
}

var MainConfig = Config{
	SampleRate:    48000,
	HistoryMs:     5000,
	FrequencyBins: 500,

	MaxSampleRateHandicap:  0.9,
	SampleRateHandicapJump: .01,
	SampleRateRatioMin:     .3,
	SampleRateRatioMax:     .5,

	InstrumentSampleMs:      20,
	InstrumentClearEventsMs: 250,

	MaxReleaseTimeMs:    10000,
	MinZeroedTimeMs:     1000,
	ZeroSampleThreshold: 0.01,

	ServerSendIntervalMs: 50,
}
