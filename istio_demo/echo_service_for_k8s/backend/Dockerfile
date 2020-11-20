FROM mcr.microsoft.com/dotnet/sdk:3.1 as builder
COPY ./EchoServer/ .
WORKDIR ./EchoServer
RUN dotnet publish

FROM mcr.microsoft.com/dotnet/runtime:3.1 as runner
COPY --from=builder ./EchoServer/bin/Debug/netcoreapp3.1 ./bin/echoserver
WORKDIR ./bin/echoserver
ENTRYPOINT [ "dotnet", "EchoServer.dll" ]
CMD [ "8020" ]
EXPOSE 8020
