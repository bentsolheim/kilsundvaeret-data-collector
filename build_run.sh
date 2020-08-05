#!/usr/bin/env bash

docker build -t bentsolheim/kilsundvaeret-data-collector .
docker run \
 --rm \
 --name kilsundvaeret-data-collector \
 --env DB_HOST=kilsundvaeret-data-collector_db_1 \
 --env DATA_RECEIVER_URL=http://something.no \
 --net kilsundvaeret-data-collector_default \
 bentsolheim/kilsundvaeret-data-collector:latest
