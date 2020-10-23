#!/bin/bash

function help_msg() {
    echo "example:"
    echo "./autoGenNewSSL.sh -S Unicorn-IM-CA -T udream-local"
    echo "./autoGenNewSSL.sh -S Lupin-Root-CA -T K8S-IM-CA -e intermediate_ca -d kudn"
    echo ""
    echo "-S 參數為 SSL 簽證者"
    echo "-T 參數為 要自動 gen 的 ssl簽證名"
    echo "-e 參數為 .conf 要使用的 session 目前有兩個選項 end_use & intermediate_ca, 預設為 end_use"
    echo "-d 參數為 .conf 要使用的 distinguished_name session , 預設為 normal_end_use"
    exit 0
}

extSession="end_use"
dnSession="normal_end_use"

while getopts "S:T:e:d:h?" arg; do
    case $arg in
        S) _SignerName=$OPTARG;;
        T) _TargetName=$OPTARG;;
        e) extSession=$OPTARG;;
        d) dnSession=$OPTARG;;
        \? | h | *) help_msg;;
    esac
done

if [ -z $_TargetName ]; then
    echo "TargetName 未設定"
    exit 1
fi

if [ -z $_SignerName ]; then
    echo "SignerName 未設定"
    exit 1
fi

sslFileSavePlace="http:\/\/pki.urad.local\/ssl-repository\/"
#sslFileSavePlace="http:\/\/pki.urad.com.tw\/ssl-repository\/"
exePath=$(dirname "$0")
curPath=$(pwd)

sourceDir="${exePath}/${_SignerName}"
sourceFileName="${sourceDir}/${_SignerName}"
if [ ! -d "${sourceDir}" ]; then
    echo "參數錯誤，來源目錄不存在"
    exit 1
fi

targetDir="${curPath}/${_TargetName}"
targetFileName="${targetDir}/${_TargetName}"
mkdir -p ${targetDir}
cp ${exePath}/temp-ssl.conf ${targetFileName}.conf
#cp ${exePath}/temp-ssl.ext ${targetFileName}.ext

# ###x509ext### change to real session
echo "changing x509 ext."
sed -i "s/###x509ext###/${extSession}/g" ${targetFileName}.conf

echo "changing distinguished_name session."
sed -i "s/###dn###/${dnSession}/g" ${targetFileName}.conf

# change ca path
echo "changing ssl ca path."
sed -i "s/###sslcapath###/${sslFileSavePlace}${_SignerName}.crt/g" ${targetFileName}.conf

echo "setting openssl ext param."
#opensslExtParam="-extensions ${extSession} -extfile ${targetFileName}.ext"
opensslExtParam="-extensions ${extSession} -extfile ${targetFileName}.conf"

# start gen key
if [ ! -f ${targetFileName}-key.pem ]; then
    echo "start gen ${targetFileName}-key.pem & ${targetFileName}.csr"
    openssl req -sha256 -utf8 -nodes -newkey rsa:2048 -keyout ${targetFileName}-key.pem -out ${targetFileName}.csr -config ${targetFileName}.conf
fi

if [ -f ${targetFileName}.pem ]; then
    echo "${targetFileName}.pem is exist!"
    exit 1
fi

# del [x509ext]
sed -i "/\[x509ext\]/d" ${targetFileName}.conf

if [ ! -f ${sourceFileName}.serial ]; then
    addSerialParam=" -CAcreateserial"
    echo "Add Create serial file param"
fi

openssl x509 -req -in ${targetFileName}.csr -CA ${sourceFileName}.pem -CAkey ${sourceFileName}-key.pem -out ${targetFileName}.pem \
-CAserial ${sourceFileName}.serial ${addSerialParam} -days 720 ${opensslExtParam}

cat ${targetFileName}.pem ${sourceFileName}-bundle.pem >> ${targetFileName}-bundle.pem

cp ${targetFileName}.pem ${targetFileName}.crt

echo "gen ssl finish. check folder:${targetDir}"