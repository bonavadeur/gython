package pipe

import (
	"fmt"
	"os"
	"syscall"
)

type Pipe struct {
	UpstreamFile   string
	DownstreamFile string
	Upstream       *os.File
	Downstream     *os.File
	Buffer         []byte
}

func NewPipe(
	upstreamFile string,
	downstreamFile string,
	bufferSize int,
) (*Pipe, error) {
	var err error

	// check if PipeFile exists
	if _, err := os.Stat(upstreamFile); os.IsNotExist(err) {
		if err := syscall.Mkfifo(upstreamFile, 0777); err != nil {
			return nil, err
		}
        if err := os.Chmod(upstreamFile, 0777); err != nil {
            return nil, err
        }
	}
	if _, err := os.Stat(downstreamFile); os.IsNotExist(err) {
		if err := syscall.Mkfifo(downstreamFile, 0777); err != nil {
			return nil, err
		}
		if err := os.Chmod(downstreamFile, 0777); err != nil {
            return nil, err
        }
	}

	newPipe := &Pipe{
		UpstreamFile:   upstreamFile,
		DownstreamFile: downstreamFile,
		Upstream:       nil,
		Downstream:     nil,
		Buffer:         make([]byte, bufferSize),
	}

	// open PipeFile
	newPipe.Upstream, err = os.OpenFile(upstreamFile, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		return nil, err
	}
	newPipe.Downstream, err = os.OpenFile(downstreamFile, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		return nil, err
	}

	return newPipe, nil
}

func (p *Pipe) Close() {
	p.Upstream.Close()
	p.Downstream.Close()
}

func (p *Pipe) Write(message string) error {
	messageByte := []byte(message)
	messageByte = append(messageByte, '\n')
	_, err := p.Upstream.Write(messageByte)
	if err != nil {
		return err
	}
	return nil
}

func (p *Pipe) Read() (string, error) {
	for {
		n, err := p.Downstream.Read(p.Buffer)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("End of data reached, closing pipe")
				return "", err
			} else {
				fmt.Println("Error reading from pipe:", err)
				return "", err
			}
		}

		if n > 0 {
			message := string(p.Buffer[:n])
			return message, nil
		}
	}
}
