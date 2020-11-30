#!/bin/sh
USERNAME=test123
PWD=test123
PRIVATE_URL=grd-dev.urad.com.tw

SERVER=$1
VERSION=$2

cd ../backend
docker login -u $USERNAME -p $PWD $PRIVATE_URL
docker build -t $PRIVATE_URL/hank/$SERVER:$VERSION -f ../backend/Dockerfile .
docker push $PRIVATE_URL/hank/$SERVER:$VERSION