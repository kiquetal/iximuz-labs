# How to Manually Download a Container Image Manifest using cURL

This guide demonstrates how to use `curl` to manually interact with a container registry's API, following the OCI Distribution Specification. This process is useful for debugging or understanding the low-level mechanics of how container images are stored and retrieved.

**For any practical, real-world use, you should always use a dedicated OCI client like `skopeo`, `podman`, or `docker`**, as they automate this complex and error-prone process.

## The `download.sh` Script

This script automates the manual `curl` steps required to fetch an image manifest. It correctly handles authentication, accepts the necessary OCI media types, and can parse a multi-architecture manifest list (also known as a "fat manifest") to find the manifest for a specific architecture.

### Usage

1.  Save the code below as `download.sh`.
2.  Make it executable: `chmod +x download.sh`.
3.  Run the script: `./download.sh`.

The script is pre-configured to download the manifest for the `linux/amd64` version of the `ghcr.io/iximiuz/labs/nginx:latest` image. You can change the environment variables inside the script to target a different image or architecture.

### Script: `download.sh`

```bash
#!/bin/bash
set -e
# set -x # Uncomment for verbose debugging

# --- Configuration ---
# The image you want to inspect
REGISTRY="ghcr.io"
IMAGE="iximiuz/labs/nginx"
TAG="latest"

# The desired architecture and OS
TARGET_ARCH="amd64"
TARGET_OS="linux"

# --- Main Script ---

echo "--- Step 1: Get Authentication Token ---"
# Even public registries like ghcr.io require a bearer token for API access.
TOKEN=$(curl -s "https://${REGISTRY}/token?scope=repository:${IMAGE}:pull" | jq -r '.token')

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
  echo "FATAL: Failed to get auth token. Cannot proceed."
  exit 1
fi
echo "Token received successfully."
echo

echo "--- Step 2: Fetch Initial Manifest (or Manifest List) ---"
# We must send the correct Accept headers to tell the server we can handle
# standard image manifests as well as multi-arch manifest lists (OCI Indexes).
MANIFEST_LIST_FILE="manifest-list.json"
curl -sL -H "Authorization: Bearer ${TOKEN}" \
         -H "Accept: application/vnd.oci.image.index.v1+json" \
         -H "Accept: application/vnd.docker.distribution.manifest.list.v2+json" \
         -H "Accept: application/vnd.oci.image.manifest.v1+json" \
         -H "Accept: application/vnd.docker.distribution.manifest.v2+json" \
         "https://${REGISTRY}/v2/${IMAGE}/manifests/${TAG}" -o "${MANIFEST_LIST_FILE}"

echo "Initial manifest downloaded to ${MANIFEST_LIST_FILE}."
# Check if this is a multi-arch manifest list (index) or a single manifest
MEDIA_TYPE=$(jq -r '.mediaType' ${MANIFEST_LIST_FILE})
echo "Detected media type: ${MEDIA_TYPE}"
echo

FINAL_MANIFEST_FILE="final-image-manifest.json"

if [[ "${MEDIA_TYPE}" == "application/vnd.oci.image.index.v1+json" || "${MEDIA_TYPE}" == "application/vnd.docker.distribution.manifest.list.v2+json" ]]; then
  echo "--- Step 3: Parse Multi-Arch Manifest List ---"
  # This is a multi-arch image. We need to find the digest for our target platform.
  IMAGE_MANIFEST_DIGEST=$(jq -r --arg ARCH "${TARGET_ARCH}" --arg OS "${TARGET_OS}" '.manifests[] | select(.platform.architecture == $ARCH and .platform.os == $OS) | .digest' ${MANIFEST_LIST_FILE})

  if [ -z "$IMAGE_MANIFEST_DIGEST" ]; then
    echo "FATAL: Could not find manifest for ${TARGET_OS}/${TARGET_ARCH} in the list."
    exit 1
  fi
  echo "Found digest for ${TARGET_OS}/${TARGET_ARCH}: ${IMAGE_MANIFEST_DIGEST}"
echo

  echo "--- Step 4: Fetch Architecture-Specific Manifest ---"
  curl -sL -H "Authorization: Bearer ${TOKEN}" \
           -H "Accept: application/vnd.oci.image.manifest.v1+json" \
           -H "Accept: application/vnd.docker.distribution.manifest.v2+json" \
           "https://${REGISTRY}/v2/${IMAGE}/manifests/${IMAGE_MANIFEST_DIGEST}" -o "${FINAL_MANIFEST_FILE}"
else
  # This was a single-architecture image, so the first file is the final one.
  echo "--- Step 3: Single-Architecture Manifest Found ---"
  mv "${MANIFEST_LIST_FILE}" "${FINAL_MANIFEST_FILE}"
fi

echo "Final image manifest saved to ${FINAL_MANIFEST_FILE}."
echo
echo "--- Final Manifest Content ---"
jq . "${FINAL_MANIFEST_FILE}"
echo "----------------------------"
echo
echo "Next step: Use the digests in '${FINAL_MANIFEST_FILE}' to download the config and layer blobs."

```



--- 

## How to Copy an Image Between Registries using `crane`

`crane` provides a simple and efficient way to copy images directly between container registries. The `crane copy` command streams the image from the source to the destination without needing to store it on your local disk, which saves time and space.

This is the recommended approach for moving images between registries, especially when dealing with different tags or private repositories.

### Usage

The command follows this format:

```bash
crane copy <source_image> <destination_image>
```

Before you begin, if your destination registry is private, you must first log in:

```bash
crane auth login <your-private-registry.com> -u <your-username>
```

### Example: Copying from `ghcr.io` to a Private Registry

Here is how to copy the `ghcr.io/iximiuz/labs/nginx:latest` image to `registry.iximiuz.com/third-party/nginx` and retag it as `alpine` in a single step.

1.  **Log in to the private registry:**
    ```bash
    crane auth login registry.iximiuz.com -u YOUR_USERNAME
    ```

2.  **Copy the image:**
    This command pulls the `:latest` tag from the source and pushes it to the destination with the `:alpine` tag.
    ```bash
    crane copy ghcr.io/iximiuz/labs/nginx:latest registry.iximiuz.com/third-party/nginx:alpine
    ```

This single command handles the entire process, including authentication with the destination registry (using the credentials you provided in the login step).

```
