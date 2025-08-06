#### Challenge

In this challenge, you will need to start a container using the default containerd CLI - ctr. Knowing how to use it may come in handy when you need to debug lower-level container issues (e.g., troubleshoot Kubernetes CRI on a containerd-powered cluster node).

#### First pull the image

sudo ctr image pull docker.io/library/nginx:alpine

sudo ctr create container docker.io/library/nginx:alpine nginx1

sudo ctr task start -d nginx1


#### Obtain the ID of the network namespace

sudo ctr task list

nginx1 13334 

sudo lsns -p 13334

#### How to obtain the IP using netns

You can use the `nsenter` command to enter the network namespace of a process.

Given the PID `13334`, you would use the following command:

```bash
sudo nsenter -t 13334 -n
```

This will execute a shell inside the network namespace of that process. From there, you can run commands like `ip a` or `ifconfig` to inspect its network interfaces and find the IP address.

If you just want to run a single command without starting an interactive shell, you can append it to the `nsenter` command:

```bash
sudo nsenter -t 13334 -n ip a
```

### Visual Explanation

Imagine your host system and the container as two separate houses, each with its own independent plumbing and electrical systems.

#### House 1: The Host System (Root Namespace)

This is your main machine. It has its own network setup: IP address, routing tables, list of network devices (`eth0`, `wlan0`, etc.). All the regular processes you run from your terminal live in this house.

```
+--------------------------------------------------+
|                  HOST SYSTEM                     |
|            (Root Network Namespace)              |
|                                                  |
|   +------------------+     +-----------------+   |
|   | Your Shell       |     | Other Processes |   |
|   | (e.g., bash)     |     | (e.g., chrome)  |   |
|   +------------------+     +-----------------+   |
|                                                  |
|   Network Stack:                                 |
|   - Device: eth0 (Physical NIC)                  |
|   - IP: 192.168.1.10                             |
|   - Can see the public internet                  |
|                                                  |
+--------------------------------------------------+
```

#### House 2: The Container (Container Namespace)

When a container starts, the Linux kernel builds a new, isolated "house" for it. This house gets its own private network setup. The process inside the container (like your app with PID 13334) only sees this private setup. It's completely unaware of the host's network.

```
+--------------------------------------------------+
|                  HOST SYSTEM                     |
|                                                  |
|   +------------------------------------------+   |
|   |           CONTAINER                      |   |
|   |      (Isolated Network Namespace)        |   |
|   |                                          |   |
|   |      +---------------------------+       |   |
|   |      | Your App (PID 13334)      |       |   |
|   |      +---------------------------+       |   |
|   |                                          |   |
|   |      Network Stack:                      |   |
|   |      - Device: eth0 (Virtual NIC)        |   |
|   |      - IP: 172.17.0.2                    |   |
|   |      - Can't see the host's NICs         |   |
|   |                                          |   |
|   +------------------------------------------+   |
|                                                  |
+--------------------------------------------------+
```

#### The `nsenter` Command: The Magic Door

The `nsenter` command is like creating a temporary, magic door from your house into the container's house.

When you run `sudo nsenter -t 13334 -n /bin/bash`:

1.  **`-t 13334` (target)**: You point at the process (PID 13334) living in the container. This tells `nsenter` which house to target.
2.  **`-n` (network)**: You specify that you are interested in the **network** namespace. You want to use the *plumbing* from the container's house.
3.  **`/bin/bash` (command)**: You choose the tool you want to bring through the door—in this case, a new shell.

The result is a new shell process that is running on the host, but looking through a window into the container's network world.

```
+--------------------------------------------------+
|                  HOST SYSTEM                     |
|                                                  |
|   +------------------------------------------+   |
|   |           CONTAINER                      |   |
|   |      (Isolated Network Namespace)        |   |
|   |                                          |   |
|   |      +---------------------------+       |   |
|   |      | Your App (PID 13334)      |       |   |
|   |      +---------------------------+       |   |
|   |                  ^                       |   |
|   |                  |                       |   |
|   +------------------|-----------------------+   |
|                      |                           |
|   +------------------+-----------------+         |
|   | Your NEW Shell (via nsenter)     |         |
|   |                                  |         |
|   |  Looks into container's network  |         |
|   |  - Sees IP 172.17.0.2            |         |
|   |  - Sees container's eth0         |         |
|   +----------------------------------+         |
|                                                  |
+--------------------------------------------------+
```

So, even though this new shell is technically a process on your host, `nsenter` has placed it *inside* the container's network namespace. When you run `ip a` in that shell, it reports the network devices and IPs that the container sees, not what the host sees.
