#!/bin/bash

./docker_stop.sh

docker build -t "projects:github-watchlists" .
docker run -p 9000:9000 -d --name github-watchlists projects:github-watchlists
