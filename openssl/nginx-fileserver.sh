#!/usr/bin/env bash
# http://terrychen.logdown.com/posts/2015/10/04/nginx-http-file-server
# http://pki.urad.local => 10.0.1.154

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

entry_port=80
proxy_target=127.0.0.1
proxy_port=8081
githubUserName=uradLupin
githubPwd=ptUuTPt59Z8ADu4

fileserverBasePath=/data/pki-ssl
fileserverPath=${fileserverBasePath}/files
fileserverConfigPath=${fileserverBasePath}/config
mkdir -p $fileserverPath/$entryPath
mkdir -p $fileserverConfigPath

if [ ! -d "${fileserverPath}/.nginx" ]; then
  cd /tmp
#  git clone https://github.com/nervo/nginx-indexer.git
  git clone https://$githubUserName:$githubPwd@github.com/nervo/nginx-indexer.git
  mv nginx-indexer/.nginx ${fileserverPath}/
  rm -rf nginx-indexer
fi

cd $fileserverConfigPath

echo "server {
  listen $entry_port;
  server_name $entry_domain;
  location / {
    proxy_pass http://$proxy_target:$proxy_port; 
    proxy_set_header Host \$host;  
    proxy_set_header X-Real-IP \$remote_addr;  
    proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;  
    proxy_set_header X-Forwarded-Proto \$scheme;  
    proxy_redirect    off;  
  }
}" > default.conf

echo "server {
  listen $proxy_port;
  server_name $proxy_target;

  location / {
    root   /home/ssl-web;
    add_before_body /.nginx/header.html;
    add_after_body /.nginx/footer.html;
    autoindex_exact_size off;
    autoindex on;
  }
} #indexer" > nginx-http-ssl-repository.conf

docker run --name hfs -v ${fileserverBasePath}/config:/etc/nginx/conf.d -v ${fileserverPath}:/home/ssl-web:ro -d -p 80:80 nginx