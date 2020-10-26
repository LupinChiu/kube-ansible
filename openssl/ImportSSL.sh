#!/bin/bash

verifySSL=()
verifySSLName=$1
if [ -z $1 ]; then
  verifySSL+=("Lupin-Root-CA.crt")
  verifySSL+=("K8S-IM-CA.crt")
  verifySSL+=("ca.pem")
  verifySSL+=("etcd-ca.pem")
  verifySSL+=("front-proxy-ca.pem")
else
  for arg in "$@"; do
    verifySSL+=("$arg")
  done
fi

echo "start copy file"

entry_domain=pki.urad.local
entryPath=ssl-repository

for verifySSLName in "${verifySSL[@]}"; do
  echo "wget $verifySSLName"
  wget http://$entry_domain/$entryPath/${verifySSLName} -P /etc/pki/ca-trust/source/anchors/
done

update-ca-trust

#如果這個import是為了要做 docker registry 的驗證才需要重啟 docker
#systemctl restart docker
