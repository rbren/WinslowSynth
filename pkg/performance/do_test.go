package performance

import (
	"testing"

	"github.com/rbren/midi/pkg/generators"
)

func TestPerformance(t *testing.T) {
	inst := generators.Library["winslow"]
	CheckPerformance(inst, 100000)
	CheckPerformance(inst, 100000)
	CheckPerformance(inst, 100000)
	CheckPerformance(inst, 100000)
	CheckPerformance(inst, 100000)
	CheckPerformance(inst, 100000)
	CheckPerformance(inst, 100000)
	CheckPerformance(inst, 100000)
}
