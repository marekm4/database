FROM golang as builder

WORKDIR /app
COPY . /app

ENV CGO_ENABLED=0

RUN go mod download
RUN go build

FROM scratch

WORKDIR /app
COPY --from=builder /app/index.html /app/index.html
COPY --from=builder /app/database /app/database

CMD ["/app/database"]
