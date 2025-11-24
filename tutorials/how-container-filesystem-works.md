### Tutorial: How Container Filesystems works: Building a docker-like container from scratch

#### Tools

- unsahre
- mount 
- pivot_root



![container-image.png](./images/container-image.png)

#### Command to mount a filesystem

```bash
mount --bind /tmp /mnt
```

This command uses `mount --bind` to create a bind mount. A bind mount makes a file or directory accessible at an alternative location in the filesystem. In this case, the content of `/tmp` is mirrored to `/mnt`.





#### Command to check mount points

```bash

cat /proc/self/mountinfo
```

