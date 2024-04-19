FROM golang:1.22.2-bullseye as builder
WORKDIR /opt
COPY . .
RUN apt-get update && apt-get install -y upx-ucl
RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -ldflags '-linkmode "external" -extldflags "-static"' -o /opt/s9cmd .
RUN upx /opt/s9cmd
RUN chmod a+x /opt/s9cmd
FROM debian:bullseye-slim
COPY --from=builder /opt/s9cmd /opt/s9cmd
ENTRYPOINT ["/opt/s9cmd"]