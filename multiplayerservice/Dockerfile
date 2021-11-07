FROM golang:1.15.2-buster

COPY . /multiplayerservice

WORKDIR /multiplayerservice

RUN go mod vendor && \
    go build --mod=vendor -gcflags "all=-N -l" -o multiplayerservice .

ENTRYPOINT ["./multiplayerservice"]