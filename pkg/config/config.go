package config

type Config struct {
	SampleRate int
}

var MainConfig = Config{
	SampleRate: 48000,
}
