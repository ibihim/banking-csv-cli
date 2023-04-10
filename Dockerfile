# Build stage
FROM golang:1.20 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build

# Final stage
FROM alpine:3.14

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/build/banking .

EXPOSE 8080

CMD ["./banking"]

