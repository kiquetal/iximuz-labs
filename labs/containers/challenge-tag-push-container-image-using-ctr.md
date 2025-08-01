### Tag and PUsh container images using ctr


#### To tag using ctr use the following

sudo ctr image tag <image_id> <the_tag>



### When push specify the correct platform or pull --all-platforms

sudo ctr image push  -u iximiuzlabs:rules! --platform amd64 registry.iximiuz.com/nginx:foo

### Delete image using ctr

sudo ctr image delete  <image_id>
