#!/usr/bin/env bash

docker build -t bentsolheim/kildsundvaeret-data-collector .
docker run \
 --rm \
 --name kildsundvaeret-data-collector \
 bentsolheim/kildsundvaeret-data-collector:latest
