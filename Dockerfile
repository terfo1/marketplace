FROM golang:1.24 AS builder

WORKDIR /app

# Сначала модули
COPY go.mod go.sum ./
RUN go mod download

# Теперь весь проект
COPY . .

# Собираем статически
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/main.go

# ---- Stage 2: distroless ----
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Копируем бинарь
COPY --from=builder /app/server .

# Копируем .env (если используешь dotenv)
COPY .env .

EXPOSE 8080

USER nonroot:nonroot

CMD ["./server"]
