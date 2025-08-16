# Challenge: Extract an Entire Filesystem from a Docker Image

This guide explains how to extract the complete filesystem from a Docker image onto your host machine. This is the closest equivalent to `ctr image mount` when using Docker. While Docker doesn't have a direct "mount image" command, you can achieve the same result by creating a container and exporting its filesystem.

## Step 1: Create a (non-running) container

First, you need to create a container from the image you want to inspect. You don't need to run it. The `docker create` command is perfect for this.

```shell
# docker create --name <container_name> <image_name>
docker create --name my-temp-container registry.iximiuz.com/tricky-one:latest
```

This creates a container named `my-temp-container` in the "Created" state. You can see it with `docker ps -a`.

## Step 2: Export the container's filesystem

Next, use the `docker export` command to create a tarball of the container's entire filesystem.

```shell
# docker export -o <output_file.tar> <container_name>
docker export -o tricky-fs.tar my-temp-container
```

This command creates a file named `tricky-fs.tar` containing the full filesystem of the `my-temp-container`.

## Step 3: Extract the filesystem on the host

Finally, you can extract the contents of the tarball into a directory on your host.

```shell
# Create a directory to hold the filesystem
mkdir /mnt/my_image

# Extract the tarball into the directory
tar -xf tricky-fs.tar -C /mnt/my_image
```

Now, the directory `/mnt/my_image` contains the complete filesystem from the container image, which you can browse and inspect.

## Step 4: Clean up

Once you're done, you can remove the temporary container and the tarball.

```shell
docker rm my-temp-container
rm tricky-fs.tar
```
