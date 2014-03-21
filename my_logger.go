package testView

import (
	"log"
	"net/http"
	"time"
	"os"
)

// filName should be called: localServerLog
func PrepareServerLogger(fileName string,) (logger *log.Logger, err error) {
	if !LogFileAlreadyExists(fileName) {
		err = CreateLogFile(fileName)
		if err != nil {
			return
		}
	}
	myFile, err = os.OpenFile(name, O_RDWR, 0666)
	if err == nil {
		logger = log.New(myFile, "[pLang] ", 0)
		defer myFile.Close()
	}
	return
}

func WriteLog(content string, logger *log.Logger) {
	logger.Printf("%v\n", content)
}

func LogFileAlreadyExists(name string) bool {
    _, err := os.Stat(name)
    if err != nil {
       if os.IsNotExist(err) {
            return false
        }
    }
    return true
}

func CreateLogFile(name string) error {
    f, err := os.Create(name)
    if err != nil {
        return err
    }
    defer f.Close()
    return nil
}