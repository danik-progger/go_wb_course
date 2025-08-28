FROM golang:1.23 AS builder
WORKDIR /consumer
COPY ./consumer/go.mod ./consumer/go.sum ./
COPY ./.env ./
RUN go mod download
COPY ./consumer .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .


FROM scratch
COPY --from=builder /consumer/main .
COPY --from=builder /consumer/.env .
CMD ["./main"]