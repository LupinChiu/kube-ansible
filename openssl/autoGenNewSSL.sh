#!/bin/bash

function help_msg() {
    echo "example:"
    echo "./autoGenNewSSL.sh -S Unicorn-IM-CA -T udream-local"
    echo ""
    echo "-S 參數為 SSL 簽證者"
    echo "-T 參數為 要自動 gen 的 ssl簽證名"
    echo "-e 參數為 .conf 要使用的 session 目前有兩個選項 end_use & intermediate_ca, 預設為 end_use"
    exit 0
}

extSession="end_use"

while getopts "S:T:e:h?" arg; do
    case $arg in
        S) _SignerName=$OPTARG;;
        T) _TargetName=$OPTARG;;
        e) extSession=$OPTARG;;
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
sed -i "s/###x509ext###/${extSession}/g" ${targetFileName}.conf

# change ca path
sed -i "s/http:\/\/pki.unicorn.udream.local\//http:\/\/pki.unicorn.udream.local\/${_SignerName}.crt/g" ${targetFileName}.conf

#opensslExtParam="-extensions ${extSession} -extfile ${targetFileName}.ext"
opensslExtParam="-extensions ${extSession} -extfile ${targetFileName}.conf"
# start gen key

if [ ! -f ${targetFileName}-key.pem ]; then
    openssl req -sha256 -utf8 -nodes -newkey rsa:2048 -keyout ${targetFileName}-key.pem -out ${targetFileName}.csr -config ${targetFileName}.conf
fi

if [ -f ${targetFileName}.pem ]; then
    echo "${targetFileName}.pem is exist!"
    exit 1
fi

# del [x509ext]
sed -i "/\[x509ext\]/d" ${targetFileName}.conf

openssl x509 -req -in ${targetFileName}.csr -CA ${sourceFileName}.pem -CAkey ${sourceFileName}-key.pem -out ${targetFileName}.pem \
-CAserial ${sourceFileName}.serial -days 720 ${opensslExtParam}

cat ${targetFileName}.pem ${sourceFileName}-bundle.pem >> ${targetFileName}-bundle.pem

echo "gen ssl finish. check folder:${targetDir}"