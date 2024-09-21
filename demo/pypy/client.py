import os
import time
import sys


PIPE_CS_FILE = "/tmp/client-server"
PIPE_SC_FILE = "/tmp/server-client"


def check_pipe_exist() -> None:
    if not os.path.exists(PIPE_CS_FILE):
        os.mkfifo(PIPE_CS_FILE)
    if not os.path.exists(PIPE_SC_FILE):
        os.mkfifo(PIPE_SC_FILE)


def call(pipe_cs, pipe_sc) -> None:
    # send message
    message = "100\n"
    pipe_cs.write(message)
    pipe_cs.flush()

    # receive
    response = pipe_sc.readline().strip()
    # if response:
    #     print(f"Received: {response}")


def main():
    check_pipe_exist()

    start_time = time.perf_counter()

    loop = int(sys.argv[1])
    with open(PIPE_CS_FILE, 'w') as pipe_cs:
        with open(PIPE_SC_FILE, 'r') as pipe_sc:
            for _ in range(loop):
                call(pipe_cs, pipe_sc)

    end_time = time.perf_counter()
    elapse_time = end_time - start_time
    print(f"Total time: {round(elapse_time, 3)}s")
    print(f"Average time: {round(elapse_time / loop * 1_000_000)}µs")


if __name__ == "__main__":
    main()
