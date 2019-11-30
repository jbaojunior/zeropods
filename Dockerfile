FROM golang:latest as build

WORKDIR /go/src/github.com/jbaojunior/zeropods
COPY . .

ENV GO111MODULE=on
ENV CGO_ENABLED=0
RUN rm -f go.mod go.sum && go mod init && go get k8s.io/client-go@kubernetes-1.15.3
RUN go build -o /tmp/zeropods

FROM alpine:3.10
COPY --from=build /tmp/zeropods /usr/local/bin/

CMD ["/usr/local/bin/zeropods -h"]
ENTRYPOINT ["/bin/sh", "-c"]