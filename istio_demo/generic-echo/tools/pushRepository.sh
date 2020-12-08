#!/bin/sh
USERNAME=test123
PWD=test123
#PRIVATE_URL=grd-dev.urad.com.tw
PUBLIC_URL=docker.io/whitewalker0506

APP_NAME=$1
VERSION=$2

if [[ -v PUBLIC_URL ]];
then
    cd ../backend
    docker build -t $PUBLIC_URL/$APP_NAME"_server":$VERSION -f ../backend/Dockerfile .
    docker push $PUBLIC_URL/$APP_NAME"_server":$VERSION
    
    cd ../frontend
    docker build -t $PUBLIC_URL/$APP_NAME"_client":$VERSION -f ../frontend/Dockerfile .
    docker push $PUBLIC_URL/$APP_NAME"_client":$VERSION
else
    cd ../backend
    docker login -u $USERNAME -p $PWD $PRIVATE_URL
    docker build -t $PRIVATE_URL/hank/$APP_NAME"_server":$VERSION -f ../backend/Dockerfile .
    docker push $PRIVATE_URL/hank/$APP_NAME"_server":$VERSION

    cd ../frontend
    docker login -u $USERNAME -p $PWD $PRIVATE_URL
    docker build -t $PRIVATE_URL/hank/$APP_NAME"_client":$VERSION -f ../frontend/Dockerfile .
    docker push $PRIVATE_URL/hank/$APP_NAME"_client":$VERSION
fi

