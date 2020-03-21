FROM golang
WORKDIR /golang/src
ADD ./cmd/* ./cmd
ADD ./github.com/* ./github.com/*
ENV GOPATH /golang/src
RUN go run cmd/consumer

# Later add support in:
# user & password for redis & rabbitmq.
# and more variables.