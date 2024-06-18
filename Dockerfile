FROM golang:1.22 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o ws-terminal .

FROM alpine:latest  

WORKDIR /app

COPY --from=builder /app/ws-terminal .

COPY public public

ENTRYPOINT ["./ws-terminal"]
