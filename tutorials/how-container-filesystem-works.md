### Tutorial: How Container Filesystems works: Building a docker-like container from scratch

#### Tools

- unshare
- mount 
- pivot_root

#### Command to mount a filesystem

```bash
mount --bind /tmp /mnt
```

This command uses `mount --bind` to create a bind mount. A bind mount makes a file or directory accessible at an alternative location in the filesystem. In this case, the content of `/tmp` is mirrored to `/mnt`.

#### Command to check mount points

```bash
cat /proc/self/mountinfo
```

This file provides detailed information about mount points in the process's mount namespace. It's more detailed than `/etc/mtab` and shows specifics like the mount ID, parent ID, propagation status (e.g., `shared`, `slave`), and other kernel-specific mount options. It's invaluable for debugging namespace issues.

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

#### Explaining the `nsenter` command

The `nsenter` command allows you to "enter" the namespaces of another running process. This is incredibly useful for debugging and inspection. Instead of creating a new, empty namespace with `unshare`, you can join an existing one to see what's happening inside it.

```bash
# Enter the mount namespace of process with PID 12345
sudo nsenter --mount=/proc/12345/ns/mnt

# Enter both the mount and PID namespaces of a process
sudo nsenter --mount=/proc/12345/ns/mnt --pid=/proc/12345/ns/pid
```

In this tutorial's example, `nsenter --mount=/proc/12345/ns/mnt` is used in Terminal 2 to make its shell run inside the exact same mount namespace as the shell in Terminal 1. This is key to setting up the master-slave relationship.


![container-image.png](./images/container-image.png)

#### Different types of propagation

- Shared: Changes propagate in both directions (master-to-slave and slave-to-master). This is useful for sharing mount points between containers or between the host and a container.

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

- Private: Changes do not propagate in either direction. This is the default and provides the strongest isolation.

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

- Slave: Changes propagate one-way, from the master to the slave. Events in the slave namespace do not propagate back to the master.

```
      Master Namespace           Slave Namespace
+---------------------+     +---------------------+
|   Mounts            |     |   Mounts            |
|  (e.g., mount /X)   |     |  (e.g., mount /Y)   |
+---------------------+     +---------------------+
              |                           ^
              |                           |
              +-----(Propagation)--------->
```

```bash
sudo unshare --mount --propagation slave bash
findmnt -o TARGET,SOURCE,FSTYPE,PROPAGATION
```

#### Explaining `findmnt`

The `findmnt` command is a utility to find a filesystem. It can search for filesystems in `/etc/fstab`, `/etc/mtab`, or `/proc/self/mountinfo`. In this tutorial, we use `findmnt -o TARGET,SOURCE,FSTYPE,PROPAGATION` to list the `TARGET` (mount point), `SOURCE` (device), `FSTYPE` (filesystem type), and `PROPAGATION` status of all mounts, which is perfect for observing the effects of our namespace and propagation changes.

#### Real OS Example of Mount Propagation

Here is a practical, step-by-step example you can run on a Linux machine to see the difference between propagation types. This example requires two separate terminal windows.

**Step 1: Prepare the Environment**

First, let's create a directory and a temporary filesystem (`tmpfs`) to work with. This needs to be done in a new mount namespace.

In **Terminal 1**, run the following:

```bash
# Become root to have permissions for mounting and namespaces.
sudo su

# Create a new mount namespace and make all mounts private initially.
unshare --mount --propagation private bash

# Create directories to work with. The -p flag creates parent directories if they don't exist.
mkdir -p /tmp/demo/master /tmp/demo/slave

# Create a new tmpfs and mount it on /tmp/demo.
# A tmpfs is a temporary filesystem that resides in RAM, making it very fast.
mount -t tmpfs -o size=1M tmpfs /tmp/demo

# Change the propagation of the new mount to 'shared'.
# This means mount events under /tmp/demo will propagate to other namespaces.
mount --make-shared /tmp/demo

# Get the PID of this shell to use in the next step.
echo "Terminal 1 PID is $$"
  # Note this PID for the next step. Let's assume it's 12345 for this example.
```

**Step 2: Create a Slave Namespace**

Now, in **Terminal 2**, we will create a new namespace that is a slave to the one in Terminal 1.

```bash
# Become root.
sudo su

# Enter the mount namespace of Terminal 1 using its PID.
# From now on, this terminal sees the same filesystem layout as Terminal 1.
nsenter --mount=/proc/12345/ns/mnt

# Create a new shell with a mount namespace that is a SLAVE to the current one.
# This means it will receive mount events from the master, but not send them.
unshare --mount --propagation slave bash

# Mount our shared tmpfs onto the slave directory.
# This makes /tmp/demo/slave a mount point so we can inspect its 'slave' property.
mount /dev/shm /tmp/demo/slave

# Check the mounts. Note the propagation type of 'slave'.
findmnt -o TARGET,PROPAGATION | grep /tmp/demo
# Expected output:
# /tmp/demo        shared
# /tmp/demo/slave  slave
```

**Step 3: Test Propagation from Master to Slave**

Now, let's create a new mount point inside the master's shared directory and see if it appears in the slave.

In **Terminal 1** (the master):

```bash
# Mount a new tmpfs inside the shared directory.
mount -t tmpfs -o size=1M tmpfs /tmp/demo/master-test

# Check if it was created.
ls /tmp/demo/master-test
# You should see the 'lost+found' directory.
```

Now, check if this new mount propagated to the slave.

In **Terminal 2** (the slave):

```bash
# Check for the new mount inside the slave's view.
ls /tmp/demo/slave/master-test
# You should see the 'lost+found' directory here too!
# The mount event propagated from the master to the slave.
```

**Step 4: Test Propagation from Slave to Master**

Finally, let's create a new mount in the slave namespace and see if it propagates back to the master.

In **Terminal 2** (the slave):

```bash
# Mount a new tmpfs inside the slave directory.
mount -t tmpfs -o size=1M tmpfs /tmp/demo/slave/slave-test

# Check if it was created.
ls /tmp/demo/slave/slave-test
# You should see 'lost+found'.
```

Now, let's check if this new mount appeared in the master.

In **Terminal 1** (the master):

```bash
# Check for the slave's mount.
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

