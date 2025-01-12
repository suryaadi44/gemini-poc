FROM golang:1.23.4 as builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o /bin/main -v ./cmd/main/
RUN mkdir config && \
    cp config.yml config/

FROM debian:stable-slim

# Copy the binary from the builder image
COPY --from=builder /bin/main .
COPY --from=builder config .

CMD [ "sh", "-c", "./main" ]