# Linux Namespaces

Linux namespaces are a feature of the Linux kernel that partitions kernel resources such that one set of processes sees one set of resources while another set of processes sees a different set of resources. Namespaces are a fundamental aspect of containers on Linux.

## Types of Namespaces

There are several types of namespaces:

*   **PID (Process ID):** Isolates the process ID number space. This means that a process inside a PID namespace can have the PID 1, but there can be another process with PID 1 in another namespace.
*   **UTS (UNIX Timesharing System):** Isolates the hostname and NIS domain name. This allows each container to have its own hostname.
*   **IPC (Inter-Process Communication):** Isolates System V IPC and POSIX message queues.
*   **Network:** Isolates network devices, stacks, ports, etc. Each network namespace can have its own virtual network device and IP addresses.
*   **Mount:** Isolates the filesystem mount points. This allows processes in different mount namespaces to have different views of the filesystem hierarchy.
*   **User:** Isolates user and group IDs. A process can have a normal unprivileged user ID outside a user namespace while having a user ID of 0 (root) inside the namespace.

## `unshare` in Go

The `unshare` functionality can be used in Go to create new namespaces programmatically. The `syscall` package provides the necessary constants and functions.

Here is an example of how to use `unshare` in Go to create a new process with its own UTS namespace.

```go
package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	cmd := exec.Command("/bin/bash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Create a new UTS namespace
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,
	}

	// Set a new hostname in the new namespace
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return
	}

	// Set hostname in the new namespace
	setHostnameCmd := exec.Command("hostname", "-b", "new-hostname")
	setHostnameCmd.Run()

	fmt.Println("Hostname in new namespace should be 'new-hostname'")
	fmt.Println("You are in a new UTS namespace. Check the hostname by typing `hostname`.")

	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for command:", err)
	}
}
```

In this example:
1.  We create a new command to run `/bin/bash`.
2.  We set the `Cloneflags` to `syscall.CLONE_NEWUTS` to create a new UTS namespace.
3.  We start the command.
4.  We then run the `hostname` command to set a new hostname within that namespace.
5.  When you run this program, you will be in a new shell where the hostname is 'new-hostname'. If you open another terminal, you will see the original hostname.

This demonstrates how Go can be used to create isolated environments, which is the foundation of containerization.

### PID Namespace Example

This example creates a new process in a new PID namespace. Inside this namespace, the process will have PID 1.

```go
package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	cmd := exec.Command("/bin/bash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Create new PID and UTS namespaces
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS,
	}

	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}
```

When you run this code, you'll get a new shell. If you run `ps aux` inside this shell, you will see that the `bash` process has a PID of 1.

### Mount Namespace Example

This example demonstrates creating a new mount namespace and changing the root filesystem within that namespace.

```go
package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	cmd := exec.Command("/bin/bash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Create a new mount namespace
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS,
	}

	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}
```

After running this code, you are in a new mount namespace. Any mounts you create or unmount will not affect the host's mount table. For example, you can remount the root filesystem as read-only, and it will only apply to this namespace.

```bash
mount -o remount,ro /
```
This will fail without root privileges, but it illustrates the isolation.