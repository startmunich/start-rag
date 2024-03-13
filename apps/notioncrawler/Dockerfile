FROM golang:1.21.3-alpine3.18 as builder

WORKDIR /app

LABEL maintainer="Khadim Fall <khadim.fall@startmunich.de>"

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/crawler


# Runtime Stage
FROM alpine:3.18

WORKDIR /root/

COPY --from=builder /app/main .

ENV PORT=5000

EXPOSE 5000

CMD ["./main"]