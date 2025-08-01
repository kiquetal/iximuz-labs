# Challenge: Add a label to a container image using ctr

## Goal

The goal of this challenge is to add a label to an existing container image using the `ctr` command-line tool.

## Background

Container image labels are key-value pairs stored in the image's metadata. They are useful for adding information like version, author, or a description. While `docker` provides a straightforward way to add labels during the build process (`docker build --label ...`), doing so with `ctr` on an *existing* image is different because `ctr` is a lower-level tool for interacting with the containerd daemon, not a full-featured image builder.

## Steps

1.  **Pull an image**

    First, let's pull a simple image to work with.

    ```bash
    sudo ctr images pull docker.io/library/busybox:latest
    ```

2.  **Inspect the image for labels**

    Let's check the existing labels on the `busybox` image. We can use `ctr images inspect` and pipe the output to `jq` to filter for the labels.

    ```bash
    sudo ctr images inspect docker.io/library/busybox:latest | jq '.Spec.Config.Labels'
    ```

    The output will likely be `null` as the base `busybox` image doesn't have any labels.

3. Add a label using ctr

 ctr image label <image_name> key=val



