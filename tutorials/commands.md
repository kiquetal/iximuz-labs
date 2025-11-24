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

*   **`ext4`, `xfs`, etc.:** For mounting regular disk partitions or image files (via loop devices).
*   **`proc`:** Mounts the `/proc` virtual filesystem, providing process information. Essential for many Linux utilities to function correctly inside a container.
    ```bash
    mount -t proc proc /path/to/container/root/proc
    ```
*   **`sysfs`:** Mounts the `/sys` virtual filesystem, exposing kernel objects and their attributes.
    ```bash
    mount -t sysfs sysfs /path/to/container/root/sys
    ```
*   **`tmpfs`:** Mounts a temporary filesystem residing in RAM. Often used for `/tmp` or other transient data within containers.
    ```bash
    mount -t tmpfs tmpfs /path/to/container/root/tmp
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


## `unshare`

The `unshare` command creates new namespaces for the calling process and then executes a specified program. This is a key command for isolating a process from the main system.

**Usage:**

```bash
unshare [options] [program [arguments]]
```

*   `[options]`: These flags determine which namespaces to unshare (create).
    *   `--fork`: Create a new process and unshare namespaces for it.
    *   `--pid`: Unshare the PID namespace.
    *   `--mount-proc`: Mount the `/proc` filesystem, which is necessary when creating a new PID namespace.
    *   `--uts`: Unshare the UTS namespace (hostname).
    *   `--ipc`: Unshare the IPC namespace.
    *   `--net`: Unshare the network namespace.
    *   `--user`: Unshare the user namespace.
    *   `--map-root-user`: Map the current user to the root user in the new user namespace.
*   `[program [arguments]]`: The program to run in the new namespaces.

**Example:**

To create a new process with its own PID and UTS namespaces, and run a shell inside it:

```bash
unshare --fork --pid --mount-proc --uts /bin/bash
```

Inside this new shell:
*   The process ID will be 1.
*   You can set a new hostname that will only be visible within this shell.

This command is fundamental for creating isolated environments that form the basis of containers.
