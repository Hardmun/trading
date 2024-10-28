package test

import (
	"github.com/Hardmun/trading.git/internal/utils"
	"testing"
)

func TestAdd(t *testing.T) {
	t.Run("testing", func(t *testing.T) {
		res := utils.AddTest(1, 4)
		if res != 5 {
			t.Error("Error")
		}
	})
}
