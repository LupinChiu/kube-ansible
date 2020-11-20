#!/bin/sh
USERNAME=test123
PWD=test123
PRIVATE_URL=grd-dev.urad.com.tw

CLIENT=echoclient
SERVER=echoserver
VERSION=2.0

docker login -u $USERNAME -p $PWD $PRIVATE_URL
# docker build -t echoserver:$VERSION -f ../backend/DockerfileGo ../backend/EchoServerGo
# docker build -t echoclient:$VERSION -f ../frontend/Dockerfile ../frontend/EchoClient
docker tag $CLIENT:$VERSION grd-dev.urad.com.tw/hank/echoclient:$VERSION
docker tag $SERVER:$VERSION grd-dev.urad.com.tw/hank/echoserver:$VERSION
docker push grd-dev.urad.com.tw/hank/echoclient:$VERSION
docker push grd-dev.urad.com.tw/hank/echoserver:$VERSION