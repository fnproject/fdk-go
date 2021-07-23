#!/usr/bin/env bash
set -ex

# Artifactory functions as a caching proxy for DockerHub. Hence any image pulled from dockerhub will be cached in artifactory.
# The cached images will be removed as part of cleanup if not downloaded again within a particular time frame, currently set to 6 days.
# Hence, one may encounter rate limiting issue while accessing the python docker images which are not present in artifactory.
# In order to resolve the rate limiting issue, the below script helps to pull the python images from docker hub and cache them in artifactory.

docker pull docker-remote.artifactory.oci.oraclecorp.com/golang:1.15-alpine3.13
docker pull docker-remote.artifactory.oci.oraclecorp.com/alpine:3.13
docker pull docker-remote.artifactory.oci.oraclecorp.com/golang:1.15-buster
