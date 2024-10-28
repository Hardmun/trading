package test

import (
	"github.com/Hardmun/trading.git/internal/utils"
	"testing"
)

func TestAdd(t *testing.T) {
	t.Run("testing", func(t *testing.T) {
		res := utils.Add(1, 4)
		if res != 5 {
			t.Error("Error")
		}
	})
}

func BenchmarkAdd(b *testing.B) {
	if testing.Short() {
		b.Skip()
	}
	for i := 0; i > b.N; i++ {
		utils.Add(1, 4)
	}
}
