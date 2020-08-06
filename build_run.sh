#!/usr/bin/env bash

docker build -t bentsolheim/kilsundvaeret-data-collector .
docker run \
 --rm \
 --name kilsundvaeret-data-collector \
 --env DB_HOST=kilsundvaeret-api_db_1 \
 --env DATA_RECEIVER_URL=http://something.no \
 --env MET_PROXY_URL=http://kilsundvaeret-api_met-proxy_1:9010 \
 --net kilsundvaeret-api_default \
 bentsolheim/kilsundvaeret-data-collector:latest
