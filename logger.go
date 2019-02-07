package logger

import (
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var fileName string
var mux sync.Mutex
var initialized bool
var printConsole bool
var fpLog *os.File
var err error
var isOpen bool
var lastLogDate string

// SetFileName will set the log file name
// and re initialize the logger
func SetFileName(newFileName string) {
	fileName = newFileName
	initialized = false // in order to stop pollingData()
	initialized = initialize()
}

func getNow() string {
	return time.Now().Local().Format("2006-01-02")
}

// Log will write log string both on console and in file
func Log(v ...string) {

	logDate := getNow()

	// If last log date is yesterday, then reset to today
	if lastLogDate != "" && strings.Compare(lastLogDate, logDate) == -1 {

		fileName = ""
		initialized = false
	}

	if fileName == "" {
		fileName = getNow() + "_log.txt"
		lastLogDate = getNow()
	}

	if initialized == false {
		initialized = initialize()
	}

	functionName := getCallingFunctionName()

	logStr := functionName + "() --> " + strings.Join(v[:], ", ")
	log.Println(logStr)
}

// Debug will log parameters with Debug string in the beginning
func Debug(v ...string) {

	v = append([]string{"***Debug***"}, v...)

	logDate := getNow()

	// If last log date is yesterday, then reset to today
	if lastLogDate != "" && strings.Compare(lastLogDate, logDate) == -1 {

		fileName = ""
		initialized = false
	}

	if fileName == "" {
		fileName = getNow() + "_log.txt"
		lastLogDate = getNow()
	}

	if initialized == false {
		initialized = initialize()
	}

	functionName := getCallingFunctionName()

	logStr := functionName + "() --> " + strings.Join(v[:], ", ")
	log.Println(logStr)
}

// Error will log parameters with Debug string in the beginning
func Error(v ...string) {

	v = append([]string{"***Error ***"}, v...)

	logDate := getNow()

	// If last log date is yesterday, then reset to today
	if lastLogDate != "" && strings.Compare(lastLogDate, logDate) == -1 {

		fileName = ""
		initialized = false
	}

	if fileName == "" {
		fileName = getNow() + "_log.txt"
		lastLogDate = getNow()
	}

	if initialized == false {
		initialized = initialize()
	}

	functionName := getCallingFunctionName()

	logStr := functionName + "() --> " + strings.Join(v[:], ", ")
	log.Println(logStr)
}

func getCallingFunctionName() string {

	fpcs := make([]uintptr, 1)
	runtime.Callers(3, fpcs)

	// get the info of the actual function that's in the pointer
	return runtime.FuncForPC(fpcs[0] - 1).Name()
}

func openFile() {

	fpLog, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {

		panic(err)
	}

	isOpen = true
	writer := io.MultiWriter(fpLog, os.Stdout)
	log.SetOutput(writer)
}

func closeFile() {

	if fpLog != nil {

		fpLog.Close()
		isOpen = false
	}
}

func initialize() bool {

	if isOpen == true {
		closeFile()
	}

	openFile()

	return true
}

