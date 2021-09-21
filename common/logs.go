package common

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

func InitLog() {
	if _, err := os.Stat("./logs/"); err != nil {
		os.MkdirAll("./logs/", 0777)
	}
	logName := fmt.Sprintf("./logs/%s_%d.log", time.Now().Format("2006_01_02"), rand.New(rand.NewSource(time.Now().UnixNano())).Intn(999999))
	f, err := os.OpenFile(logName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	writers := []io.Writer{
		f,
		os.Stdout}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	log.SetOutput(fileAndStdoutWriter)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}
