# Linux Commands for Containerization

This file documents common Linux commands used in containerization, with a focus on process and filesystem isolation.

## `pivot_root`

The `pivot_root` system call moves the root filesystem of the current process to a different directory and mounts the old root filesystem at a specified location. It is a crucial tool for creating a new root environment for a container, effectively "pivoting" from the host's root to the container's root filesystem.

**Usage:**

```bash
pivot_root new_root put_old
```

*   `new_root`: The directory that will become the new root filesystem.
*   `put_old`: A directory under `new_root` where the old root filesystem will be mounted.

**Visual Explanation:**

Imagine the initial state of the filesystem tree:

```
      / (root)
      |
      +-- /dev
      +-- /proc
      +-- /sys
      +-- /home
      |   |
      |   +-- /user
      |       |
      |       +-- /new_root
      |           |
      |           +-- /put_old
      |           +-- /bin
      |           +-- /lib
      |           ...
      ...
```

After `pivot_root /home/user/new_root /home/user/new_root/put_old`, the process's view of the filesystem changes to:

```
      / (now points to new_root)
      |
      +-- /put_old (old root is mounted here)
      |   |
      |   +-- /dev
      |   +-- /proc
      |   +-- /sys
      |   +-- /home
      |       ...
      +-- /bin
      +-- /lib
      ...
```

The process is now contained within `new_root`, and the old root is accessible under `/put_old` (which can then be unmounted).

## `nsenter`

`nsenter` is a command-line tool that allows you to run a program with namespaces of other processes. This is incredibly useful for debugging containers or executing commands within a container from the host system.

**Usage:**

```bash
nsenter [options] -t <pid> [program [arguments]]
```

*   `-t <pid>`: The target process ID whose namespaces will be entered.
*   `[options]`: Can specify which namespaces to enter (e.g., `-m` for mount, `-u` for UTS, `-i` for IPC, `-n` for network, `-p` for PID, `-U` for user).

**Example:**

To enter the mount and network namespaces of a process with PID 1234 and run a shell:

```bash
nsenter -t 1234 -m -n /bin/bash
```

This will give you a shell inside the container-like environment of that process, allowing you to inspect its filesystems and network configuration as if you were inside it.

## `chroot`

`chroot` (change root) is a command that changes the root directory of the current running process and its children. A program that is run in such a modified environment cannot name (and therefore normally cannot access) files outside the designated directory tree.

**Usage:**
```bash
chroot NEW_ROOT [COMMAND [ARGS]...]
```

* `NEW_ROOT`: The directory to become the new root.
* `COMMAND`: The command to execute. If not specified, it defaults to the user's shell.

**Visual Explanation:**

If the filesystem looks like this:
```
      / (root)
      |
      +-- /bin
      +-- /lib
      +-- /home
          |
          +-- /user
              |
              +-- /jail
                  |
                  +-- /bin
                  +-- /lib
```

And you run `chroot /home/user/jail /bin/bash`, the new process will see `/home/user/jail` as its root (`/`). It won't be able to access files in the original `/bin`, `/lib`, or `/home`.

**`chroot` vs `pivot_root`**

- `chroot` is a system call, and it's simpler than `pivot_root`.
- `pivot_root` is more powerful and provides better isolation because it moves the old root away, making it possible to unmount it. With `chroot`, the original root remains a part of the process's filesystem, which can sometimes be escaped.
- `pivot_root` is used in modern containerization technologies, while `chroot` is an older and less secure method for filesystem isolation.

## `mount`

The `mount` command attaches a filesystem (from a device, a network share, or a special virtual filesystem) to a specified directory in the filesystem hierarchy, making its contents accessible. It is fundamental for setting up isolated environments in containers by controlling what filesystems are visible and where.

**Usage:**
```bash
mount [-t fstype] [-o options] device dir
```

