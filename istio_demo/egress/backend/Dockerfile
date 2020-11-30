FROM golang:1.15.5 as builder
WORKDIR /go/src/simple_login/
COPY /simple_login/main.go .
RUN go mod init main
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o webapp main.go


FROM scratch as runner
WORKDIR /bin/
COPY --from=builder /go/src/simple_login .
ENTRYPOINT [ "/bin/webapp" ]
CMD [ "8033" ]
EXPOSE 8033