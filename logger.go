package logger

import (
	"container/list"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var fileName string
var logQueue *list.List
var mux sync.Mutex
var initialized bool
var printConsole bool
var fpLog *os.File
var err error
var isOpen bool

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

	if fileName == "" {
		fileName = getNow() + "_log.txt"
	}

	if initialized == false {
		initialized = initialize()
	}

	fpcs := make([]uintptr, 1)
	runtime.Callers(2, fpcs)

	// get the info of the actual function that's in the pointer
	fun := runtime.FuncForPC(fpcs[0] - 1)

	log := fun.Name() + "() --> " + strings.Join(v[:], ", ")
	logQueue.PushBack(log)
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

	logQueue = list.New()
	logQueue.Init()

	go pollingData()

	return true
}

func writeLog() {

	if logQueue.Len() > 0 {
		if isOpen == false {
			openFile()
		}

		logStr := logQueue.Front()
		log.Println(logStr.Value)
		logQueue.Remove(logStr)
	} else {
		if isOpen == true {
			isOpen = false
			closeFile()
		}
	}
}

func pollingData() {

	for {
		if initialized == true {
			mux.Lock()
			go writeLog()
			mux.Unlock()
		}
		<-time.After(10 * time.Millisecond)
	}
}
