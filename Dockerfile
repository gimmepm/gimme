FROM golang:1.10
WORKDIR /go/src/github.com/gimmepm/gimme
COPY . .
RUN make

FROM debian:9.4-slim
WORKDIR /usr/local/bin
COPY --from=0 /go/src/github.com/gimmepm/gimme/bin/gimme .
ENTRYPOINT ["/usr/local/bin/gimme"]
