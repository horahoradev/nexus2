FROM golang

RUN apt-get update && \
    apt install -y protobuf-compiler && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1 && \
    apt-get install -y python3 python3-pip && \
    python3 -m pip install --upgrade pip && \
    python3 -m pip install grpcio && \
    python3 -m pip install grpcio-tools

COPY gen_all.sh /bin/gen_all.sh

WORKDIR /

ENTRYPOINT ["/bin/gen_all.sh"]