FROM golang:1.23.4 as builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o /bin/main -v ./cmd/main/

FROM debian:stable-slim

# Copy the binary from the builder image
COPY --from=builder /bin/main .
RUN cp config.yml config/

CMD [ "sh", "-c", "./main" ]