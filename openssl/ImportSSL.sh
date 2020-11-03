#!/usr/bin/env bash

source ./define.sh
if [ -z ${entry_domain} ]; then 
    entry_domain=pki.urad.local
fi
if [ -z ${entryPath} ]; then 
    entryPath=ssl-repository
fi
if [ -z ${cert_ext_name} ]; then 
    cert_ext_name="pem"
fi
if [ -z ${key_ext_name} ]; then 
    key_ext_name="pem"
fi

verifySSL=()
verifySSLName=$1
if [ -z $1 ]; then
  verifySSL+=("Lupin-Root-CA.${cert_ext_name}")
  verifySSL+=("K8S-IM-CA.${cert_ext_name}")
  verifySSL+=("ca.${cert_ext_name}")
  verifySSL+=("etcd-ca.${cert_ext_name}")
  verifySSL+=("front-proxy-ca.${cert_ext_name}")
else
  for arg in "$@"; do
    verifySSL+=("$arg")
  done
fi

echo "start copy file"

for verifySSLName in "${verifySSL[@]}"; do
  echo "wget $verifySSLName"
  wget http://$entry_domain/$entryPath/${verifySSLName} -P /etc/pki/ca-trust/source/anchors/
done

update-ca-trust

#如果這個import是為了要做 docker registry 的驗證才需要重啟 docker
#systemctl restart docker
