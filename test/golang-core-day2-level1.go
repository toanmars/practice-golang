package main

import "fmt"

type Logger interface {
	Log(msg string)
}

type ConsoleLogger struct {
	message string
}
type FileLogger struct {
	fileName string
}

func (c *ConsoleLogger) Log(msg string) {
	fmt.Println(c.message + msg)
}

type Server struct {
	Logger Logger
}

func (f *FileLogger) Log(msg string) {
	fmt.Println(f.fileName + msg)
}

func main() {
	s1 := Server{Logger: &ConsoleLogger{message: "Console: "}}
	s1.Logger.Log("Server started")

	// Sau này bạn có FileLogger, bạn chỉ cần thay vào mà không sửa Struct Server
	s2 := Server{Logger: &FileLogger{fileName: "file.log"}}
	s2.Logger.Log("Server started")
}
