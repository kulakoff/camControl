# base image
FROM golang:1.22 AS builder

WORKDIR /app

# Copy go.mod и go.sum for load dependency
COPY go.mod go.sum ./
RUN go mod download

#  opy project code
COPY . .

# make app
RUN CGO_ENABLED=0 GOOS=linux go build -o ptz-camera-service ./cmd

# Final image
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/ptz-camera-service .
EXPOSE 8080
CMD ["./ptz-camera-service"]