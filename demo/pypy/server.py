import os


PIPE_CS_FILE = "/tmp/client-server"
PIPE_SC_FILE = "/tmp/server-client"


def check_pipe_exist() -> None:
    if not os.path.exists(PIPE_CS_FILE):
        os.mkfifo(PIPE_CS_FILE)
    if not os.path.exists(PIPE_SC_FILE):
        os.mkfifo(PIPE_SC_FILE)


def main():
    check_pipe_exist()

    # receive
    with open(PIPE_CS_FILE, 'r') as pipe_cs:
        with open(PIPE_SC_FILE, 'w') as pipe_sc:
            while True:
                while True:
                    try:
                        data = pipe_cs.readline().strip()
                        if not data:
                            break
                    except Exception as e:
                        print(f"Error: {e}")
                        break
                    
                    # Process data
                    i = int(data)
                    i += 1

                    # send
                    message = f"{i}\n"
                    pipe_sc.write(message)
                    pipe_sc.flush()

if __name__ == "__main__":
    main()
