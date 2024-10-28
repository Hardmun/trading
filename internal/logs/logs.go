package logs

import (
	"fmt"
	"github.com/Hardmun/trading.git/internal/utils"
	"log"
	"os"
	"path/filepath"
)

type LogStruct struct {
	logType string
	logFile *os.File
}

func (l *LogStruct) Write(errMsg ...any) {
	newLine := log.New(l.logFile, fmt.Sprintf("[%s]", l.logType), log.LstdFlags)
	newLine.Println(errMsg...)
}

func (l *LogStruct) Close() {
	err := l.logFile.Close()
	if err != nil {
		l.Write(err)
	}
}

func NewLog(logType string) (*LogStruct, error) {
	logDir, err := utils.DirPath("logs")
	if err != nil {
		return nil, err
	}

	filename := "error.log"

	var lFile *os.File
	lFile, err = os.OpenFile(filepath.Join(logDir, filename), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	return &LogStruct{
		logType: logType,
		logFile: lFile,
	}, nil
}
