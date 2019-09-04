#! /usr/bin/env bash

set -eux -o pipefail

CONTAINER_NAME=s3_anti_virus

docker pull amazonlinux:latest

docker build . --tag "${CONTAINER_NAME}"

docker rm -f "${CONTAINER_NAME}" || true
docker run --name "${CONTAINER_NAME}" "${CONTAINER_NAME}" 

# Get the files
mkdir -p clamav
docker cp "${CONTAINER_NAME}":/usr/bin/clamscan ./clamav/
docker cp "${CONTAINER_NAME}":/usr/bin/freshclam ./clamav/
docker cp -L "${CONTAINER_NAME}":/usr/lib64/libclamav.so.9 ./clamav/
docker cp "${CONTAINER_NAME}":/var/lib/clamav .
rm ./clamav/*.dat
docker stop "${CONTAINER_NAME}"

# Cleanup!
docker rm -f "${CONTAINER_NAME}"
docker rmi "${CONTAINER_NAME}"
