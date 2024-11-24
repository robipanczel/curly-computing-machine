FROM golang:1.22-alpine AS builder

WORKDIR /src

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/api/main.go

FROM alpine:3.18

ARG PORT

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE ${PORT}

CMD ["./main"]