*   `-t fstype`: Specifies the filesystem type. This is crucial for virtual filesystems used in containerization.
*   `-o options`: Specifies mount options, such as `ro` (read-only), `rw` (read-write), `bind` (bind mount), `defaults`, `loop`, etc.
*   `device`: The special file (e.g., `/dev/sdb1`), a directory for bind mounts, or a pseudo-filesystem name (e.g., `proc`).
*   `dir`: The directory where the filesystem will be mounted.

**Common `-t` (filesystem type) parameters for containerization:**

*   **`ext4`, `xfs`, etc.:** For mounting regular disk partitions or image files (via loop devices). This is typically used for the container's root filesystem.

*   **`proc`:** Mounts the `/proc` virtual filesystem.
    *   **Why it's used:** The `proc` filesystem provides an interface to kernel data structures. It's where information about processes, system memory, hardware configuration, and more is exposed as a hierarchical file-like structure.
    *   **When to use it in a container:** You should mount `/proc` in almost every container. Without it, many fundamental Linux commands like `ps`, `top`, `free`, `netstat`, and even `ls` (in some cases) will not work correctly because they rely on `/proc` to gather information about the system and its processes. When combined with a PID namespace, mounting `/proc` ensures that the container only sees its own processes.
    *   **Example:**
        ```bash
        # Mount the proc filesystem in the container's /proc directory
        mount -t proc proc /path/to/container/root/proc
        ```

*   **`sysfs`:** Mounts the `/sys` virtual filesystem.
    *   **Why it's used:** `sysfs` provides information about devices, drivers, and some kernel features.
    *   **When to use it in a container:** While not as universally critical as `/proc`, `/sys` is often mounted to provide read-only information about the system's hardware. Some applications and monitoring tools might require access to `/sys` to function correctly. For security, it's common practice to mount it as read-only.
    *   **Example:**
        ```bash
        # Mount the sysfs filesystem read-only
        mount -t sysfs sysfs /path/to/container/root/sys -o ro
        ```

*   **`tmpfs`:** Mounts a temporary filesystem that resides entirely in RAM.
    *   **Why it's used:** `tmpfs` is extremely fast because it avoids disk I/O. Data stored in a `tmpfs` mount is volatile and will be lost when the container stops or the filesystem is unmounted. This is useful for security and performance.
    *   **When to use it in a container:**
        *   For directories like `/tmp` and `/var/run` where applications store temporary files, sockets, or PID files that should not persist.
        *   For sharing secrets or sensitive data with a container. You can mount a `tmpfs` volume and write the secret to it, ensuring it never touches a physical disk and is automatically cleaned up.
        *   For performance-sensitive applications that need a high-speed scratch space for transient data.
    *   **Example:**
        ```bash
        # Mount a tmpfs for the /tmp directory with a size limit of 64MB
        mount -t tmpfs -o size=64m tmpfs /path/to/container/root/tmp
        ```

*   **`devtmpfs` or `devfs`:** Mounts a filesystem containing device files (e.g., `/dev/null`, `/dev/random`). Often combined with `mknod` to create specific devices.
    ```bash
    mount -t devtmpfs devtmpfs /path/to/container/root/dev
    ```
*   **`none` with `bind` option:** Creates a "bind mount," effectively making a directory or file available at another location in the filesystem hierarchy. This is extensively used in containers to share specific host directories or files without granting full access to the host's filesystem.
    ```bash
    mount --bind /host/path /path/to/container/root/guest/path
    # or using -o bind
    mount -o bind /host/path /path/to/container/root/guest/path
    ```

**Visual Explanation of a Bind Mount:**

Initial state:

```
      / (Host Root)
      |
      +-- /app (Host Application Code)
      |
      +-- /container_root
          |
          +-- /usr
          +-- /etc
          +-- /mnt/app_data (empty directory in container)
```

After `mount --bind /app /container_root/mnt/app_data`:

```
      / (Host Root)
      |
      +-- /app (Host Application Code)
      |   |
      |   +-- app.py
      |   +-- lib/
      |
      +-- /container_root
          |
          +-- /usr
          +-- /etc
          +-- /mnt/app_data <--- Now points to /app
              |
              +-- app.py
              +-- lib/
```

