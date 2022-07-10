#!/bin/bash
set -e

MAGE_IMAGE=pp-builder:latest
TARGET_ARCH=$(uname -m)

build_mage_image() {
  docker build -t $MAGE_IMAGE -f ./infra/dev/mage-docker/Dockerfile --build-arg TARGET_ARCH=${TARGET_ARCH} ./infra/dev/mage-docker/
}

# in case if a need the mage docker image can be rebuilt.
if [[ $1 == "--build-mage" ]]; then
  build_mage_image
  exit 0
fi

# mage docker image will be built automatically if it does not exist.
IMAGE_EXISTS=$(docker images -q $MAGE_IMAGE)
if [[ ! "$IMAGE_EXISTS" ]]; then
  echo "Mage docker image does not exist, building it now..."
  build_mage_image
fi

MAGE_SESSION_ID=$(uuidgen | cut -c 1,2,3 | tr "[:upper:]" "[:lower:]")
MAGE_NETWORK_ID=$(docker network create "${MAGE_SESSION_ID}")
printf "Session ID: $MAGE_SESSION_ID\nNetwork ID: $MAGE_NETWORK_ID\n\n"

cleanup () {
    docker ps --format '{{.Names}}' | grep "^$MAGE_SESSION_ID-" | awk '{print $1}' | xargs -I {} docker rm -f {} > /dev/null
    docker image list --format '{{.Repository}}:{{.Tag}}' | grep ":$MAGE_SESSION_ID\$" | awk '{print $1}' | xargs -I {} docker rmi -f {} > /dev/null
    docker network rm "$MAGE_NETWORK_ID" > /dev/null
    echo "Cleanup complete"
}
trap cleanup EXIT

# GID to be added to user groups in the running container
# so that the user can interact with docker.
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
   docker_sock="/var/run/docker.sock"
   docker_gid=$(stat -c "%g" $docker_sock)
elif [[ "$OSTYPE" == "darwin"* ]]; then
   docker_sock="/var/run/docker.sock.raw"
   docker_gid=0
else
   echo "Unsupported OS"
   exit 1
fi

# Detect interactive
[[ -t 0 ]] && interactive='-it'

docker run --rm \
  $interactive \
  --name "$MAGE_SESSION_ID-mage" \
  --network $MAGE_NETWORK_ID \
  --user $(id -u):$(id -g) \
  --group-add $docker_gid \
  --volume ${docker_sock}:/var/run/docker.sock \
  --volume "$PWD":/home/mage \
  --env REPO_DIRECTORY="$PWD" \
  --env MAGE_NETWORK_ID="$MAGE_NETWORK_ID" \
  --env MAGE_SESSION_ID="${MAGE_SESSION_ID}" \
  --env TARGET_ARCH="${TARGET_ARCH}" \
  $MAGE_IMAGE "$@"
