package test

import (
	"github.com/Hardmun/trading.git/internal/logs"
	"testing"
)

func TestNewLog(t *testing.T) {
	l, err := logs.NewLog("ERROR")
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 10000000; i++ {
		go l.Write(i)
		//t.Run(fmt.Sprintf("logg #%v", i), func(t *testing.T) {
		//	l.Write(i)
		//})
	}
}

//func BenchmarkAdd(b *testing.B) {
//	if testing.Short() {
//		b.Skip()
//	}
//	for i := 0; i > b.N; i++ {
//		utils.Add(1, 4)
//	}
//}
