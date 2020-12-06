FROM golang:alpine as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o reverse-proxy .

FROM scratch

COPY --from=builder /build/reverse-proxy /app/

EXPOSE 9090

ENTRYPOINT ["/app/reverse-proxy"]
