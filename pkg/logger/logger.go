package logger

import (
	"log"
	"os"
)

var InfoLog *log.Logger
var ErrorfoLog *log.Logger

func SetGlobalLogger() {
	InfoLog = log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
	ErrorfoLog = log.New(os.Stdout, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)
}
