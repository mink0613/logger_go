package logger

import "time"

func main() {

	Log("test")
	<-time.After(10 * time.Second)

}
