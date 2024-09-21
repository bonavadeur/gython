# Named-pipe IPC between Go and Python

[![LICENSE](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)
![Linux](https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black)
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Python](https://img.shields.io/badge/python-3670A0?style=for-the-badge&logo=python&logoColor=ffdd54)

Call Python function from Golang code with overhead in tens of microseconds.

## 1. Introduction

This repo introduces how to call a Python function from a Golang code, using a named pipeline IPC (Inter-Process Communication). Although this type of IPC has a very high penalty, in my case the delay is acceptable.

## 2. Experiments

We have one client code (Go or Python) which need to call a function from another server code (Go or Python). On each call, the client passes value "100\n" into Linux named-pipeline to the server, then server adds 1 to this value and returns to client.

First, you need to start a Server (Python), then start a Client. I don't know why Go Server doesn't work with Python Client :D.

```bash
# start go server
go run demo/gogo/server/main.go
# start python server
python3 demo/pypy/server.py
# start go client with <number-of-calls>
# example: go run demo/gogo/client/main.go 100
go run demo/gogo/client/main.go <number-of-calls>
# start python client with <number-of-calls>
# example: python3 demo/pypy/client.py 100
python3 demo/pypy/client.py <number-of-calls>
```

If this function (the plus-one-function) is implemented inside one programming language, the execution time will be in nanoseconds only. The result indicate that, by using this type of IPC, the penalty will be in tens of microseconds. In my usecase, it meets my requirements perfectly.

### Result of average response time of each call (µs)

| | 100 | 1K | 10K | 100K | 1M |
|---|---|---|---|---|---|
| Go-Go | 26 | 24 | 19 | 18 | 17 |
| Go-Python | 23 | 21 | 20 | 19 | 18 |
| Python-Go | NaN | NaN | NaN | NaN | NaN | 
| Python-Python | 24 | 21 | 19 | 19 | 19 |

## 3. Intergration

### 3.1. Problem

We have a Go code, specially a web server written in Go using Echo framework. This web server exposes API `GET /native/<sleepTime>` which will return `sleepTime` and time for processing this request.

```bash
# let modify this tag to be :root
$ cat manifests/gython.yaml | grep image:
        image: docker.io/bonavadeur/gython:root

# modify this ConfigMap
$ cat manifests/gython.yaml | grep enable-python:
    enable-python: "false"

$ kubectl apply -f manifests/gython.yaml
service/gython created
deployment.apps/gython created
configmap/gython created

$ curl gython.default/native/10
{"result":"10","execTime":"10.249563ms"}
```

This API let server sleeps in 10ms and return `sleepTime`. Total time actually elapsed is 10.249563ms. **Sleep** is a fake processing procedure that executed inside Go code. Now, we need execute this procedure in Python (for easier implementation other algorithms). You can compare two functions do the same work in `cmd/gython/main.go func fakeProcessing()` and `outlier/main.py def fake_processing()`

```go
// cmd/gython/main.go
func fakeProcessing(sleepTime string) string {
	sleep, _ := strconv.Atoi(sleepTime)
	time.Sleep(time.Duration(sleep) * time.Millisecond)
	return sleepTime
}
```

```python
# outlier/main.py
def fake_processing(n : str) -> str:
    time.sleep(float(n) / 1000)
    return n
```

Python code is shorter significantly than Go code, and much easier for algorithm implementing. So, my work is design a pattern for calling Python func from Go code, with acceptable latency.

### 3.2. Experiment

#### Build Go code

Take a quick look in [build.sh](build.sh), let configure registry and image name. Build Go code to Docker image:

```bash
$ ./build.sh ko
$ ./build.sh push root
```

Assume that you have an closed image `docker.io/bonavadeur/gython:root`. This image is a web server that serve API `GET /native/10` in Go code. Now, you will define another Python code located in [outlier/main.py]([outlier/main.py]) do the same work, and design another API `GET /pipe/10` to do the same work in Go code. You will build another Docker image named `docker.io/bonavadeur/gython:python` base on the first image `docker.io/bonavadeur/gython:root`. Go code in base image can interact with new Python code.

```bash
# let modify this tag to be :python
$ cat manifests/gython.yaml | grep image:
        image: docker.io/bonavadeur/gython:python

# modify this ConfigMap
$ cat manifests/gython.yaml | grep enable-python:
    enable-python: "true"

# base root image, build new image with Python code
$ docker build -t bonavadeur/gython:python .
$ docker push bonavadeur/gython:python

# delete current gython pod
$ kubectl get pod | grep gython
gython-6f84b6ffb7-kchpx             1/1     Running   0                 22m

$ kubectl delete pod gython-6f84b6ffb7-kchpx
pod "gython-6f84b6ffb7-kchpx" deleted

$ kubectl get pod | grep gython
gython-6f84b6ffb7-gzvcv             1/1     Running   0                 12s
```

Check if API work correctly:

```bash
$ curl gython.default/native/10
{"result":"10","execTime":"10.251909ms"}
$ curl gython.default/pipe/10
{"result":"10","execTime":"10.399817ms"}
```

As you can see, in doing the same works, Go code takes 10.251909ms, Python code take 10.399817ms, 87.9µs slower. The overhead for calling from Go to Python is about 90µs.

### 3.3. For Development

```bash
# build root Go code
$ ./build.sh ko
$ ./build.sh push root

# build Python intergrated image:
$ docker build -t bonavadeur/gython:python .
$ docker push bonavadeur/gython:python
```

## 4. Épilogue

This technique is used in [Katyusha](https://github.com/bonavadeur/katyusha) and [Nonna](https://github.com/bonavadeur/nonna)

## 5. Contributeur

Đào Hiệp - Bonavadeur - ボナちゃん  
The Future Internet Laboratory, E711 C7 Building, Hanoi University of Science and Technology, Vietnam.   
未来のインターネット研究室, C7 の E ７１１、ハノイ百科大学、ベトナム。  
![](images/github-wp.png)
