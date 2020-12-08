  
#!/bin/bash

echo ---- Gen gRPC Begin ----
SCRIPTPATH=$(readlink -f "$0")
BASEDIR=$(dirname $SCRIPTPATH)
SRCDIR=${BASEDIR}/../shared
CLIENT_OUTDIR=${BASEDIR}/../frontend/echo_client/proto
SERVER_OUTDIR=${BASEDIR}/../backend/echo_server/proto

if [ ! -d "${CLIENT_OUTDIR}" ];then
	mkdir -p ${CLIENT_OUTDIR}
elif [ ! -d "${SERVER_OUTDIR}" ]; then
	mkdir -p ${SERVER_OUTDIR}
fi

# gen server code
protoc --go_out=${SERVER_OUTDIR} --go_opt=paths=source_relative  \
    --go-grpc_out=${SERVER_OUTDIR} --go-grpc_opt=paths=source_relative  \
    -I ${SRCDIR}/ ${SRCDIR}/*.proto

# # gen client code
protoc --go_out=${CLIENT_OUTDIR} --go_opt=paths=source_relative  \
    --go-grpc_out=${CLIENT_OUTDIR} --go-grpc_opt=paths=source_relative  \
    -I ${SRCDIR}/ ${SRCDIR}/*.proto

echo ---- Gen gRPC End ----

