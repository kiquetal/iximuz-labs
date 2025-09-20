### First we need to download the contaienrd

- Following the information from the gihub getting started with containers

https://github.com/containerd/containerd/blob/main/docs/getting-started.md


### Add a systemd unit file to the OS

To register and use a systemd unit file for containerd, follow these steps:

1.  **Create the systemd unit file:**

    Create a file named `containerd.service` in `/etc/systemd/system/` with the following content:

    ```
    [Unit]
    Description=containerd container runtime
    Documentation=https://containerd.io
    After=network.target local-fs.target

    [Service]
    ExecStartPre=-/sbin/modprobe overlay
    ExecStart=/usr/local/bin/containerd
    KillMode=process
    Delegate=yes
    LimitNPROC=infinity
    TasksMax=infinity
    OOMScoreAdjust=-999

    [Install]
    WantedBy=multi-user.target
    ```

    You can create this file using a text editor like `nano` or `vi`, or by using `tee`:
    ```bash
    sudo tee /etc/systemd/system/containerd.service <<EOF
    [Unit]
    Description=containerd container runtime
    Documentation=https://containerd.io
    After=network.target local-fs.target

    [Service]
    ExecStartPre=-/sbin/modprobe overlay
    ExecStart=/usr/local/bin/containerd
    KillMode=process
    Delegate=yes
    LimitNPROC=infinity
    TasksMax=infinity
    OOMScoreAdjust=-999

    [Install]
    WantedBy=multi-user.target
    EOF
    ```

2.  **Reload systemd daemon:**

    After creating the unit file, inform systemd about the new service:
    ```bash
    sudo systemctl daemon-reload
    ```

3.  **Enable and start containerd:**

    Enable the containerd service to start on boot and then start it immediately:
    ```bash
    sudo systemctl enable containerd --now
    ```
    (This command combines `sudo systemctl enable containerd` and `sudo systemctl start containerd`)

4.  **Check the status of containerd:**

    Verify that the containerd service is running correctly:
    ```bash
    sudo systemctl status containerd
    ```

### Running a container with CNI

Once CNI is installed and configured on the host, you can run containers with network connectivity. Here's how to run an nginx container using `ctr` and attach it to the CNI network.

1.  **Pull the nginx image:**
    ```bash
    sudo ctr image pull docker.io/library/nginx:latest
    ```

2.  **Run the nginx container with CNI:**
    ```bash
    sudo ctr run --cni -d docker.io/library/nginx:latest nginx
    ```
    - `--cni`: Enables CNI networking for the container.
    - `-d`: Runs the container in detached mode.
    - `docker.io/library/nginx:latest`: The image to use.
    - `nginx`: The name of the container.

3.  **Verify the container is running and has an IP address:**
    You can inspect the container to get its IP address.
    ```bash
    sudo ctr container info nginx
    ```
    Look for the IP address in the container's network information.

### You need to sure to add the following installation for the cni

```bash

{
  "type": "bridge",
  "bridge": "bridge0",
  "name": "bridge",
  "isGateway": true,
  "ipMasq": true,
  "ipam": {
    "type": "host-local",
    "ranges": [
      [{"subnet": "172.18.0.0/24"}]
    ],
    "routes": [{"dst": "0.0.0.0/0"}]
  },
  "cniVersion": "1.0.0"
}
```




 sudo mv /etc/cni/net.d/bridge.cnf /etc/cni/net.d/10-bridge.conf 
