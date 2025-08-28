FROM golang:1.23 AS builder
WORKDIR /back
COPY ./back/go.mod ./back/go.sum ./
COPY ./.env ./
RUN go mod download
COPY ./back .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .


FROM scratch
COPY --from=builder /back/main .
COPY --from=builder /back/.env .
CMD ["./main"]