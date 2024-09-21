package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/bonavadeur/gython/pkg/pipe"
	"github.com/labstack/echo/v4"
)

var (
	PORT       = ":8080"
	UPSTREAM   = "/tmp/upstream"
	DOWNSTREAM = "/tmp/downstream"
	PIPE       *pipe.Pipe
)

type Response struct {
	Result   string `json:"result"`
	ExecTime string `json:"execTime"`
}

func NewResponse() *Response {
	return &Response{
		Result:   "",
		ExecTime: "",
	}
}

func fakeProcessing(sleepTime string) string {
	sleep, _ := strconv.Atoi(sleepTime)
	time.Sleep(time.Duration(sleep) * time.Millisecond)
	return sleepTime
}

func callByNative(c echo.Context) error {
	param := c.Param("param")
	response := NewResponse()

	// processing
	startTime := time.Now()
	response.Result = fakeProcessing(param)
	elapseTime := time.Since(startTime)
	response.ExecTime = elapseTime.String()

	return c.JSON(http.StatusOK, response)
}

func callByPipe(c echo.Context) error {
	response := NewResponse()

	// processing
	startTime := time.Now()
	var err error
	PIPE.Write(c.Param("param"))
	response.Result, err = PIPE.Read()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	elapseTime := time.Since(startTime)
	response.ExecTime = elapseTime.String()

	return c.JSON(http.StatusOK, response)
}

func init() {
	enablePython := os.Getenv("enable-python")
	if enablePython == "false" {

	} else if enablePython == "true" {
		var err error
		PIPE, err = pipe.NewPipe(UPSTREAM, DOWNSTREAM, 1024)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Konnichiwa!\n")
	})
	e.GET("/native/:param", callByNative)
	e.GET("/pipe/:param", callByPipe)
	e.Logger.Fatal(e.Start(PORT))
	fmt.Println("Creating pipe")
}
