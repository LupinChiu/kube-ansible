FROM golang:1.15.5 as builder
WORKDIR /go/src/echo_client/
COPY /echo_client/ .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o generic_echo_client main.go

FROM alpine:3.12.1 as runner
WORKDIR /bin/
COPY --from=builder /go/src/echo_client/generic_echo_client .
COPY --from=builder /go/src/echo_client/conf/client.yaml .
ENTRYPOINT [ "/bin/generic_echo_client"]
CMD [ "-c", "client.yaml" ]
EXPOSE 5000