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
