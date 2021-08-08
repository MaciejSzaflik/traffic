package smallMath

import (
	"math/rand"
	"testing"
)

func TestCantorPairDepair(t *testing.T) {
	for i := 0; i < 100; i++ {
		x := rand.Intn(2000)
		y := rand.Intn(2000)

		pair := CantorPair(x, y)

		xTest, yTest := CantorDepair(pair)

		if x != xTest || y != yTest {
			t.Errorf("Pair: %d x: %d, y:%d diffrent than: x1: %d, y1:%d", pair, x, y, xTest, yTest)
		}
	}
}
