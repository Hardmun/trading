package logs

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type Logger interface {
	write(errMsg ...any)
}

type logStruct struct {
	logType string
	logFile *os.File
	mtx     *sync.RWMutex
}

func (l *logStruct) write(errMsg ...any) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	newLine := log.New(l.logFile, fmt.Sprintf("[%s]", l.logType), log.LstdFlags)
	newLine.Println(errMsg...)
}

//
//func GetLog() {
//	lg := &LogStruct{}
//	lg.logType = logType
//	if err := lg.initialize(filename); err != nil {
//		log.Fatal(err)
//	}
//
//	return lg
//}
