FROM golang:1.18 as builder

WORKDIR /app

ENV GO111MODULE=on

COPY . .

RUN go mod download
RUN go build -ldflags="-w -s" -o bin/jackbot/jackbot/main cmd/jackbot/main.go
RUN go build -ldflags="-w -s" -o bin/jackbot/migrate/main cmd/migrate/main.go

FROM golang:1.18 as app

WORKDIR /app

COPY --from=builder /go/pkg/mod /go/pkg/mod
COPY --from=builder /app /app

RUN app/bin/jackbot/migrate/main

CMD ["/app/bin/jackbot/jackbot/main"]