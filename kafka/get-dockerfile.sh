#!/usr/bin/env sh

if [ "$#" -lt 1 ]
then
    printf "Image name needed\n" >&2
    exit 1
fi

image_id="$(docker images | grep "^$1 " | awk '{print $3}')"
if [ -z "$image_id" ]
then
    printf "Image not found\n" >&2
    exit 2
fi

docker run -v /var/run/docker.sock:/var/run/docker.sock --rm chenzj/dfimage "$image_id"