The container now sees the contents of `/host/path` at `/container_root/mnt/app_data`, effectively sharing the application code without copying it.

## `mknod` and `chown` for `/dev` Setup

The `mknod` command is used to create special files like device nodes (also known as device files) in the filesystem. These files provide an interface to hardware devices or pseudo-devices, which are crucial for the basic functionality of a Linux system, including within a container. The `chown` command is then used to set proper ownership for these device files.

**Usage Example for Container `/dev` Directory:**

```bash
mknod -m 666 "$ROOTFS_DIR/dev/null"    c 1 3
mknod -m 666 "$ROOTFS_DIR/dev/zero"    c 1 5
mknod -m 666 "$ROOTFS_DIR/dev/full"    c 1 7
mknod -m 666 "$ROOTFS_DIR/dev/random"  c 1 8
mknod -m 666 "$ROOTFS_DIR/dev/urandom" c 1 9
mknod -m 666 "$ROOTFS_DIR/dev/tty"     c 5 0

chown root:root "$ROOTFS_DIR/dev/"{null,zero,full,random,urandom,tty}
```

**Explanation:**

*   `mknod`: Creates a device file.
    *   `-m 666`: Sets the file permissions to `rw-rw-rw-`.
    *   `"$ROOTFS_DIR/dev/..."`: The path to the device file within the container's root filesystem (represented by `$ROOTFS_DIR`).
    *   `c`: Specifies a character device file (for sequential I/O). Block devices (for random access I/O) would use `b`.
    *   `1 3` (and other pairs): These are the **major** and **minor** device numbers.
        *   **Major Number**: Identifies the device driver. For example, `1` typically corresponds to the `mem` driver (which handles `/dev/null`, `/dev/zero`, etc.), and `5` often corresponds to the `tty` driver.
        *   **Minor Number**: Identifies a specific device or partition managed by that driver. For `mem` (major 1), minor `3` is `/dev/null`, minor `5` is `/dev/zero`, and so on. For `tty` (major 5), minor `0` is `/dev/tty`.

*   `chown root:root ...`: Changes the owner and group of the created device files to `root:root`, which is the standard ownership for most files in the `/dev` directory.
    *   `"$ROOTFS_DIR/dev/"{null,zero,full,random,urandom,tty}`: Uses shell brace expansion to concisely list all the device files that were just created.

**Purpose in Containerization:**

Creating these essential device files within a container's `/dev` directory is critical for its functionality:
*   `/dev/null`, `/dev/zero`, `/dev/full`: Fundamental pseudo-devices used by many applications for discarding output, providing zero-filled input, or simulating disk full conditions.
*   `/dev/random`, `/dev/urandom`: Provide sources of randomness, vital for cryptographic operations and secure applications.
*   `/dev/tty`: Allows processes to interact with their controlling terminal, essential for command-line applications.

Without these, a container would often fail to run basic commands or applications, as they expect these standard interfaces to the kernel's device management to be present.




## `unshare`

The `unshare` command creates new namespaces for a calling process and then executes a specified program within them. It is a high-level utility that provides a convenient wrapper around the `unshare()` and `clone()` system calls, making it a key tool for creating and experimenting with isolated environments from the command line.

### Comprehensive Example

The following command creates a new, highly isolated `bash` shell that resembles a basic container environment.

**Usage:**
```bash
sudo unshare --mount --pid --fork --cgroup --uts --net bash
```

