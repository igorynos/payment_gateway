FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/app \
    && CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/worker ./cmd/worker

FROM alpine:3.22

RUN apk add --no-cache ca-certificates \
    && addgroup -S payment \
    && adduser -S -G payment payment

WORKDIR /app

COPY --from=builder /out/app /app/app
COPY --from=builder /out/worker /app/worker
COPY deploy/config.yml /app/config.yml

# cmd/app expects an .env file even when all variables come from Compose.
RUN touch /app/.env && chown -R payment:payment /app

USER payment

EXPOSE 9000 9001

CMD ["/app/app"]
