### Tag and Push Container Images using `ctr`

#### Tagging an Image

To tag a container image using `ctr`, use the following command:

```bash
sudo ctr image tag <image_id> <the_tag>
```

#### Pushing an Image

When pushing an image, it's important to specify the correct platform, or pull with `--all-platforms`.

```bash
sudo ctr image push -u iximiuzlabs:rules! --platform amd64 registry.iximiuz.com/nginx:foo
```

#### Deleting an Image

To delete a container image using `ctr`, use the following command:

```bash
sudo ctr image delete <image_id>
```