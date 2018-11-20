package main

import "time"
import logger ".."

func main() {

	logger.Log("test")
	<-time.After(10 * time.Second)
}
