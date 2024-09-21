package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"
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

	pipeCS, err := os.OpenFile(pipeCSFile, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		fmt.Println("Error opening pipe for writing:", err)
		return
	}
	defer pipeCS.Close()

	pipeSC, err := os.OpenFile(pipeSCFile, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		fmt.Println("Error opening pipe for reading:", err)
		return
	}
	defer pipeSC.Close()

	// benchmark
	start := time.Now()
	loop, _ := strconv.Atoi(os.Args[1])
	for i := 1; i <= loop; i++ {
		call(buffer, pipeCS, pipeSC)
	}
	elapse := time.Since(start)
	avgResp := elapse.Microseconds() / int64(loop)
	fmt.Println("Total time:", elapse)
	fmt.Printf("Average response time: %dµs\n", avgResp)
}

func checkPipeExist(pipe string) {
	if _, err := os.Stat(pipe); os.IsNotExist(err) {
		if err := syscall.Mkfifo(pipe, 0777); err != nil {
			fmt.Println("Error creating named pipe:", err)
			return
		}
	}
}

func call(buffer []byte, pipeCS *os.File, pipeSC *os.File) {
	// send
	message := "100\n"
	_, err := pipeCS.Write([]byte(message))
	if err != nil {
		fmt.Println("Error writing to pipeCS:", err)
		return
	}

	// receive
	for {
		n, err := pipeSC.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("End of data reached, closing pipe")
				break
			} else {
				fmt.Println("Error reading from pipe:", err)
				return
			}
		}

		if n > 0 {
			// message := string(buffer[:n])
			// fmt.Println(message)
			break
		}
	}
}
