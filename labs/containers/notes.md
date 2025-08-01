#### How to work with container images using ctr

#### List images

```
sudo ctr images ls

```

#### Build images

It does not provide a out-of-the-box image building functionality

```

docker build -t example.com/iximiuz/test:latest - <<EOF
FROM busybox:latest
CMD ["echo", "just a test"]
EOF

docker save -o iximiuz-test.tar example.com/iximiuz/test:latest

sudo ctr image import iximiuz-test.tar

```

### Mount images to see the content as one os file

```

mkdir /tmp/nginx_rootfs
sudo ctr image mount docker.io/library/nginx:latest /tmp/nginx_rootfs



```


### Tag images

The command to tag an image is:

```bash
ctr images tag <source_image> <new_tag>
```

This command creates a new tag (an alias) for an existing image. It's a very efficient operation as it doesn't duplicate the image data.

From your notes, the example is:

```bash
ctr images tag example.com/iximiuz/test:latest registry.iximiuz.com/test:latest
```

Here is a breakdown of the command:

*   `ctr images tag`: The `ctr` command for tagging images.
*   `example.com/iximiuz/test:latest`: This is the source image you want to tag.
*   `registry.iximiuz.com/test:latest`: This is the new tag you are applying to the source image.

A common reason to tag an image is to prepare it for pushing to a different container registry, which is what this example shows by changing the registry part of the name from `example.com/iximiuz` to `registry.iximiuz.com`.

After running the tag command, you can verify it by listing the images:

```bash
sudo ctr images ls
```

You will see both `example.com/iximiuz/test:latest` and `registry.iximiuz.com/test:latest` in the output, pointing to the same image digest.

### Where are image labels stored?

When you export an image, the labels are not in a single, obvious file. They are stored within the image's configuration file, which is a JSON file located in the `blobs/sha256/` directory of the exported tarball.

However, the easiest way to read an image's labels is to use the `ctr image inspect` command, which provides a full JSON output of the image's metadata.

```bash
ctr image inspect <image_name> | jq '.Spec.Config.Labels'
```
This is more direct than exporting and searching through the tar file.

### Troubleshooting Killed Pods

When a pod is killed, it's often due to resource constraints or errors. Here are some commands to help you investigate.

#### Check kubelet logs

The `kubelet` is the primary "node agent" that runs on each node. It's responsible for managing pods and their containers.

```bash
sudo journalctl -u kubelet -n 100 --no-pager
```

Look for messages related to the pod in question, especially any "Out of Memory" (OOM) messages.

#### Inspecting containers with crictl

`crictl` is a command-line interface for CRI-compatible container runtimes. It's a more direct way to inspect containers on a Kubernetes node than using `docker` or `ctr`.

First, list all pods to find the ID of the pod you're interested in:

```bash
sudo crictl pods
```

Then, inspect the pod to get more details, including the state of its containers:

```bash
sudo crictl inspectp <POD_ID>
```

If a container in the pod has been killed, you can look at its logs:

```bash
sudo crictl logs <CONTAINER_ID>
```

You can also inspect the container itself to see details about its state, including the reason it was terminated:

```bash
sudo crictl inspect <CONTAINER_ID>
```

Look for the `reason` field in the output.


### Challenges

- [How to manually download a container image manifest using cURL and how to copy an image between registries using `crane`](./copy-from-container-to-another.md) 
