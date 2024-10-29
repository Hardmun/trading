package test

import (
	"strconv"
	"testing"
	"trading/internal/api"
	"trading/test/cases"
)

func TestRequestKlineData(t *testing.T) {
	for _, v := range cases.CasesKline {
		t.Run(strconv.FormatInt(v.TimeStart, 10), func(t *testing.T) {
			err := api.RequestKlineData(v)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func BenchmarkRequestKlineData(b *testing.B) {
	if testing.Short() {
		b.Skip()
	}

	for i := 0; i < b.N; i++ {
		for _, v := range cases.CasesKline {
			err := api.RequestKlineData(v)
			if err != nil {
				b.Error(err)
			}
		}
	}
}
