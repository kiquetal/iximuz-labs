### Execute a Command in a Docker Container Using ctr

### First we obtain the information about process already running via docker

- sudo ctr -n moby task ls


### Then we use the following command

- sudo ctr -n moby task exec --exec-id 0 --tty <your-container-id> sh



