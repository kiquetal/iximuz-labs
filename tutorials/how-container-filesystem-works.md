### Tutorial: How Container Filesystems works: Building a docker-like container from scratch

#### Tools

- unshare
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



#### Different types of propagation

- Private: Changes do not propagate in either direction.

```
      Namespace A                Namespace B
+---------------------+     +---------------------+
|   Mounts            |     |   Mounts            |
|   (e.g., mount /X)  |     |   (e.g., mount /Y)  |
+---------------------+     +---------------------+
        ^                           ^
        |                           |
        +-----(No Propagation)------+
```

```bash
sudo unshare --mount --propagation private bash
findmnt -o TARGET,SOURCE,FSTYPE,PROPAGATION
```

- Shared: Changes propagate in both directions.

```
      Namespace A                Namespace B
+---------------------+     +---------------------+
|   Mounts            |     |   Mounts            |
|   (e.g., mount /X)  |     |   (e.g., mount /X)  |
+---------------------+     +---------------------+
        ^                           ^
        |                           |
        +-----(Bidirectional)-------+
           <-- Propagation -->
```

```bash
sudo unshare --mount --propagation shared bash
findmnt -o TARGET,SOURCE,FSTYPE,PROPAGATION
```

#### Real OS Example of Mount Propagation

Here is a practical, step-by-step example you can run on a Linux machine to see the difference between propagation types. This example requires two separate terminal windows.

**Step 1: Prepare the Environment**

First, let's create a directory and a temporary filesystem (`tmpfs`) to work with. This needs to be done in a new mount namespace.

In **Terminal 1**, run the following:

```bash
# Become root
sudo su

# Create a new mount namespace and make all mounts private initially
unshare --mount --propagation private bash

# Create directories to work with
mkdir -p /tmp/demo/master /tmp/demo/slave

# Create a new tmpfs and mount it on /tmp/demo, making it shared
mount -t tmpfs -o size=1M tmpfs /tmp/demo
mount --make-shared /tmp/demo

# Get the PID of this shell
echo "Terminal 1 PID is $$"
  # Note this PID for the next step. Let's assume it's 12345 for this example.
```

**Step 2: Create a Slave Namespace**

Now, in **Terminal 2**, we will create a new namespace that is a slave to the one in Terminal 1.

```bash
# Become root
sudo su

# Enter the mount namespace of Terminal 1 (replace 12345 with the actual PID)
nsenter --mount=/proc/12345/ns/mnt

# Create a new shell with a mount namespace that is a slave to the current one
unshare --mount --propagation slave bash

# Mount our shared tmpfs onto the slave directory
mount /dev/shm /tmp/demo/slave

# Check the mounts. Note the propagation type.
findmnt -o TARGET,PROPAGATION | grep /tmp/demo
# Expected output:
# /tmp/demo        shared
# /tmp/demo/slave  slave
```

**Step 3: Test Propagation from Master to Slave**

Now, let's create a new mount point inside the master's shared directory and see if it appears in the slave.

In **Terminal 1** (the master):

```bash
# Mount a new tmpfs inside the shared directory
mount -t tmpfs -o size=1M tmpfs /tmp/demo/master-test

# Check if it was created
ls /tmp/demo/master-test
# You should see the 'lost+found' directory.
```

Now, check if this new mount propagated to the slave.

In **Terminal 2** (the slave):

```bash
# Check for the new mount inside the slave's view
ls /tmp/demo/slave/master-test
# You should see the 'lost+found' directory here too!
# The mount event propagated from the master to the slave.
```

**Step 4: Test Propagation from Slave to Master**

Finally, let's create a new mount in the slave namespace and see if it propagates back to the master.

In **Terminal 2** (the slave):

```bash
# Mount a new tmpfs inside the slave directory
mount -t tmpfs -o size=1M tmpfs /tmp/demo/slave/slave-test

# Check if it was created
ls /tmp/demo/slave/slave-test
# You should see 'lost+found'.
```

Now, let's check if this new mount appeared in the master.

In **Terminal 1** (the master):

```bash
# Check for the slave's mount
ls /tmp/demo/slave/slave-test
# ls: cannot access '/tmp/demo/slave/slave-test': No such file or directory
# The mount event did NOT propagate back from the slave to the master.
```

This example clearly demonstrates that with a `slave` propagation type, mount events only flow one way: from the master to the slave.

```
   Master Namespace          Slave Namespace
+------------------+       +------------------+
|   Mounts         |------>|   Mounts         |
| (e.g., mount /A) |       | (receives /A)    |
+------------------+       +------------------+
|                  |       |                  |
|   mount /B       |<--X-- |   mount /C       |
| (not propagated) |       | (not propagated) |
+------------------+       +------------------+

```

