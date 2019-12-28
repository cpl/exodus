FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go build -o server cmd/exodus-server/*.go

FROM alpine
RUN adduser -S -D -H -h /exodus exodus
USER exodus
COPY --from=builder /build/server /exodus/server
WORKDIR /exodus
EXPOSE 5353/udp
CMD ["./server", "-v", "--port=5353"]
