package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
)

var (
	pipeCSFile = "/tmp/client-server"
	pipeSCFile = "/tmp/server-client"
)

func main() {
	checkPipeExist(pipeCSFile)
	checkPipeExist(pipeSCFile)

	// prepare
	buffer := make([]byte, 1024)

	pipeCS, err := os.OpenFile(pipeCSFile, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		fmt.Println("Error opening pipe for reading:", err)
		return
	}
	defer pipeCS.Close()

	pipeSC, err := os.OpenFile(pipeSCFile, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		fmt.Println("Error opening pipe for writing:", err)
		return
	}
	defer pipeSC.Close()

	// receive
	for {
		n, err := pipeCS.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				// do nothing
			} else {
				fmt.Println("Error reading from pipe:", err)
				return
			}
		}

		if n > 0 {
			// some processing
			message := string(buffer[:n])
			i, _ := strconv.Atoi(message)
			i++
			
			// send
			message = strconv.Itoa(i)
			_, err = pipeSC.Write([]byte(message))
			if err != nil {
				fmt.Println("Error writing to pipeSC:", err)
			}
		}
	}
}

func checkPipeExist(pipe string) {
	if _, err := os.Stat(pipe); os.IsNotExist(err) {
		if err := syscall.Mkfifo(pipe, 0777); err != nil {
			fmt.Println("Error creating named pipe:", err)
			return
		}
	}
}
