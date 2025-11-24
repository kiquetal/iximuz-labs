### Tutorial: How Container Filesystems works: Building a docker-like container from scratch

#### Tools

- unsahre
- mount 
- pivot_root



![container-image.png](./images/container-image.png)

#### Command to mount a filesystem

```bash
mount --bind /tmp /mnt
```

This command uses `mount --bind` to create a bind mount. A bind mount makes a file or directory accessible at an alternative location in the filesystem. In this case, the content of `/tmp` is mirrored to `/mnt`.





#### Command to check mount points

```bash

cat /proc/self/mountinfo
```

#### Command to inspect a process's mount namespace

```bash
readlink /proc/$PID/ns/mnt
```

The `readlink` command is used to display the value of a symbolic link. In the context of containers, `readlink /proc/$PID/ns/mnt` is particularly useful. It shows the inode number of the mount namespace for a given process ID (`$PID`). By comparing the output for different processes, you can determine if they share the same mount namespace (i.e., see the same filesystem hierarchy) or if they are isolated in their own namespaces. This is a fundamental concept for understanding how containers provide filesystem isolation.

#### Explaining the `unshare` command

The `unshare` command in Linux is used to run a program with new namespaces. Namespaces are a fundamental Linux kernel feature that provides resource isolation for processes. This is how containers achieve their isolation from the host system.

Here are some common examples of using `unshare`:

**1. Unsharing the Mount Namespace (`-m` or `--mount`)**

This creates a new mount namespace, allowing processes within it to have an independent view of the filesystem mount points. Changes to mount points inside this new namespace do not affect the parent namespace. However, the reverse is not true by default. Due to mount propagation (specifically shared subtrees), mount events from the parent namespace can still propagate into the new namespace. To achieve true isolation, you must change the propagation policy.

```bash
# Create a new mount namespace
sudo unshare --mount bash

# The new namespace can still receive mount events from the parent.
# To prevent this, mark the entire filesystem as a private mount recursively.
mount --make-rprivate /

# Now you are in a new shell with a truly isolated mount namespace.
# Mount/unmount events will not propagate in either direction.
# Example: mount -t tmpfs none /tmp_new
# Then exit to return to the original namespace.
```

**2. Unsharing the PID Namespace (`-p` or `--pid`)**

This creates a new PID namespace, where the first process in the new namespace has PID 1. Processes inside this namespace have a different view of process IDs compared to the parent namespace.

```bash
sudo unshare --pid --fork bash
# In the new shell, 'ps aux' will show a different set of PIDs,
# and the shell itself will likely be PID 1 within this new namespace.
```

**3. Unsharing the UTS Namespace (`-u` or `--uts`)**

This creates a new UTS (UNIX Time-sharing System) namespace, allowing the hostname and NIS domain name to be changed without affecting the host system.

```bash
sudo unshare --uts bash
hostname new-container-hostname
# The hostname is changed only for this new namespace.
```

**4. Unsharing Multiple Namespaces**

You can combine multiple options to unshare several namespaces simultaneously, which is typical for container environments.

```bash
sudo unshare --mount --pid --uts --fork bash
# This creates a new environment with isolated mount, PID, and UTS namespaces.
```

The `unshare` command is a crucial tool for understanding and building container technologies, as it directly manipulates the namespaces that provide isolation.



