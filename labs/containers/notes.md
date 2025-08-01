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