*   `sudo`: Required for creating namespaces that demand elevated privileges (like the network and mount namespaces) and configuring them.
*   `--mount`: Creates a **new Mount namespace**. The shell gets an isolated view of the filesystem. Mounts and unmounts performed inside this shell will not affect the host's mount table.
*   `--pid`: Creates a **new PID namespace**. The new `bash` process will be PID 1 inside this namespace, unable to see or signal processes on the host.
*   `--fork`: Used with `--pid`, this forks a new child process to become PID 1 in the new PID namespace.
*   `--cgroup`: Creates a **new Cgroup namespace**. This isolates the shell's view of its own control groups, preventing it from seeing the host's full cgroup hierarchy.
*   `--uts`: Creates a **new UTS namespace**, giving the shell its own isolated hostname and domain name.
*   `--net`: Creates a **new Network namespace**. The shell gets its own private network stack, including network interfaces (initially just a down `lo` loopback), routing tables, and firewall rules, completely disconnecting it from the host's network.
*   `bash`: The program to run inside the newly created set of namespaces.

### `unshare` Command vs. Programmatic Creation

There are two primary ways to create namespaces: using the high-level `unshare` command or by making direct system calls (`clone`, `unshare`) in a low-level language like C or Go.

*   **`unshare` Command (High-Level):**
    *   **Pros:** Simple, concise, and ideal for shell scripting, system administration tasks, and experimentation. It allows you to create a complex, multi-namespace environment in a single, readable line.
    *   **Cons:** Less flexible. It is designed to run a command within new namespaces, but it doesn't offer fine-grained control over the setup process *before* the target program starts.

*   **Programmatic Syscalls (Low-Level, e.g., in C/Go):**
    *   **Pros:** Offers maximum power and flexibility. When building a container runtime, you need to perform many setup steps after creating namespaces but *before* executing the final user process. This includes setting up virtual network devices, mounting the root filesystem, changing the root with `pivot_root`, setting resource limits (cgroups), and dropping privileges. This can only be done programmatically.
    *   **Cons:** Far more complex and verbose. Requires writing and compiling a dedicated program. This is the approach taken by container runtimes like Docker/containerd and by the example code in `codes/how-container-works/main.go`.

In essence, the `unshare` command is a fantastic tool for learning and for simple isolation tasks, while the programmatic approach is essential for building robust, full-featured containerization software.


## `strace`

`strace` is a powerful diagnostic and debugging tool that intercepts and logs system calls made by a process and the signals it receives. It provides a low-level view of how a program interacts with the Linux kernel, making it invaluable for understanding container runtimes like `containerd`.

**Usage:**

```bash
# Terminal 1
sudo strace -f -qqq -e \
    trace=/clone,/execve,/unshare,/mount,/mknod,/mkdir,/link,/chdir,/chroot \
    -p $(pgrep containerd)
```

*   `sudo strace`: Run `strace` with root privileges, which are necessary to inspect other processes like `containerd`.
*   `-f`: Follows all child processes created by the traced process. This is crucial because container runtimes launch containers as new child processes.
*   `-qqq`: A "super quiet" mode that suppresses `strace`'s own attachment and detachment messages, resulting in cleaner output focused only on the traced calls.
*   `-p $(pgrep containerd)`: Attaches `strace` to an existing process. The `pgrep containerd` command finds the Process ID (PID) of the running `containerd` daemon.
*   `-e trace=...`: Specifies which system calls to trace. This example uses a custom set of calls that are fundamental to container creation:
    *   `clone`, `execve`: For creating new processes and executing programs within them. The `clone` syscall is the basis for creating new namespaces.
    *   `unshare`: For explicitly creating new namespaces for a process without creating a new process.
    *   `mount`: For setting up the container's filesystem, including the rootfs and virtual filesystems like `/proc`.
    *   `mknod`, `mkdir`, `link`: For creating filesystem structures like device nodes, directories, and hard links.
    *   `chdir`, `chroot`: For changing the working directory and, most importantly, changing the process's root directory to isolate its filesystem view.

**Why it's used:**

This command allows you to observe the precise, low-level sequence of kernel interactions that a container runtime like `containerd` performs to start a container. By watching these specific syscalls, you can see exactly how namespaces are created, how the container's root filesystem is constructed and isolated, and which binaries are executed. It is an essential tool for deep-diving into how containers actually work under the hood.
