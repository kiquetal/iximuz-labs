#### Challenge

In this challenge, you will need to start a container using the default containerd CLI - ctr. Knowing how to use it may come in handy when you need to debug lower-level container issues (e.g., troubleshoot Kubernetes CRI on a containerd-powered cluster node).

#### First pull the image

sudo ctr image pull docker.io/library/nginx:alpine

sudo ctr create container docker.io/library/nginx:alpine nginx1

sudo ctr task start -d nginx1
