* Openssl 簽發中繼憑證 config 檔
  sign-ssl.conf

# 產生 csr 及 key
```
openssl req -sha256 -utf8 -nodes -newkey rsa:2048 -keyout ./Unicorn-IM-CA-key.pem -out ./Unicorn-IM-CA.csr -config ./sign-ssl.conf

openssl req -sha256 -utf8 -nodes -newkey rsa:2048 -keyout ./K8S-IM-CA-key.pem -out ./K8S-IM-CA.csr -config ./sign-ssl.conf
```

# 自行加簽 (-CAcreateserial 只有簽第一個的時侯要加，簽過之後，如果加上這個參數， serial 號會重頭開始)
```
openssl x509 -req -in ./Unicorn-IM-CA/K8S-IM-CA/K8S-IM-CA.csr -CA ./Unicorn-IM-CA/Unicorn-IM-CA.pem -CAkey ./Unicorn-IM-CA/Unicorn-IM-CA-key.pem -out ./Unicorn-IM-CA/K8S-IM-CA/K8S-IM-CA.pem \
-CAserial Unicorn-IM-CA.serial -CAcreateserial \
-days 730 \
-extensions intermediate_ca -extfile K8S-IM.ext
```
extensions 這裡指的是要使用 extfile 內的什麼 session


# 說明
主要是利用openssl 產生 key 及 csr 檔，再去 內部自動憑證申請(http://dc01/certsrv/)完成 加簽動作，取得已簽好的 pem 檔
後續會利用此可繼續加簽的 Unicorn-IM-CA 往下加簽各種憑證，如：K8S, etcd, proxy, web side, docker, git.....等

已生成可繼續往下加簽的 中繼憑證：
1. Unicorn-IM-CA
2. K8S-IM-CA


** cfssl 指令：利用上述中繼憑證繼續簽出指定的 etcd ca 進行初步測試是否正確

cfssl gencert -ca=Unicorn-IM-CA.pem -ca-key=Unicorn-IM-CA-key.pem -config=ca-config.json -profile=intermediate-ca ./ca-csr.json | cfssljson -bare etcd-ca

** 測試加簽指令
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=intermediate-ca ./test-csr.json | cfssljson -bare k8s-ca

cfssl gencert -initca etcd-ca-csr.json | cfssljson -bare etcd-ca

cfssl sign -ca=../intermediate/IM-CA.pem -ca-key=../intermediate/IM-CA-key.pem -config=/tmp/ca-config.json -profile=intermediate-ca etcd-ca.csr | cfssljson -bare etcd-ca

cfssl gencert -ca=../intermediate/IM-CA.pem -ca-key=../intermediate/IM-CA-key.pem -config=/tmp/ca-config.json -profile=intermediate-ca etcd-ca-csr.json | cfssljson -bare etcd-ca