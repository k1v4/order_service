FROM golang:1.23-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY ["go.mod","go.sum","./"]
RUN go mod download

COPY . .
RUN go build -o ./bin/app cmd/main/main.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/app /

EXPOSE 50051 50052 50053 8080

COPY .env /

CMD ["/app"]