import os
import time
import math


class Pipe():
    def __init__(self, upstream_file: str, downstream_file: str, buffer: int) -> None:
        # check if pipe_file exist
        if not os.path.exists(upstream_file):
            os.mkfifo(upstream_file)
        if not os.path.exists(downstream_file):
            os.mkfifo(downstream_file)

        self.upstream_file = upstream_file
        self.downstream_file = downstream_file
        self.buffer = buffer
        self.upstream = open(self.upstream_file, 'r')
        self.downstream = open(self.downstream_file, 'w')


    def read(self) -> str:
        while True:
            while True:
                try:
                    data = self.upstream.readline().strip()
                    if not data:
                        return None
                except Exception as e:
                    return None
                return data


    def write(self, message: str) -> None:
        message = f"{message}"
        self.downstream.write(message)
        self.downstream.flush()


    def __del__(self):
        self.upstream.close()
        self.downstream.close()


def fake_processing(n : str) -> str:
    time.sleep(float(n) / 1000)
    return n


def main():
    pipe = Pipe(
        upstream_file="/tmp/upstream",
        downstream_file="/tmp/downstream",
        buffer=1024,
    )
    while True:
        message = int(pipe.read())
        if message != None:
            resp = fake_processing(message)
        pipe.write(resp)


if __name__ == "__main__":
    main()
