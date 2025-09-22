### We create a continer without running it.

- CONT_ID=$(docker create --platform linux/mips64le ghcr.io/iximiuz/labs/redis)

### We export using docker 

- docker export ${CONT_ID} -o redis.tar.gz

### We untar in the selected directory

-  tar -xf redis.tar.gz -C ~/imagefs
