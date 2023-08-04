FROM golang as builder

WORKDIR /app
COPY . /app

RUN go mod download
RUN go build -o main websocket.go query.go database.go dump.go command.go

FROM google/cloud-sdk:alpine

WORKDIR /app
COPY --from=builder /app/index.html /app/index.html
COPY --from=builder /app/main /app/main

CMD ["/app/main"]
