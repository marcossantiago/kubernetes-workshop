#!/usr/bin/env bash -x

NAME="kubernetes-training"

docker kill $NAME
docker rm $NAME
docker build -t nodejs-training .
docker run --rm --name $NAME -d  -p 8000:8000 -v "$PWD":/home/node/app nodejs-training
