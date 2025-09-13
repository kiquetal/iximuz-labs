# Challenge: Execute an HTTP Request in a Container Without Exposed Ports

This guide explains how to access a web service running inside a container even when its ports are not published to the host. This technique involves finding the container's process ID (PID) and using the `nsenter` utility to execute commands inside the container's isolated network namespace.

## Concept Diagram

The diagram below illustrates how `nsenter` provides a "portal" from the host's default namespace into the container's private network namespace, allowing commands to run as if they were inside the container.

```text
+-------------------------------------------------------------------------+
| Host System (Default Namespace)                                         |
|                                                                         |
|   +-----------------------------------------------------------------+   |
|   | Container's Network Namespace (Isolated)                         |   |
|   |                                                                 |   |
|   |   +-----------------------------+                               |   |
|   |   | Container                   |                               |   |
|   |   |                             |                               |   |
|   |   |   [ Web App on port 80 ]    <-----+                         |   |
|   |   |   (localhost:80)            |     |                         |   |
|   |   |                             |     |                         |   |
|   |   +-----------------------------+     |                         |   |
|   |                                       |                         |   |
|   +-----------------------------------------------------------------+   |
|                                           ^                             |
|   Host PID for Container: [67890]         |                             |
|                 |                         |                             |
|                 |     +-----------------------------------------+       |
|                 +-----> sudo nsenter --target 67890 --net curl ...|       |
|                       +-----------------------------------------+       |
|   [ Host Shell ]                                                        |
|                                                                         |
+-------------------------------------------------------------------------+
```

## The "Why"

Containers achieve isolation using Linux Namespaces. For networking, a **network namespace** provides the container with its own private network stack, including its own `localhost` interface, IP address, and routing table.

Because the container's network is isolated, you cannot access its ports from the host's `localhost`. The `nsenter` command is a powerful utility that allows you to "enter" one or more namespaces of a target process and execute a command there.

## Step 1: Find the Container's Host PID

You need the Process ID (PID) of the container's main process as it runs on the host. There are several ways to find it, from high-level to low-level.

### Method 1: Using `docker inspect` (High-Level)

This is the most common and straightforward way when using Docker.

```sh
docker inspect -f '{{.State.Pid}}' <container_name_or_id>
```

### Method 2: Using `ctr` (Containerd/Mid-Level)

If you are interacting with `containerd` directly, you can list the running tasks. Docker uses the `moby` namespace by default.

```sh
# Use grep to find your specific container
sudo ctr -n moby task ls | grep <container_name>
```
The second column in the output is the PID.

### Method 3: Using `lsns` (OS/Low-Level)

You can list the system's network namespaces and find the one running your container's command.

```sh
# Grep for the known process name, e.g., "nginx", "python"
sudo lsns -t net -o NS,NPROCS,PID,COMMAND | grep <process_name>
```
The third column in the output is the PID.

## Step 2: Enter the Namespace and Execute the Command

Once you have the PID, use `nsenter` to run a command inside the container's network namespace.

- `--target <PID>`: Specifies the process whose namespaces we want to enter.
- `--net`: Specifies we want to enter the network namespace.

```sh
# Replace <PID> and <port> accordingly
sudo nsenter --target <PID> --net curl http://localhost:<port>
```

You can also use the direct path to the namespace file, which can be found via `lsns`.

```sh
# The path is typically /proc/<PID>/ns/net
sudo nsenter --net=/proc/<PID>/ns/net curl http://localhost:<port>
```

## Full Example

1.  **Start an `nginx` container with no published ports:**
    ```sh
    docker run --name my-secret-nginx -d nginx
    ```

2.  **Get its PID using any method:**
    ```sh
    NGINX_PID=$(docker inspect -f '{{.State.Pid}}' my-secret-nginx)
    echo "Nginx PID is: $NGINX_PID"
    ```

3.  **Execute an HTTP request inside its network namespace:**
    ```sh
    sudo nsenter --target $NGINX_PID --net curl http://localhost:80
    ```

You will see the default Nginx welcome page HTML, proving you successfully accessed the web server in its isolated network.
