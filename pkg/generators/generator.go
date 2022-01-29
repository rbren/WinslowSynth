package generators

type Generator interface {
	GetValue(elapsed uint64, releasedAt uint64) float32
}
