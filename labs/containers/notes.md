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

 
