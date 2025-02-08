FROM golang:1.23-alpine AS builder

WORKDIR /usr/local/src

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .
RUN go build -o ./bin/app cmd/main/main.go

FROM alpine AS runner

RUN apk add --no-cache iputils

COPY --from=builder /usr/local/src/bin/app /

EXPOSE 8081

CMD ["/app"]
