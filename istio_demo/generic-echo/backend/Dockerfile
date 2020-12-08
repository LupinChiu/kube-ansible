FROM golang:1.15.5 as builder
WORKDIR /go/src/echo_server/
COPY /echo_server/ .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o generic_echo_server main.go

FROM scratch as runner
WORKDIR /bin/
COPY --from=builder /go/src/echo_server/generic_echo_server .
COPY --from=builder /go/src/echo_server/conf/server.yaml .
ENTRYPOINT [ "/bin/generic_echo_server" ]
CMD [ "8033", "-c", "server.yaml" ]
EXPOSE 8033