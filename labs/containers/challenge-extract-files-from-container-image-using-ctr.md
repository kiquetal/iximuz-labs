### Save tar to file using docker

docker save -o tricky.tar registry.iximiuz.com/tricky-one:latest


### Load using ctr

sudo ctr image import tricky.tar

### Mount using ctr

sudo ctr image mount <image_existing_in_ctr>
