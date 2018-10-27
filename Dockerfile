FROM golang:1.11 as builder

ENV GOPATH="/go"\
    CGO_ENABLED='0'\
    GOOS='linux'\
    GOARCH='amd64'
    

WORKDIR /go/src/github.com/gonzalesraul/meow
COPY . .

RUN adduser --system --no-create-home --quiet appuser &&\
    go get ./... && \
    go build -a -installsuffix cgo -ldflags="-w -s" ./...

FROM scratch

COPY --from=builder /go/bin/*-service /bin/
COPY --from=builder /etc/passwd /etc/passwd

USER appuser
EXPOSE 8080