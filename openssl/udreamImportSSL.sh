#!/bin/bash

verifySSLName=$1
if [ -z $1 ]; then
	verifySSLName=Unicorn-IM-CA-bundle.pem
fi

wget http://pki.unicorn.udream.local/${verifySSLName} -P /etc/pki/ca-trust/source/anchors/
update-ca-trust

#如果這個import是為了要做 docker registry 的驗證才需要重啟 docker
#systemctl restart docker
