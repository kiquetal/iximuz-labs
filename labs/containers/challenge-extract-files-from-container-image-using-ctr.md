# Challenge: Extract Files from a Container Image using `ctr`

This guide walks you through the process of saving a container image from Docker, importing it into `containerd` using `ctr`, and then mounting the image to access its filesystem.

## Step 1: Save the Image from Docker

First, we'll save the container image to a local tarball.

```shell
docker save -o tricky.tar registry.iximiuz.com/tricky-one:latest
```

This command saves the `registry.iximiuz.com/tricky-one:latest` image to a file named `tricky.tar`.

## Step 2: Import the Image into `ctr`

Next, we'll import the tarball into `containerd`'s image store using the `ctr` command-line tool.

```shell
sudo ctr image import tricky.tar
```

This command needs to be run with `sudo` because `ctr` typically requires root privileges to interact with the `containerd` daemon.

## Step 3: Mount the Image Filesystem

Finally, we can mount the image's filesystem to a temporary location to inspect its contents.

```shell
sudo ctr image mount <image_existing_in_ctr> /mnt/my_image
```

**Note:**
*   Replace `<image_existing_in_ctr>` with the actual image name or digest that was imported in the previous step. You can list the available images using `sudo ctr image ls`.
*   The mount point (e.g., `/mnt/my_image`) must exist. You can create it with `sudo mkdir -p /mnt/my_image`.

After mounting, you can browse the contents of the container image under `/mnt/my_image`.