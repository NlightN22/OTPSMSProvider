# build stage
FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go test ./... && \
    CGO_ENABLED=0 GOOS=linux go build -o otp-sms-provider .

# final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/otp-sms-provider .
EXPOSE 8080
ENTRYPOINT ["./otp-sms-provider"]
