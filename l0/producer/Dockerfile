FROM golang:1.23 AS builder
WORKDIR /producer
COPY ./producer/go.mod ./producer/go.sum ./producer/orders ./
COPY ./.env ./
RUN go mod download
COPY ./producer .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .


FROM scratch
COPY --from=builder /producer/main /main
COPY --from=builder /producer/orders /orders
COPY --from=builder /producer/healthcheck.json /healthcheck.json
COPY --from=builder /producer/.env /.env
CMD ["./main"]