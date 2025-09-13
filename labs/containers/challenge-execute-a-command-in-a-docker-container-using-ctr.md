### Execute a Command in a Docker Container Using ctr

### First we obtain the information about process already running via docker

- sudo ctr -n moby task ls

### Containers vs. Tasks in Containerd

In `containerd`'s terminology, a container and a task are different concepts:

-   **Container**: This is a metadata object that defines the container. It includes the container's ID, the OCI specification (what command to run, environment variables, etc.), and a reference to the root filesystem (the snapshot of the image). A container object is static; it describes what the running container *should* be. You can list them with `sudo ctr -n moby container ls`.

-   **Task**: This is the live, running process that executes inside the container. When you start a container, `containerd` creates a `task`. The task is what has a PID and consumes system resources. The `ctr task ls` command lists these active tasks. A single container can have multiple tasks if you use `ctr task exec` to run additional processes.

The `ctr task exec` command requires the **Container ID**, not the task's process ID (PID). The output of `sudo ctr -n moby task ls` conveniently shows both the `CONTAINER` ID and the `PID` for each running task, which is why it's a useful first step. You should use the value from the `CONTAINER` column in the `exec` command.

### Then we use the following command

- sudo ctr -n moby task exec --exec-id 0 --tty <container-id> sh

### Command Explanation

The `ctr task exec` command allows you to run a new process inside an already running container. Let's break down the flags used:

-   `--exec-id <id>`: This flag provides a unique identifier for the new process you are starting within the container. This ID is crucial for managing the process later, for example, to send signals to it or to kill it without affecting the main container process. You can choose any unique string for the ID; in the example, `0` is used, but you could use something more descriptive like `debug-shell`.

    **Example:**
    If you start a process with `--exec-id my-debug-session`, you can later interact with that specific process using this ID.

-   `--tty` or `-t`: This allocates a pseudo-terminal (TTY), which is essential for interactive sessions like a shell. It connects your terminal to the container's process, allowing you to type commands and see the output.

-   `<container-id>`: This is the target container where the command will be executed. You can get this ID from the `CONTAINER` column in the output of `sudo ctr -n moby task ls`.

-   `sh`: This is the command to run inside the container. In this case, it's launching a Bourne shell, giving you a command prompt inside the container. You could replace `sh` with any other command, like `ls /` to list the root directory's contents.

**Example Usage:**

1.  List running tasks to find your container ID:
    ```bash
    sudo ctr -n moby task ls
    ```

2.  Execute an interactive shell in the container:
    ```bash
    sudo ctr -n moby task exec --exec-id shell-1 --tty my-container-abc sh
    ```

3.  Execute a non-interactive command:
    ```bash
    sudo ctr -n moby task exec --exec-id list-files my-container-abc ls -l /app
    ```



